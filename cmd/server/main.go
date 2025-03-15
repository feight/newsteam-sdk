package main

import (
	"github.com/feight/newsteam-sdk"
	"github.com/feight/newsteam-sdk/importers/cosmos"
)

func main() {

	newsteam.InitializeBuckets([]newsteam.Bucket{
		&cosmos.Importer{Host: "https://www.businesslive.co.za/apiv1", Bucket: "news", AccessToken: "62a75039712d063ab27da709ecaf24c607acbd2a"},
	})
}
