package api

import (
	"errors"
	"fmt"
	"net/http"

	"kua-mei-a-api/model"

	"github.com/PuerkitoBio/goquery"
)

// FetchPageAmount get latest page number
func FetchPageAmount() (int, error) {
	url := "https://www.ptt.cc/bbs/Beauty/index.html"

	client := http.DefaultClient
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", "over18=1")

	res, _ := client.Do(req)
	doc, _ := goquery.NewDocumentFromResponse(res)

	prevPageSelector := "div.btn-group.btn-group-paging a:nth-child(2)"
	href, _ := doc.Find(prevPageSelector).Attr("href")

	var n int
	fmt.Sscanf(href, "/bbs/Beauty/index%d.html", &n)

	if n == 0 {
		return 0, errors.New("Cannot connect to PTT")
	}
	return n + 1, nil
}

// FetchPage get all posts in a page
func FetchPage(prefix string, page int) ([]model.Post, error) {
	baseURL := "https://www.ptt.cc/bbs/Beauty/"
	url := fmt.Sprintf("%sindex%d.html", baseURL, page)

	// TODO: refactor HTTP client
	client := http.DefaultClient
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", "over18=1")

	res, _ := client.Do(req)
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	posts := parseDoc2Posts(doc, prefix)
	return posts, nil
}

// Search use PTT search to get search result
// sometimes PTT cache search result
func Search(prefix string, page, recommend int) ([]model.Post, error) {
	// page from 1, 2, ...
	baseURL := "https://www.ptt.cc/bbs/Beauty/search"
	url := fmt.Sprintf("%s?page=%d&q=%s+recommend:%d", baseURL, page, prefix, recommend)

	client := http.DefaultClient
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", "over18=1")

	res, _ := client.Do(req)
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	posts := parseDoc2Posts(doc, prefix)
	return posts, nil
}
