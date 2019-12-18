package main

import (
	"kua-mei-a-api/db"
	"kua-mei-a-api/model"
	"kua-mei-a-api/ptt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func startDailyBeautyCrawler() {
	log.Println("getting daily beauty...")
	// TODO: do parallelly
	beauties, err := ptt.FetchBeauties()
	if err != nil {
		panic(err)
	}

	imageURLs := []*model.ImageData{}
	for _, data := range beauties {
		if data.NVote > 0 && data.NImage > 0 {
			log.Println(data)
			for _, url := range data.Images {
				img := model.ImageData{
					URL: url,
				}
				imageURLs = append(imageURLs, &img)
			}
		}
	}
	db.SessionInsert(imageURLs)
	log.Println("Finish")
}
func crawlerHandler(c *gin.Context) {
	startDailyBeautyCrawler()
	c.String(200, "Crawler successfully")
}
func main() {
	r := gin.Default()
	r.GET("/crawler", crawlerHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listen on port %s", port)
	err := r.Run(":" + port)
	panic(err)
}
