package cosmos

import (
	"fmt"
	"log"

	"buf.build/gen/go/dgroux/newsteam/protocolbuffers/go/admin"
	"github.com/pkg/errors"
)

const organizationId = "newsteam"

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

		/*
		 * Get the feed
		 */
		// _, err := client.Feed.Get(&admin.GetFeedRequest{
		// 	Id: feed.Id,
		// })

		// /*
		//  * Create the feed if it does not exist
		//  */
		// if err != nil {

		// 	fmt.Println("Creating feed", feed.Name)

		// 	_, err := client.Feed.Create(&admin.CreateFeedRequest{
		// 		Id:             feed.Id,
		// 		OrganizationId: organizationId,
		// 		Feed: &admin.FeedInput{
		// 			Name: &feed.Name,
		// 			Wire: getWire(feed),
		// 		},
		// 	})

		// 	if err != nil {
		// 		log.Fatal(errors.Wrap(err, "could not create feed"))
		// 	}

		// 	/*
		// 	 * Create the sections
		// 	 */
		// 	createSections(feed, pub)
		// }
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

	// response, err := netclient.Json[Env](
	// 	netclient.Options{
	// 		Path:           "https://businesslive.co.za/apiv1/config/env",
	// 		AuthCookieName: "_cosmos_auth",
	// 		AccessToken:    "348c758db9458109244ddbefe4549bde73133324",
	// 	},
	// )

	// if err != nil {
	// 	return nil, errors.Wrap(err, "could not get articles from cosmos")
	// }

	// return response, nil
	return nil, nil
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
