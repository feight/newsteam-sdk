package main

import (
	"github.com/feight/newsteam-sdk"
	"github.com/feight/newsteam-sdk/importers/cosmos"
)

func main() {
	newsteam.InitializeFeeds([]newsteam.Feed{
		&cosmos.Importer{Project: "bl", Host: "https://businesslive.co.za/apiv1", AccessToken: "348c758db9458109244ddbefe4549bde73133324"},
		&cosmos.Importer{Project: "bd", Host: "https://businesslive.co.za/apiv1", AccessToken: "348c758db9458109244ddbefe4549bde73133324"},
	})
}
