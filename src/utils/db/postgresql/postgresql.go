package postgresql

// postgresql localhost connection  code write
// Path: src/utils/db/postgresql/postgresql.go
// Compare this snippet from src/features/auth/repository/signupAuthRepository.go:

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var GormPGDB *gorm.DB

const DBTimeOut = 8 * time.Second

func InitPostgreSQL() error {

	// PostgreSQL 연결 문자열
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_DB_NAME"))
	// 실제 사용시에는 연결 정보를 안전하게 다루어야 합니다.
	// GormDB 초기화
	var err error
	GormPGDB, err = gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer GormPGDB.Close()

	// GormDB 연결 테스트
	err = GormPGDB.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to PostgreSQL using Gorm!")

	// GormDB 로그 모드 설정 (선택 사항)
	GormPGDB.LogMode(true)
	return nil
}
