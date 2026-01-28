package models

import "errors"

type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var articleList = []Article{
	{ID: 1, Title: "Article 1", Content: "Article 1 body"},
	{ID: 2, Title: "Article 2", Content: "Article 2 body"},
}

// ----- 모든 기사 목록을 반환합니다 ----- //
func GetAllArticles() []Article {
	return articleList
}

func GetArticleByID(id int) (*Article, error) {
	for _, a := range articleList {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("기사를 찾을 수 없습니다")
}

func CreateNewArticle(title, content string) {
	newID := len(articleList) + 1
	articleList = append(articleList, Article{ID: newID, Title: title, Content: content})
}

func DeleteArticleById(id int) {
	idx := 0
	for i, a := range articleList {
		if a.ID == id {
			idx = i
			break
		}
	}

	articleList = append(articleList[:idx], articleList[idx+1:]...)
}

func UpdateArticle(id int, title, content string) {
	for i, a := range articleList {
		if a.ID == id {
			articleList[i].Title = title
			articleList[i].Content = content
			break
		}
	}
}
