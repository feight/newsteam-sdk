package main

import (
	"github.com/feight/newsteam-sdk"
	"github.com/feight/newsteam-sdk/importers/cosmos"
)

func main() {

	newsteam.InitializeBuckets([]newsteam.Bucket{
		&cosmos.Importer{Host: "https://cosmos-dot-iol-prod.appspot.com/apiv1", Bucket: "news", AccessToken: "ba9b292dd015069f5a941921fcb313c6dde7de2c"},
	})
}
