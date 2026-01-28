//------------------- handlers.article_test.go -----------------//

package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ----- 홈페이지에 대한 GET 요청이 인증되지 않은 사용자에 대해 HTTP 200 코드가 있는 홈페이지를 반환하는지 테스트합니다 ----- //
func TestShowIndexPageUnauthenticated(t *testing.T) {
	r := getRouter(true)

	r.GET("/", showIndexPage)

	// ---- 위의 라우터로 보낼 요청을 생성합니다 ----- //
	req, _ := http.NewRequest("GET", "/", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// ----- HTTP status 코드가 200인지 확인합니다 ----- //
		statusOK := w.Code == http.StatusOK

		// ------ 페이지의 제목이 'Home Page' 인지 테스트합니다 ------- //
		// --- HTML 페이지를 파싱하고 처리할 수 있는 라이브러리를 사용하여 훨씬 디테일한 테스트를 수행합니다 --- //
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Home Page</title>") > 0

		return statusOK && pageOK
	})
}
