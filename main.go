package main

import (
	"kua-mei-a-api/ptt"
	"log"
)

func startDailyBeautyCrawler() {
	log.Println("getting daily beauty...")
	// TODO: do parallelly
	beauties, err := ptt.FetchBeauties()
	if err != nil {
		panic(err)
	}

	for _, data := range beauties {
		if data.NVote > 0 && data.NImage > 0 {
			log.Println(data)
		}
	}

	log.Println("Finish")
}

func main() {
}
