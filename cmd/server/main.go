package main

import (
	"github.com/feight/newsteam-sdk"
	"github.com/feight/newsteam-sdk/importers/cosmos"
)

func main() {

	newsteam.InitializeBuckets([]newsteam.Bucket{
		&cosmos.Importer{Host: "https://cosmos-dot-iol-prod.appspot.com/apiv1", Bucket: "news"},
	})
}
