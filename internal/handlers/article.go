package handlers

import (
	"fmt"
	"math"
	"net/http"
	"os"
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
	articleID := c.Param("article_id")
	var article models.Article

	// ⭐ Preload("Comments")를 추가하면 Article 구조체의 Comments 필드에 데이터가 자동으로 채워집니다.
	if err := models.DB.Preload("Comments").First(&article, articleID).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.HTML(http.StatusOK, "article.html", gin.H{
		"articles": article, // 이제 .articles.Comments 로 템플릿에서 접근 가능!
		"user":     GetCurrentUser(c),
	})
	/*
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
	*/
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

	var imagePath string
	file, err := c.FormFile("image")

	if err == nil {
		// 2. 파일이 있다면 저장 경로 설정 (중복 방지를 위해 파일명 앞에 시간을 붙이기도 함)
		// 예: uploads/1672531200_photo.jpg
		dst := "uploads/" + file.Filename

		// 3. 서버 하드디스크에 파일 실제 저장
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.String(http.StatusInternalServerError, "파일 저장 실패")
			return
		}
		imagePath = "/" + dst // DB에는 나중에 웹에서 접근할 경로 저장
	}

	title := c.PostForm("title")
	content := c.PostForm(("content"))

	article := models.Article{Title: title, Content: content, UserID: user.ID, Image: imagePath}
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

	if article.Image != "" {
		filePath := article.Image[1:]

		fmt.Printf("파일경로: %s\n", filePath)

		err := os.Remove(filePath)

		if err != nil {
			fmt.Println("파일 삭제 실패:", err)
		}
	}

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
	user := GetCurrentUser(c)

	if user == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	id, _ := strconv.Atoi(c.Param("article_id"))

	var article models.Article
	models.DB.First(&article, id)

	title := c.PostForm("title")
	content := c.PostForm("content")

	// ---- 데이터 검증 로직 추가 ---- //
	if strings.TrimSpace(title) == "" {
		article.Title = title
		article.Content = content
		article.UserID = user.ID
		// 제목이 비어있거나 공백만 있다면 에러 메시지와 함께 중단
		c.HTML(http.StatusBadRequest, "edit_article.html", gin.H{
			"ErrorTitle": "제목은 필수입니다.",
			"articles":   article,
		})
		return
	}

	models.UpdateArticle(id, user.ID, title, content)

	// PerformUpdateArticle 마지막 줄
	c.Redirect(http.StatusSeeOther, "/article/view/"+strconv.Itoa(id))
}

func CreateComment(c *gin.Context) {
	user := GetCurrentUser(c)
	articleID, _ := strconv.Atoi(c.Param("article_id"))
	content := c.PostForm("content")

	comment := models.Comment{
		Username:  user.Username,
		Content:   content,
		ArticleID: uint(articleID),
		UserID:    user.ID,
	}

	models.DB.Create(&comment)

	// 저장 후 다시 보던 게시글 상세 페이지로 리다이렉트!
	c.Redirect(http.StatusSeeOther, "/article/view/"+strconv.Itoa(articleID))
}

func UpdateComment(c *gin.Context) {
	user := GetCurrentUser(c)
	commentID, _ := strconv.Atoi(c.Param("comment_id"))

	var comment models.Comment
	models.DB.First(&comment, commentID)

	if user.ID != comment.UserID {
		fmt.Printf("[ERROR] ID 불일치! USER ID[%s], COMMENT ID[%s]\n", user.ID, comment.UserID)
	}

	content := c.PostForm("content")

	models.UpdateComment(commentID, content)
	c.Redirect(http.StatusSeeOther, "/article/view/"+strconv.Itoa(int(comment.ArticleID)))
}

func DeleteComment(c *gin.Context) {
	user := GetCurrentUser(c)
	commentID := c.Param("comment_id")

	var comment models.Comment
	models.DB.First(&comment, commentID)

	result := models.DB.Where("id = ? AND user_id = ?", commentID, user.ID).Delete(&models.Comment{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "삭제 권한이 없습니다."})
		return
	}

	c.Redirect(http.StatusSeeOther, "/article/view/"+strconv.Itoa(int(comment.ArticleID)))
}

func GetIndex(c *gin.Context) {
	searchQuery := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 5
	user := GetCurrentUser(c) // 헬퍼 함수 사용!

	var articles []models.Article
	var totalCount int64 // 전체 게시글 수를 담을 변수

	// 1. 쿼리 초기화
	query := models.DB.Model(&models.Article{})

	// 2. 검색 조건 적용
	if searchQuery != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// 3. 페이징 적용 전! 전체 개수부터 센다 (이게 포인트!)
	query.Count(&totalCount)

	// 4. 페이징 적용해서 실제 데이터 가져오기
	offset := (page - 1) * pageSize
	query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&articles)

	// 5. 마지막 페이지 계산 (예: 13개면 3페이지까지 있어야 함)
	// float64로 변환해서 계산 후 올림(math.Ceil)하는 게 정석입니다.
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	var pageNums []int
	for i := 1; i <= totalPages; i++ {
		pageNums = append(pageNums, i)
	}

	// 핸들러 내부 예시 (이미 totalPages를 구한 상태라면)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"articles":   articles,
		"query":      searchQuery,
		"page":       page,
		"totalPages": totalPages,
		"pageNums":   pageNums,
		"hasPrev":    page > 1,          // 이전 페이지가 있는가? (bool)
		"hasNext":    page < totalPages, // 다음 페이지가 있는가? (bool)
		"prevPage":   page - 1,          // 이전 페이지 번호
		"nextPage":   page + 1,          // 다음 페이지 번호
		"user":       user,
	})
}
