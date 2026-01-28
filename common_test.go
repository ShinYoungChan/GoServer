package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var tmpArticleList []article

// ------- 이 함수는 테스트 함수를 실행하기 전, 설정에 사용합니다 ----- //
func TestMain(m *testing.M) {
	// ----- Gin을 테스트 모드로 설정합니다 ----- //
	gin.SetMode(gin.TestMode)

	// ------ 다른 테스트를 실행합니다 ------ //
	os.Exit(m.Run())
}

// ---- 테스트 중 라우터를 만드는 헬퍼 함수 ---- //
func getRouter(withTemplates bool) *gin.Engine {
	r := gin.Default()
	if withTemplates {
		r.LoadHTMLGlob("templates/*")
	}
	return r
}

// ----- 요청을 처리하고 응답을 테스트하는 헬퍼 함수 ------- //
func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	// ---- 응답 레코더 생성 ---- //
	w := httptest.NewRecorder()

	// ------ 서비스를 생성하고 위의 요청을 처리합니다 ----- //
	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

// ---- 이 함수는 테스트를 위해 메인 리스트를 임시 리스트에 저장하는데 사용됩니다 ----- //
func saveLists() {
	tmpArticleList = articleList
}

// ------ 이 함수는 임시 리스트에서 메인 리스트를 복구하는데 사용됩니다 ------ //
func restoreLists() {
	articleList = tmpArticleList
}
