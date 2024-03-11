package cosmos

import (
	"encoding/json"
	"fmt"

	"buf.build/gen/go/dgroux/newsteam/protocolbuffers/go/admin"
	"github.com/feight/newsteam-sdk"
	"github.com/feight/newsteam-sdk/lib/client"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type CosmosImporter struct {
	Host    string
	Project string
}

func (s *CosmosImporter) ProjectId() string {
	return s.Project
}

/*
 * GetLogfiles
 */
func (s *CosmosImporter) GetLogfiles() ([][]byte, error) {

	fmt.Println("getting logfiles...")

	type ResponseT = []map[string]any

	type QueryParams struct {
		Limit   int32  `json:"limit"`
		Section string `json:"section"`
	}

	response, err := client.Json[ResponseT](
		client.Options{
			Path:           s.Host + "/pub/articles/get-all",
			AuthCookieName: "_cosmos_auth",
			AccessToken:    "", // TODO

			Data: QueryParams{
				Limit:   100,
				Section: "lifestyle",
			},
		},
	)

	if err != nil {

		return nil, errors.Wrap(err, "could not get articles from cosmos")

	}

	ret := [][]byte{}

	for _, a := range *response {

		record, _ := json.Marshal(a)

		ret = append(ret, record)

	}

	return ret, nil
}

/*
 * ProcessLogfile
 */
func (s *CosmosImporter) ProcessLogfile(content []byte) []*admin.ArticleInput {

	a := Article{}

	err := json.Unmarshal(content, &a)

	if err != nil {

		panic(errors.Wrap(err, "could not unmarshal logfile"))

	}

	return []*admin.ArticleInput{s.createArticle(a)}
}

/*
 * createArticle
 */
func (s *CosmosImporter) createArticle(ca Article) *admin.ArticleInput {

	m := &admin.ArticleInput{}
	m.AdTagCustom = &ca.AdTagCustom
	m.ContentType = &ca.ContentType
	m.Assigned = ca.Assigned
	// m.AuthorIds = ca.AuthorKeys
	m.Authors = ca.Authors
	m.Breaking = &ca.Breaking
	m.CanonicalUrl = &ca.CanonicalURL
	m.CommentsEnabled = &ca.CommentsEnabled
	m.EditorsChoice = &ca.EditorsChoice
	m.ExternalUrl = &ca.ExternalURL
	m.Groups = ca.Groups
	m.HideInApp = &ca.HideInApp
	// m.ImageHeader
	// m.ImageThumbnail
	// m.Source = &ca.Key // TODO: Create source ID
	m.Keywords = ca.Keywords
	// m.LocationGeo = &ca.LocationGeo
	// m.LocationKeywords = ca.LocationKeywords
	// m.LocationLat = ca.LocationLat
	// m.LocationName
	m.Native = &ca.Native
	// m.Notes
	m.OnHold = &ca.OnHold
	// m.PlainText = ca.PlainText
	m.Priority = &ca.Priority
	// m.Published
	m.PushNotify = &ca.PushNotify
	// m.RelatedArticles = ca.RelatedArticles
	// m.RelatedArticlesKeys
	m.Sensitive = &ca.Sensitive
	m.SourceId = proto.String(fmt.Sprint(ca.Key))
	// m.ShareID
	m.ShowCreatorCard = &ca.ShowAuthorCard
	// m.Slug
	// m.SlugCustom
	// m.SlugOld
	m.SocialShareText = &ca.SocialShareText
	m.Source = &ca.Source
	// m.SourceID
	// m.Sponsor
	// m.SponsorIds = ca.SponsorKeys
	// m.Status = &ca.Status
	m.Style = &ca.Style
	// m.Summary
	m.Syndicate = &ca.Syndicate
	// m.SyndicateStatus
	// m.Synopsis
	m.SynopsisCustom = &ca.SynopsisCustom
	m.Tags = ca.Tags
	m.TitleCustom = &ca.TitleCustom
	// m.TitleListing
	// m.TitleListingText
	// m.TitleSectionText
	// m.Weight

	m.Title = &admin.TextProperty{
		Raw:  ca.Title,
		Html: ca.Title,
		Text: ca.TitleText,
	}
	m.Title2 = &admin.TextProperty{
		Raw:  ca.Title2,
		Html: ca.Title2,
		Text: ca.Title2Text,
	}
	m.Title3 = &admin.TextProperty{
		Raw:  ca.Title3,
		Html: ca.Title3,
		Text: ca.Title3Text,
	}
	m.Intro = &admin.TextProperty{
		Raw:  ca.Intro,
		Html: ca.Intro,
		Text: ca.IntroText,
	}

	m.Labels = map[string]string{"source_id": "cosmos"}

	// TODO: fix:
	//
	// for _, placement := range ca.Sections {

	// 	// TODO: Fix, this is slow
	// 	section, err := s.client.Section.GetBySID(&admin.GetSectionBySIDRequest{
	// 		ProjectId: placement.Publication,
	// 		Sid:       *s.getSid(placement.Publication, placement.Section, placement.Subsection),
	// 	})

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	m.Placements = append(m.Placements, section.Id)
	// }

	for _, w := range ca.Widgets {
		s.mapWidget(m, w.Type, w.Data)
	}

	return m
}

/*
 * mapWidget
 */
func (s *CosmosImporter) mapWidget(m *admin.ArticleInput, typ string, data map[string]any) {

	switch typ {

	case "text":
		m.Widgets = append(m.Widgets, &admin.WidgetCollection{
			Widget: &admin.WidgetCollection_Text{
				Text: mapTextWidget(typ, getWidget[TextWidget](data))}})
	case "image":
		m.Widgets = append(m.Widgets, &admin.WidgetCollection{
			Widget: &admin.WidgetCollection_Image{
				Image: mapImageWidget(typ, getWidget[ImageWidget](data))}})

		// case "accordion":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapAccordionWidget(typ, getWidget[TextWidget](data))}})
		// case "article-list":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapArticleListWidget(typ, getWidget[TextWidget](data))}})
		// case "chartblocks":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapChartblocksWidget(typ, getWidget[TextWidget](data))}})
		// case "crowdsignal":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapCrowdsignalWidget(typ, getWidget[TextWidget](data))}})
		// case "facebook-page":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapFacebookPageWidget(typ, getWidget[TextWidget](data))}})
		// case "facebook-post":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapFacebookPostWidget(typ, getWidget[TextWidget](data))}})
		// case "facebook-video":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapFacebookVideoWidget(typ, getWidget[TextWidget](data))}})
		// case "image-gallery":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapImageGalleryWidget(typ, getWidget[TextWidget](data))}})
		// case "giphy":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapGiphyWidget(typ, getWidget[TextWidget](data))}})
		// case "google-map":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapGoogleMapWidget(typ, getWidget[TextWidget](data))}})
		// case "horizontal-line":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapHorizontalLineWidget(typ, getWidget[TextWidget](data))}})
		// case "html":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapHtmlWidget(typ, getWidget[TextWidget](data))}})
		// case "iframely":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapIframelyWidget(typ, getWidget[TextWidget](data))}})
		// case "infogram":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapInfogramWidget(typ, getWidget[TextWidget](data))}})
		// case "instagram":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapInstagramWidget(typ, getWidget[TextWidget](data))}})
		// case "iono":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapIonoWidget(typ, getWidget[TextWidget](data))}})
		// case "issuu":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapIssuuWidget(typ, getWidget[TextWidget](data))}})
		// case "jwplayer":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapJwplayerWidget(typ, getWidget[TextWidget](data))}})
		// case "kickstarter":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapKickstarterWidget(typ, getWidget[TextWidget](data))}})
		// case "link-list":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapLinkListWidget(typ, getWidget[TextWidget](data))}})
		// case "oovvuu":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapOovvuuWidget(typ, getWidget[TextWidget](data))}})
		// case "quote":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapQuoteWidget(typ, getWidget[TextWidget](data))}})
		// case "related-articles":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapRelatedArticlesWidget(typ, getWidget[TextWidget](data))}})
		// case "scribd":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapScribdWidget(typ, getWidget[TextWidget](data))}})
		// case "soundcloud":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapSoundcloudWidget(typ, getWidget[TextWidget](data))}})
		// case "text-block":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapTextBlockWidget(typ, getWidget[TextWidget](data))}})
		// case "tiktok":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapTiktokWidget(typ, getWidget[TextWidget](data))}})
		// case "twitter":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapTwitterWidget(typ, getWidget[TextWidget](data))}})
		// case "youtube":
		// 	m.Widgets = append(m.Widgets, &content.WidgetCollection{
		// 		Widget: &content.WidgetCollection_Text{
		// 			Text: mapYoutubeWidget(typ, getWidget[TextWidget](data))}})

	}
}

func getWidget[T any](data map[string]any) T {

	var widget T

	mapstructure.Decode(data, &widget)

	return widget

}

func mapTextWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapImageWidget(t string, w ImageWidget) *admin.ImageWidget {

	// Upload image...
	image := newsteam.UploadImageFromUrl(fmt.Sprintf("https:%s=s2000", w.Image.Filepath))
	// Id:          w.String("text"),
	image.ContentType = w.Image.ContentType
	image.Width = w.Image.Width
	image.Height = w.Image.Width
	// image.// Filesize=    w.Image.Width
	image.Keywords = w.Image.Keywords
	image.Palette = w.Image.Palette
	image.Creator = w.Image.Author
	image.FileName = w.Image.Filename
	image.Blur = w.Image.Blur
	image.Description = w.Image.Description
	image.Title = w.Image.Title
	image.FocalY = w.Image.Width
	image.FocalX = w.Image.Width
	image.Average = w.Image.Average

	return &admin.ImageWidget{
		Type: t,
		Data: &admin.ImageWidget_Data{
			Image: image,
		},
	}
}

func mapAccordionWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapArticleListWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapChartblocksWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapCrowdsignalWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapFacebookPageWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapFacebookPostWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapFacebookVideoWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapImageGalleryWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapGiphyWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapGoogleMapWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapHorizontalLineWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapHtmlWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapIframelyWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapInfogramWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapInstagramWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapIonoWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapIssuuWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapJwplayerWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapKickstarterWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapLinkListWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapOovvuuWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapQuoteWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapRelatedArticlesWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapScribdWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapSoundcloudWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapTextBlockWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapTiktokWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapTwitterWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}

func mapYoutubeWidget(t string, w TextWidget) *admin.TextWidget {
	return &admin.TextWidget{
		Type: t,
		Data: &admin.TextWidget_Data{
			Text: w.Text,
			Html: w.Clean,
			Raw:  w.HTML,
			// Clear: widget.Type,
		},
	}
}
