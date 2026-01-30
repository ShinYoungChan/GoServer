package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"gin/internal/models"

	"github.com/gin-gonic/gin"
)

func ShowIndexPage(c *gin.Context) {
	articles := models.GetAllArticles()
	user := GetCurrentUser(c) // 헬퍼 함수 사용!

	c.HTML(http.StatusOK, "index.html", gin.H{
		"articles": articles,
		"user":     user, // 템플릿으로 유저 정보 전달
	})
}

func GetArticle(c *gin.Context) {
	// ----- 기사 ID가 유효한지 확인합니다 ----- //
	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
		// ----- 기사가 존재하는지 확인합니다 ----- //
		if article, err := models.GetArticleByID(articleID); err == nil {
			// Call the HTML method of the Context to render a template
			c.HTML(
				http.StatusOK,
				"article.html",
				gin.H{
					"title":    article.Title,
					"articles": article,
				},
			)

		} else {
			// ---- 기사를 찾을 수 없는 경우 오류와 함께 중단합니다 ---- //
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		// ---- URL에 잘못된 기사 ID가 지정된 경우 오류와 함께 중단합니다 ---- //
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func ShowArticleCreatePage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"create_article.html",
		gin.H{
			"title": "새 글 작성",
		},
	)
}

func PerformCreateArticle(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm(("content"))

	models.CreateNewArticle(title, content)

	c.Redirect(http.StatusMovedPermanently, "/")
}

func DeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("article_id"))
	models.DeleteArticleById(id)

	c.JSON(http.StatusOK, gin.H{"message": "삭제완료"})
}

func ShowArticleEditPage(c *gin.Context) {
	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
		if article, err := models.GetArticleByID(articleID); err == nil {
			c.HTML(http.StatusOK, "edit_article.html", gin.H{
				"articles": article,
			})
		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}
	}
}

func PerformUpdateArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("article_id"))
	title := c.PostForm("title")
	content := c.PostForm("content")

	// ---- 데이터 검증 로직 추가 ---- //
	if strings.TrimSpace(title) == "" {
		// 제목이 비어있거나 공백만 있다면 에러 메시지와 함께 중단
		c.HTML(http.StatusBadRequest, "edit_article.html", gin.H{
			"ErrorTitle": "제목은 필수입니다.",
			"articles":   models.Article{Title: title, Content: content}, // 기존 데이터 유지
		})
		return
	}

	models.UpdateArticle(id, title, content)

	c.Redirect(http.StatusMovedPermanently, "/article/view/"+strconv.Itoa(id))
}
