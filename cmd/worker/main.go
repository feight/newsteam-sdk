package main

import (
	"github.com/feight/newsteam-sdk"
	"github.com/feight/newsteam-sdk/wires"
)

func main() {
	newsteam.InitializeFeeds([]newsteam.Feed{
		&wires.CosmosImporter{
			Project:     "bd",
			Host:        "https://businesslive.co.za/apiv1",
			AccessToken: "",
		},
	})
}
