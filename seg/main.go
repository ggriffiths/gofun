package main

import (
	"fmt"
	"os"

	analytics "gopkg.in/segmentio/analytics-go.v3"
)

func main() {
	key := os.Getenv("SEGMENT_KEY")

	client := analytics.New(key)

	if err := client.Enqueue(analytics.Identify{
		UserId: "1234123",
		Traits: analytics.NewTraits().
			SetName("Tim Smith").
			SetEmail("tmsith@gmail.com").
			SetAddress("123 something street").
			Set("FavoriteColor", "blue"),
	}); err != nil {
		fmt.Println("Error occurred:", err)
	}

	if err := client.Close(); err != nil {
		fmt.Println("Error occurred:", err)
	}
}
