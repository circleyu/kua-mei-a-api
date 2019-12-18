package ptt

import (
	"kua-mei-a-api/model"
	"kua-mei-a-api/ptt/api"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// TODO: split ptt api layer and utils layer
func init() {
	rand.Seed(time.Now().UnixNano())
}

func fetchYesterdayPosts() ([]model.Post, error) {
	prefix := "[正妹]"
	recentPosts := make([]model.Post, 0, 20)

	// get recent posts
	page, err := api.FetchPageAmount()
	if err != nil {
		return nil, err
	}

	for ; ; page-- {
		posts, err := api.FetchPage(prefix, page)

		if err != nil {
			return nil, err
		}

		if len(posts) == 0 {
			continue
		}

		recentPosts = append(recentPosts, posts...)
		oldestDate := recentPosts[len(recentPosts)-1].Date
		if isBeforeYesterday(oldestDate) {
			break
		}
	}

	// filter yesterday post
	yesterdayPosts := make([]model.Post, 0, 10)
	for _, p := range recentPosts {
		if isYesterday(p.Date) {
			yesterdayPosts = append(yesterdayPosts, p)
		}
	}

	return yesterdayPosts, nil
}

// FetchRandomBeauty randomly fetch a model.Beauty
func FetchRandomBeauty() (model.Beauty, error) {
	prefix := "[正妹]"
	page := rand.Intn(40) + 11 // 11 ~ 50
	posts, err := api.Search(prefix, page, 99)

	if err != nil {
		return model.Beauty{}, err
	}

	idx := rand.Intn(len(posts)) // 0 ~ len(posts)-1
	p := posts[idx]
	b := p.ToBeauty()
	return b, nil
}

func getBestBeauties(posts []model.Post) []model.Beauty {
	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].NVote > posts[j].NVote
	})

	nBeauty := len(posts)
	//champions := posts[:nBeauty]
	beauties := make([]model.Beauty, nBeauty)

	var wg sync.WaitGroup
	wg.Add(nBeauty)
	for i, p := range posts {
		go func(i int, p model.Post) {
			beauties[i] = p.ToBeauty()
			wg.Done()
		}(i, p)
	}
	wg.Wait()

	return beauties
}

// FetchBeauties send a request to get beauties from getDailyBeauties api
func FetchBeauties() ([]model.Beauty, error) {
	posts, err := fetchYesterdayPosts()
	if err != nil {
		return nil, err
	}
	beauties := getBestBeauties(posts)
	return beauties, nil
}
