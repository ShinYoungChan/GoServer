package handlers

import (
	"fmt"
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
	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
		if article, err := models.GetArticleByID(articleID); err == nil {
			// 1. 현재 로그인한 유저 정보를 가져옴
			user := GetCurrentUser(c)
			c.HTML(
				http.StatusOK,
				"article.html", // 상세 페이지 설계도
				gin.H{
					"title":    article.Title,
					"articles": article, // 상세 페이지에선 이게 '글 하나'를 의미함
					"user":     user,    // 2. 유저 정보를 보따리에 넣어줍니다
				},
			)
		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
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
	user := GetCurrentUser(c)

	if user == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	title := c.PostForm("title")
	content := c.PostForm(("content"))

	article := models.Article{Title: title, Content: content, UserID: user.ID}
	models.DB.Create(&article)

	c.Redirect(http.StatusSeeOther, "/")
}

func PerformDeleteArticle(c *gin.Context) {
	user := GetCurrentUser(c)
	articleID := c.Param("article_id")

	var article models.Article
	models.DB.First(&article, articleID)

	// PerformArticleDelete 내부
	fmt.Printf("현재 로그인 유저 ID: %d\n", user.ID)
	fmt.Printf("글에 저장된 주인 ID: %d\n", article.UserID)

	if user == nil || article.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "본인 글만 삭제할 수 있습니다",
		})
		return
	}

	models.DB.Delete(&article)

	c.Redirect(http.StatusSeeOther, "/")
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
