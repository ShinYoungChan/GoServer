package models

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// 2. 주소는 똑같습니다.
	DB, err = gorm.Open(sqlite.Open("article.db"), &gorm.Config{})
	if err != nil {
		// 에러 내용을 더 자세히 출력하면 디버깅이 쉬워집니다.
		panic("DB 연결 실패: " + err.Error())
	}

	DB.AutoMigrate(&Article{})
}
