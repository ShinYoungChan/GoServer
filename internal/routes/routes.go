package routes

import (
	"gin/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {

	// 인덱스 라우터 처리(Handle)
	r.GET("/", handlers.ShowIndexPage)

	aritlcleRoutes := r.Group("/article")
	{
		// /article/view/some_article_id 부분에 대한 GET 요청 처리
		aritlcleRoutes.GET("/view/:article_id", handlers.GetArticle)
		// article 생성 페이지
		aritlcleRoutes.GET("/create", handlers.ShowArticleCreatePage)
		// article 생성 이후 저장 post
		aritlcleRoutes.POST("/create", handlers.PerformCreateArticle)
		// 기존 r.DELETE를 r.POST로 변경
		//aritlcleRoutes.DELETE("/:article_id", handlers.PerformDeleteArticle)
		r.POST("/article/delete/:article_id", handlers.PerformDeleteArticle)
		// 수정 페이지 보여주기 (기존 데이터가 채워진 폼)
		aritlcleRoutes.GET("/edit/:article_id", handlers.ShowArticleEditPage)
		// 수정된 데이터 처리하기
		aritlcleRoutes.POST("/edit/:article_id", handlers.PerformUpdateArticle)
	}

	// 로그인
	r.GET("/login", handlers.ShowLoginPage)
	r.POST("/login", handlers.PerformLogin)

	// 회원가입
	r.GET("/register", handlers.ShowRegisterPage)
	r.POST("/register", handlers.PerformCreateUser)

	// 로그아웃
	r.GET("/logout", handlers.PerformLogout)
}
