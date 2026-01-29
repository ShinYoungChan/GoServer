package models

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
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
