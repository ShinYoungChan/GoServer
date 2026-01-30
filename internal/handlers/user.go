package handlers

import (
	"fmt"
	"gin/internal/models" // 본인 프로젝트 경로에 맞게
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowRegisterPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"register.html",
		gin.H{
			"title": "회원가입",
		},
	)
}

func PerformCreateUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	err := models.CreateUser(username, password)
	if err != nil {
		// 이미 존재하는 아이디일 가능성이 높음
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"Error": "이미 사용 중인 아이디거나 가입에 실패했습니다.",
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"login.html",
		gin.H{
			"title": "로그인",
		},
	)
}

func PerformLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	fmt.Printf("로그인 시도 ID: %s, PW 길이: %d\n", username, len(password))

	user, err := models.LoginCheck(username, password)

	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"Error": "아이디 또는 비밀번호가 일치하지 않습니다.",
		})
		return
	}

	fmt.Println("로그인 성공! 유저 ID: ", user.ID)

	// 로그인 성공 시 쿠키 설정 예시
	c.SetCookie("user_id", fmt.Sprintf("%d", user.ID), 3600, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}

func GetCurrentUser(c *gin.Context) *models.User {
	uid, err := c.Cookie("user_id")
	if err != nil {
		return nil
	}

	var user models.User
	// DB는 전역변수니까 models.DB로 접근!
	if err := models.DB.First(&user, uid).Error; err != nil {
		return nil
	}

	return &user
}

func PerformLogout(c *gin.Context) {
	// 쿠키의 Max-Age를 -1로 설정하면 브라우저가 즉시 삭제합니다.
	c.SetCookie("user_id", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}
