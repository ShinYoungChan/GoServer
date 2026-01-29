package models

import (
	"fmt"
	"os" // 환경 변수를 읽기 위해 필요

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 1. .env 파일 로드
	err := godotenv.Load("postgre.env")
	if err != nil {
		fmt.Println(".env 파일을 찾을 수 없습니다. 시스템 환경 변수를 사용합니다.")
	}

	// 2. os.Getenv()로 값 가져오기
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		host, user, password, dbname, port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("PostgreSQL 연결 실패: " + err.Error())
	}

	DB.AutoMigrate(&Article{}, &User{})
}
