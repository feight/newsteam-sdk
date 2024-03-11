package main

import (
	"github.com/feight/newsteam-sdk"
	"github.com/feight/newsteam-sdk/wires/cosmos"
)

func main() {
	newsteam.InitializeFeeds([]newsteam.Feed{
		&cosmos.CosmosImporter{
			Project: "bd", Host: "https://businesslive.co.za/apiv1",
		},
	})
}
