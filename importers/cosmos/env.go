package cosmos

import (
	"fmt"
	"log"

	"buf.build/gen/go/dgroux/newsteam/protocolbuffers/go/admin"
	v1 "buf.build/gen/go/dgroux/newsteam/protocolbuffers/go/v1"
	"github.com/feight/newsteam-sdk/lib"
	"github.com/pkg/errors"
)

/*
 * GetEnv
 */
func (s *Importer) GetEnv() (*v1.GetEnvResponse, error) {

	ret := &v1.GetEnvResponse{
		Capabilities: &v1.WireCapabilities{
			Article: true,
		},
	}

	// 	env, err := getEnvironment(s)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	publication := &v1.Publication{
	// 		Id:          env.Publication.ID,
	// 		Name:        env.Publication.Name,
	// 		Description: &env.Publication.Meta.Description,
	// 	}
	//
	// 	for _, p := range env.Publications {
	// 		publication.Menu = append(publication.Menu, &v1.Publication_MenuItem{
	// 			Id:       p.ID,
	// 			Text:     &p.Name,
	// 			Title:    &p.Name,
	// 			Children: appendSections(p.Sections),
	// 			Type: &v1.Publication_MenuItem_Page_{
	// 				Page: &v1.Publication_MenuItem_Page{},
	// 			},
	// 		})
	// 	}
	//
	// 	ret.Publications = append(ret.Publications, publication)

	return ret, nil
}

func getPlacement(s *Importer, publication, section, subsection string) *admin.ArticlePlacement {

	env, err := getEnvironment(s)

	if err != nil {
		panic(err)
	}

	new := func(id, name string) *admin.ArticlePlacement_SourceDescriptor {
		return &admin.ArticlePlacement_SourceDescriptor{
			Id:   id,
			Name: name,
		}
	}

	for _, pub := range env.Publications {
		if pub.ID == publication {
			ret := &admin.ArticlePlacement{
				Source: []*admin.ArticlePlacement_SourceDescriptor{new(pub.ID, pub.Name)}}

			for _, s := range pub.Sections {
				if s.ID == section {
					ret.Source = append(ret.Source, new(s.ID, s.Name))

					for _, ss := range s.Sections {
						if ss.ID == subsection {
							ret.Source = append(ret.Source, new(ss.ID, ss.Name))
						}
					}
				}
			}
			return ret
		}
	}

	return nil
}

// func appendSections(sections []Section) []*v1.Publication_MenuItem {
//
// 	mi := []*v1.Publication_MenuItem{}
// 	for _, section := range sections {
// 		mi = append(mi, &v1.Publication_MenuItem{
// 			Id:       section.ID,
// 			Text:     &section.Name,
// 			Title:    &section.Name,
// 			Children: appendSections(section.Sections),
// 			Type: &v1.Publication_MenuItem_Page_{
// 				Page: &v1.Publication_MenuItem_Page{},
// 			},
// 		})
// 	}
//
// 	return mi
// }

/*
 * updateBuckets
 */
func updateBuckets() {

	/*
	 * Update buckets
	 */
	env, err := getEnvironment(nil)

	if err != nil {
		log.Fatal(errors.Wrap(err, "could not get cosmos environment"))
	}

	for _, pub := range env.Publications {

		_ = &admin.Bucket{
			Id:   pub.ID,
			Name: pub.Name,
		}
	}
}

/*
 * TEMP
 */
func getWire(bucket *admin.Bucket) *admin.Wire {
	if bucket.Id == "bd" {
		return &admin.Wire{
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
func createSections(bucket *admin.Bucket, pub Publication) {

	for range pub.Sections {

		// if section.ID == "" {
		// 	continue
		// }

		// ss, err := client.Section.Create(&admin.CreateSectionRequest{
		// 	BucketId: bucket.Id,
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
		// 		BucketId: bucket.Id,
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

var environ *Env

/*
 * getEnvironment
 */
func getEnvironment(s *Importer) (*Env, error) {

	if environ != nil {
		return environ, nil
	}

	fmt.Println("Getting environment from Cosmos...")

	env, err := lib.Json[Env](fmt.Sprintf("%s/config/env", s.Host), s.AccessToken)

	if err != nil {
		return nil, errors.Wrap(err, "could not get articles from cosmos")
	}

	environ = env

	return env, nil
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
