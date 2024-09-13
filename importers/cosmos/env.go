package cosmos

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"buf.build/gen/go/dgroux/newsteam/protocolbuffers/go/admin"
	v1 "buf.build/gen/go/dgroux/newsteam/protocolbuffers/go/v1"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

/*
 * GetEnv
 */
func (s *Importer) GetEnv() (*v1.GetEnvResponse, error) {

	ret := &v1.GetEnvResponse{}

	env, err := getEnvironment()

	if err != nil {
		return nil, err
	}

	for _, pub := range env.Publications {

		toAppend := &v1.Publication{
			Name: pub.Name,
		}

		if pub.Meta.Description != "" {
			toAppend.Description = &pub.Meta.Description
		}

		for _, section := range pub.Sections {
			toAppend.Menu = append(toAppend.Menu, &v1.Publication_MenuItem{
				Id:    section.ID,
				Title: &section.Name,
				Type: &v1.Publication_MenuItem_Page_{
					Page: &v1.Publication_MenuItem_Page{},
				},
			})
		}

		ret.Publications = append(ret.Publications, toAppend)

	}

	return ret, nil
}

/*
 * updateFeeds
 */
func updateFeeds() {

	/*
	 * Update feeds
	 */
	env, err := getEnvironment()

	if err != nil {
		log.Fatal(errors.Wrap(err, "could not get cosmos environment"))
	}

	for _, pub := range env.Publications {

		_ = &admin.Feed{
			Id:   pub.ID,
			Name: pub.Name,
		}
	}
}

/*
 * TEMP
 */
func getWire(feed *admin.Feed) *admin.Feed_Wire {
	if feed.Id == "bd" {
		return &admin.Feed_Wire{
			Active:          true,
			AddToList:       true,
			AutoPublish:     true,
			Premium:         true,
			ProcessorUrl:    "http://localhost:3333",
			UpdateFrequency: 60,
		}
	}

	return nil
}

/*
 * createSections
 */
func createSections(feed *admin.Feed, pub Publication) {

	for range pub.Sections {

		// if section.ID == "" {
		// 	continue
		// }

		// ss, err := client.Section.Create(&admin.CreateSectionRequest{
		// 	FeedId: feed.Id,
		// 	Section: &admin.SectionInput{
		// 		Name:  &section.ID,
		// 		Title: &section.Name,
		// 		Sid:   getSid(pub.ID, section.ID),
		// 	},
		// })

		// if err != nil {
		// 	log.Fatal(errors.Wrap(err, "could not create section"))
		// }

		// for _, subsection := range section.Sections {

		// 	if subsection.ID == "" {
		// 		continue
		// 	}

		// 	_, err := client.Section.Create(&admin.CreateSectionRequest{
		// 		FeedId: feed.Id,
		// 		Section: &admin.SectionInput{
		// 			Name:     &subsection.ID,
		// 			Title:    &subsection.Name,
		// 			ParentId: &ss.Id,
		// 			Sid:      getSid(pub.ID, section.ID, subsection.ID),
		// 		},
		// 	})

		// 	if err != nil {
		// 		log.Fatal(errors.Wrap(err, "could not create subsection"))
		// 	}
		// }
	}
}

/*
 * getEnvironment
 */
func getEnvironment() (*Env, error) {

	fmt.Println("Getting environment from Cosmos...")

	client := resty.New()

	response, err := client.R().
		SetCookie(&http.Cookie{
			Name:  "_cosmos_auth",
			Value: "348c758db9458109244ddbefe4549bde73133324",
		}).
		Get("https://businesslive.co.za/apiv1/config/env")

	if err != nil {
		return nil, errors.Wrap(err, "could not get articles from cosmos")
	}

	dst := &Env{}

	json.Unmarshal(response.Body(), dst)

	return dst, nil

	// response, err := netclient.Json[Env](
	// 	netclient.Options{
	// 		Path:           "https://businesslive.co.za/apiv1/config/env",
	// 		AuthCookieName: "_cosmos_auth",
	// 		AccessToken:    "348c758db9458109244ddbefe4549bde73133324",
	// 	},
	// )

}

/*
 * Types
 */

// type Env struct {
// 	Publication  Publication
// 	Publications []Publication
// 	Sections     []struct {
// 		ID       string
// 		Name     string
// 		URLKey   string
// 		Sections []struct {
// 			ID     string
// 			Name   string
// 			URLKey string
// 		}
// 	}
// }

// type Publication struct {
// 	ID            string
// 	Name          string
// 	Primary       bool
// 	Routed        bool
// 	Show          bool
// 	URLKey        string
// 	UsePrimaryNav bool
// 	Sections      []struct {
// 		ID       string
// 		Name     string
// 		URLKey   string
// 		Sections []struct {
// 			ID     string
// 			Name   string
// 			URLKey string
// 		}
// 	}
// 	Meta struct {
// 		Description string
// 		Keywords    string
// 	}
// }
