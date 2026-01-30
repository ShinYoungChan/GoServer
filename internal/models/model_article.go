package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
}

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

// Article 구조체 하단에 추가
func (a Article) FormattedCreatedAt() string {
	// 2006-01-02 15:04:05는 Go 언어의 날짜 포맷 규칙입니다 (바꾸면 안 돼요!)
	return a.CreatedAt.Format("2006-01-02 15:04")
}

func (a Article) FormattedUpdatedAt() string {
	return a.UpdatedAt.Format("2006-01-02 15:04")
}

/* gorm model 추가로 주석처리
var articleList = []Article{
	{ID: 1, Title: "Article 1", Content: "Article 1 body"},
	{ID: 2, Title: "Article 2", Content: "Article 2 body"},
}
*/

// ----- 모든 기사 목록을 반환합니다 ----- //
func GetAllArticles() []Article {
	var articles []Article
	DB.Find(&articles) // SELECT * FROM articles;

	fmt.Println("가져온 글 개수:", len(articles))
	return articles
}

func GetArticleByID(id int) (*Article, error) {
	var article Article
	if err := DB.First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func CreateNewArticle(title, content string) (*Article, error) {
	newArticle := Article{Title: title, Content: content}
	if err := DB.Create(&newArticle).Error; err != nil {
		return nil, err
	}
	return &newArticle, nil
}

func DeleteArticleById(id int) {
	DB.Delete(&Article{}, id)
}

func UpdateArticle(id int, title, content string) {
	var article Article
	DB.First(&article, id)
	DB.Model(&article).Updates(Article{Title: title, Content: content})
}
