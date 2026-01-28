package main

import "testing"

// ------ 모든 기사 목록을 가져오는 기능을 테스트합니다 ------- //
func TestGetAllArticles(t *testing.T) {
	alist := getAllArticles()

	// --- 반환된 기사 목록의 길이가, 리스트가 포함하고 있는 전역변수의 길이와 같은지 확인합니다 --- //
	if len(alist) != len(articleList) {
		t.Fail()
	}

	// ---- 3가지 필드값이 동일한지 확인해봅니다 ----- //
	for i, v := range alist {
		if v.Content != articleList[i].Content ||
			v.ID != articleList[i].ID ||
			v.Title != articleList[i].Title {

			t.Fail()
			break
		}
	}
}
