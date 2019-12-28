package main

import (
	"kua-mei-a-api/db"
	"kua-mei-a-api/model"
	"kua-mei-a-api/ptt"
	"log"
	"math/rand"
	"net/http"
	"os"
)

func startDailyBeautyCrawler() {
	log.Println("getting daily beauty...")
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

func getRandomImageURL() (string, error) {
	count, err := db.Count()
	if err != nil {
		return "", err
	}
	id := rand.Int63n(count) + 1
	img := db.SelectOne(id)
	return img.URL, nil
}

func getImageHandler(w http.ResponseWriter, r *http.Request) {
	url, err := getRandomImageURL()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte(url))
}

func showImageHandler(w http.ResponseWriter, r *http.Request) {
	url, err := getRandomImageURL()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte(`<html>
					<head><title>kua mei a</title></head>
					<style type="text/css">
					body {
						margin: 0;
						padding: 0;
						display: flex;
						justify-content: center;
						align-items: center;
						height: 100vh;
						width: 100%;
					}
					img {
						height: 100vh;
						width: auto;
					}
					</style>
					<body>
					<img src="` + url + `"/>
					</body>
					</html>`))
}

func crawlerHandler(w http.ResponseWriter, r *http.Request) {
	startDailyBeautyCrawler()
	w.Write([]byte(`Crawler successfully`))
}

func main() {
	http.HandleFunc("/crawler", crawlerHandler)
	http.HandleFunc("/image", getImageHandler)
	http.HandleFunc("/", showImageHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listen on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
