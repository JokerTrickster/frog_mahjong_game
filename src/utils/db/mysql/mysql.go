package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_aws "main/utils/aws"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlDB *sql.DB
var GormMysqlDB *gorm.DB

const DBTimeOut = 8 * time.Second

// 게임 타입 지정
const (
	FIND_IT   = 1
	SLIME_WAR = 2
)

func InitMySQL() error {
	var connectionString string
	var err error
	isLocal := os.Getenv("IS_LOCAL")
	if isLocal == "true" {
		// MySQL 연결 문자열
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DATABASE"),
		)
	} else {
		dbInfos, err := _aws.AwsSsmGetParams([]string{"dev_frog_mysql_user", "dev_common_mysql_password", "dev_common_mysql_host", "dev_common_mysql_port", "dev_frog_mysql_db"})
		if err != nil {
			return err
		}
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			dbInfos[4], //user
			dbInfos[1], //password
			dbInfos[0], //host
			dbInfos[2], //port
			dbInfos[3], //db
		)
	}
	fmt.Println(connectionString)
	// MySQL에 연결
	MysqlDB, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("Failed to connect to MySQL!")
		fmt.Sprintln("에러 메시지 %s", err)
	}
	fmt.Println("Connected to MySQL!")

	/*
		GORM perform write (create/update/delete) operations run inside a transaction to ensure data consistency,
		you can disable it during initialization if it is not required, you will gain about 30%+ performance improvement after that
	*/
	GormMysqlDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: MysqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: false,
	})
	if err != nil {
		fmt.Println("Failed to connect to Gorm MySQL!")
		fmt.Sprintln("에러 메시지 %s", err)
	}

	// gen 패키지를 사용하여 쿼리를 생성할 때 사용할 DB를 설정
	// SetDefault(GormMysqlDB)

	return nil
}

func PKIDGenerate() string {
	//uuid 로 생성
	result := (uuid.New()).String()
	return result
}

func NowDateGenerate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func EpochToTime(t int64) time.Time {
	return time.Unix(t, t%1000*1000000)
}
func EpochToTimeString(t int64) string {
	return time.Unix(t, t%1000*1000000).String()
}

func TimeStringToEpoch(t string) int64 {
	date, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", t)
	return date.Unix()
}

func TimeToEpoch(t time.Time) int64 {
	return t.Unix()
}

// 트랜잭션 처리 미들웨어
func Transaction(db *gorm.DB, fc func(tx *gorm.DB) error) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	if err = tx.Error; err != nil {
		return err
	}

	err = fc(tx)
	return
}
