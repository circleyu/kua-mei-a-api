package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"kua-mei-a-api/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/vjeantet/jodaTime"
)

// parseNVote parses vote text to int
// "50" => 50, "爆" => 100
// "" => 0
// "X7" => -1
// there is no need to handle nVote <= 0
// because they are filterer out when searching
func parseNVote(nVoteText string) int {
	if nVoteText == "爆" {
		return 100
	}
	if nVoteText == "" {
		return 0
	}
	if strings.HasPrefix(nVoteText, "X") {
		return -1
	}
	nVote, _ := strconv.Atoi(nVoteText)
	return nVote
}

func parseDoc2Posts(doc *goquery.Document, prefix string) []model.Post {
	// TODO: remove 置頂文
	posts := make([]model.Post, 0, 20)
	doc.Find(".r-ent").Each(func(i int, el *goquery.Selection) {
		nVoteText := el.Find(".hl").Text()
		nVote := parseNVote(nVoteText)

		titleEl := el.Find(".title > a")
		title := titleEl.Text()

		if !strings.HasPrefix(title, prefix) {
			return
		}

		hrefText, _ := titleEl.Attr("href")
		href := "https://www.ptt.cc" + hrefText

		currentYear := time.Now().Year()
		mmdd := strings.TrimSpace(el.Find(".meta .date").Text())
		dateText := fmt.Sprintf("%d/%s", currentYear, mmdd)
		date, _ := jodaTime.ParseInLocation("YYYY/M/dd", dateText, "Asia/Taipei")

		p := model.Post{
			Title: title,
			Href:  href,
			NVote: nVote,
			Date:  date,
		}

		posts = append(posts, p)
	})
	return posts
}
