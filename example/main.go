package main

import (
	"fmt"
	"time"

	"github.com/KidiXDev/rich-go/client"
)

func main() {
	err := client.Login("DISCORD_APP_ID")
	if err != nil {
		panic(err)
	}

	now := time.Now()
	activityType := client.ActivityTypeListening
	statusDisplayType := client.StatusDisplayTypeDetails
	err = client.SetActivity(client.Activity{
		Type:              activityType,
		StatusDisplayType: &statusDisplayType,
		Name:              "Never Gonna Give You Up",
		Details:           "Rick Astley",
		State:             "Whenever You Need Somebody",
		LargeImage:        "album_cover",
		LargeText:         "Never Gonna Give You Up",
		SmallImage:        "music",
		SmallText:         "Listening to...",
		Timestamps: &client.Timestamps{
			Start: &now,
		},
		Buttons: []*client.Button{
			{
				Label: "GitHub",
				Url:   "https://github.com/KidiXDev/rich-go",
			},
		},
	})

	if err != nil {
		panic(err)
	}

	// Discord will only show the presence if the app is running
	// Sleep for a few seconds to see the update
	fmt.Println("Sleeping...")
	time.Sleep(time.Second * 10)
}
