package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlDB *sql.DB
var GormMysqlDB *gorm.DB

const DBTimeOut = 8 * time.Second

func InitMySQL() error {
	// MySQL 연결 문자열
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
	fmt.Println(connectionString)

	// MySQL에 연결
	MysqlDB, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	// MySQL 연결 테스트
	err = MysqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL!")
	/*
		GORM perform write (create/update/delete) operations run inside a transaction to ensure data consistency,
		you can disable it during initialization if it is not required, you will gain about 30%+ performance improvement after that
	*/
	GormMysqlDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: MysqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return err
	}
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
