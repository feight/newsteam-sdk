package cosmos

import (
	"encoding/json"
	"fmt"
	"time"

	"buf.build/gen/go/dgroux/newsteam/protocolbuffers/go/admin"
	"github.com/feight/newsteam-sdk"
	"github.com/feight/newsteam-sdk/lib"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Importer struct {
	Host        string
	Bucket      string
	AccessToken string
}

func (s *Importer) Id() string {
	return s.Bucket
}

/*
 * GetLogfiles
 */
func (s *Importer) GetLogfiles(state *admin.Cursor) ([][]byte, error) {

	var (
		offset       = state.SeekPos
		limit  int64 = 100
	)

	fmt.Println("getting logfiles... offset", offset)

	url := fmt.Sprintf("%s/pub/articles/get-all?limit=%d&offset=%d", s.Host, limit, offset)

	dst, err := lib.Json[[]map[string]any](url, s.AccessToken)

	if err != nil {
		return nil, errors.Wrapf(err, "could not get articles from cosmos: %s", url)
	}

	ret := [][]byte{}

	for _, a := range *dst {

		entry, _ := json.Marshal(a)

		ret = append(ret, entry)
	}

	state.SeekPos = offset + limit
	state.SeekDate = 0 // TODO: Set latest article publish date

	return ret, nil
}

/*
 * ProcessLogfile
 */
func (s *Importer) ProcessLogfile(bucket *admin.Bucket, content []byte) []*admin.Article {

	a := Article{}

	err := json.Unmarshal(content, &a)

	if err != nil {
		panic(errors.Wrap(err, "could not unmarshal logfile"))
	}

	article := s.createArticle(bucket, a)

	return []*admin.Article{article}
}

/*
 * createArticle
 */
func (s *Importer) createArticle(bucket *admin.Bucket, ca Article) *admin.Article {

	m := &admin.Article{}
	// m.AuthorIds = ca.AuthorKeys
	m.Authors = ca.Authors
	m.Breaking = &ca.Breaking
	m.CanonicalUrl = &ca.CanonicalURL
	m.RedirectUrl = &ca.ExternalURL
	m.Groups = ca.Groups
	// m.ImageHeader
	// m.ImageThumbnail
	// m.Source = &ca.Key // TODO: Create source ID
	m.Keywords = ca.Keywords
	// m.LocationGeo = &ca.LocationGeo
	// m.LocationKeywords = ca.LocationKeywords
	// m.LocationLat = ca.LocationLat
	// m.LocationName
	m.Sponsored = &ca.Native
	// m.Notes
	// m.PlainText = ca.PlainText
	// m.Published
	// m.RelatedArticles = ca.RelatedArticles
	// m.RelatedArticlesKeys
	m.Sensitive = &ca.Sensitive
	m.Source = &admin.SourceProperty{
		Id: fmt.Sprint(ca.Key),
	}
	// m.ShareID
	// m.Slug
	// m.SlugCustom
	// m.SlugOld
	m.ShareText = &admin.ShareTextProperty{
		Text: ca.SocialShareText,
	}
	// m.SourceID
	// m.Sponsor
	// m.SponsorIds = ca.SponsorKeys
	// m.Status = &ca.Status
	// m.Summary
	// m.SyndicateStatus
	// m.Synopsis
	m.SynopsisCustom = &ca.SynopsisCustom
	m.Tags = ca.Tags
	m.TitleListing = &admin.TextProperty{
		Html: ca.TitleCustom,
		Raw:  ca.TitleCustom,
		Text: ca.TitleCustom,
	}
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

	m.Attributes = map[string]string{"source_id": "cosmos"}

	for _, section := range ca.Sections {
		if placement := getPlacement(ca, section); s != nil {
			m.Placements = append(m.Placements, placement)
		}
	}

	for _, w := range ca.Widgets {
		s.mapWidget(m, w.Type, w.Data)
	}

	return m
}

/*
 * mapWidget
 */
func (s *Importer) mapWidget(m *admin.Article, typ string, data map[string]any) {

	switch typ {

	case "text":
		m.Widgets = append(m.Widgets, &admin.Widget{
			Data: &admin.Widget_Text{
				Text: mapTextWidget(getWidget[TextWidget](data))}})
	case "image":
		m.Widgets = append(m.Widgets, &admin.Widget{
			Data: &admin.Widget_Image{
				Image: mapImageWidget(getWidget[ImageWidget](data))}})

		// case "accordion":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapAccordionWidget(typ, getWidget[TextWidget](data))}})
		// case "article-list":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapArticleListWidget(typ, getWidget[TextWidget](data))}})
		// case "chartblocks":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapChartblocksWidget(typ, getWidget[TextWidget](data))}})
		// case "crowdsignal":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapCrowdsignalWidget(typ, getWidget[TextWidget](data))}})
		// case "facebook-page":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapFacebookPageWidget(typ, getWidget[TextWidget](data))}})
		// case "facebook-post":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapFacebookPostWidget(typ, getWidget[TextWidget](data))}})
		// case "facebook-video":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapFacebookVideoWidget(typ, getWidget[TextWidget](data))}})
		// case "image-gallery":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapImageGalleryWidget(typ, getWidget[TextWidget](data))}})
		// case "giphy":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapGiphyWidget(typ, getWidget[TextWidget](data))}})
		// case "google-map":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapGoogleMapWidget(typ, getWidget[TextWidget](data))}})
		// case "horizontal-line":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapHorizontalLineWidget(typ, getWidget[TextWidget](data))}})
		// case "html":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapHtmlWidget(typ, getWidget[TextWidget](data))}})
		// case "iframely":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapIframelyWidget(typ, getWidget[TextWidget](data))}})
		// case "infogram":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapInfogramWidget(typ, getWidget[TextWidget](data))}})
		// case "instagram":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapInstagramWidget(typ, getWidget[TextWidget](data))}})
		// case "iono":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapIonoWidget(typ, getWidget[TextWidget](data))}})
		// case "issuu":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapIssuuWidget(typ, getWidget[TextWidget](data))}})
		// case "jwplayer":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapJwplayerWidget(typ, getWidget[TextWidget](data))}})
		// case "kickstarter":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapKickstarterWidget(typ, getWidget[TextWidget](data))}})
		// case "link-list":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapLinkListWidget(typ, getWidget[TextWidget](data))}})
		// case "oovvuu":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapOovvuuWidget(typ, getWidget[TextWidget](data))}})
		// case "quote":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapQuoteWidget(typ, getWidget[TextWidget](data))}})
		// case "related-articles":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapRelatedArticlesWidget(typ, getWidget[TextWidget](data))}})
		// case "scribd":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapScribdWidget(typ, getWidget[TextWidget](data))}})
		// case "soundcloud":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapSoundcloudWidget(typ, getWidget[TextWidget](data))}})
		// case "text-block":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapTextBlockWidget(typ, getWidget[TextWidget](data))}})
		// case "tiktok":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapTiktokWidget(typ, getWidget[TextWidget](data))}})
		// case "twitter":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapTwitterWidget(typ, getWidget[TextWidget](data))}})
		// case "youtube":
		// 	m.Widgets = append(m.Widgets, &content.Widget{
		// 		Widget: &content.Widget_Text{
		// 			Text: mapYoutubeWidget(typ, getWidget[TextWidget](data))}})

	}
}

func getWidget[T any](data map[string]any) T {

	var widget T

	mapstructure.Decode(data, &widget)

	return widget

}

func mapTextWidget(w TextWidget) *admin.TextWidgetData {
	return &admin.TextWidgetData{
		Text: w.Text,
		Html: w.Clean,
		Raw:  w.HTML,
		// Clear: widget.Type,
	}
}

func mapImageWidget(w ImageWidget) *admin.ImageWidgetData {

	// Upload image...
	image := newsteam.UploadImageFromUrl(fmt.Sprintf("%s/raw", w.Image.Filepath))
	// Id:          w.String("text"),
	image.ContentType = w.Image.ContentType
	image.Width = w.Image.Width
	image.Height = w.Image.Height
	// image.// Filesize=    w.Image.Width
	image.Keywords = w.Image.Keywords
	image.Palette = w.Image.Palette
	image.Credit = &w.Image.Author
	image.Sha256Hash = w.Image.Filename
	image.Caption = &w.Image.Description
	image.Title = w.Image.Title
	image.FocalY = w.Image.FocalY
	image.FocalX = w.Image.FocalX
	image.Average = w.Image.Average

	return &admin.ImageWidgetData{
		Image: image,
	}
}

/*
 * Types ------------------------------------------------------------------------------------
 */

type Env struct {
	AppID              string   `json:"app_id"`
	Clients            []string `json:"clients"`
	ContentTypeDefault string   `json:"content_type_default"`
	ContentTypes       []struct {
		Key  string `json:"key"`
		Text string `json:"text"`
	} `json:"content_types"`
	Firebase struct {
		APIKey            string `json:"apiKey"`
		AuthDomain        string `json:"authDomain"`
		DatabaseURL       string `json:"databaseURL"`
		MessagingSenderID string `json:"messagingSenderId"`
		ProjectID         string `json:"projectId"`
		StorageBucket     string `json:"storageBucket"`
	} `json:"firebase"`
	Logrocket struct {
		APIKey string `json:"apiKey"`
	} `json:"logrocket"`
	NovaFeeds []struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"nova_feeds"`
	Publication  EnvPublication   `json:"publication"`
	Publications []EnvPublication `json:"publications"`
	Sections     []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Sections []struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			URLKey string `json:"urlKey"`
		} `json:"sections,omitempty"`
		URLKey string `json:"urlKey"`
	} `json:"sections"`
}

type EnvPublication struct {
	ID   string `json:"id"`
	Meta struct {
		Description string `json:"description"`
		Keywords    string `json:"keywords"`
	} `json:"meta"`
	Name     string       `json:"name"`
	Primary  bool         `json:"primary"`
	Routed   bool         `json:"routed"`
	Sections []EnvSection `json:"sections"`
	Show     bool         `json:"show"`
	Strings  struct {
	} `json:"strings"`
	URLKey        string `json:"urlKey"`
	UsePrimaryNav bool   `json:"usePrimaryNav"`
}

type EnvSection struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Sections []EnvSection `json:"sections,omitempty"`
	URLKey   string       `json:"urlKey"`
}

type Article struct {
	Access      bool     `json:"access"`
	Active      bool     `json:"active"`
	AdTagCustom string   `json:"ad_tag_custom"`
	Assigned    []string `json:"assigned"`
	// Author                *Author     `json:"author"`
	AuthorKeys            []int64     `json:"author_keys"`
	Authors               []string    `json:"authors"`
	Breaking              bool        `json:"breaking"`
	CanonicalURL          string      `json:"canonical_url"`
	CommentsEnabled       bool        `json:"comments_enabled"`
	ContentType           string      `json:"content_type"`
	CountAudio            int         `json:"count_audio"`
	CountBlockquote       int         `json:"count_blockquote"`
	CountEmbeddedArticles int         `json:"count_embedded_articles"`
	CountImage            int         `json:"count_image"`
	CountInfographic      int         `json:"count_infographic"`
	CountSocial           int         `json:"count_social"`
	CountVideo            int         `json:"count_video"`
	CountWords            int         `json:"count_words"`
	Created               int64       `json:"created"`
	CreatedBy             string      `json:"created_by"`
	Domain                string      `json:"domain"`
	EditURL               string      `json:"edit_url"`
	EditorsChoice         bool        `json:"editors_choice"`
	EmbeddedArticlesList  []int64     `json:"embedded_articles_list"`
	ExternalURL           string      `json:"external_url"`
	Groups                []string    `json:"groups"`
	HideInApp             bool        `json:"hide_in_app"`
	Image                 Image       `json:"image"`
	ImageHeader           Image       `json:"image_header"`
	ImageThumbnail        Image       `json:"image_thumbnail"`
	Images                []Image     `json:"images"`
	Intro                 string      `json:"intro"`
	IntroText             string      `json:"intro_text"`
	Key                   int64       `json:"key"`
	Keywords              []string    `json:"keywords"`
	LocationGeo           string      `json:"location_geo"`
	LocationKeywords      []string    `json:"location_keywords"`
	LocationLat           string      `json:"location_lat"`
	LocationName          string      `json:"location_name"`
	Modified              int64       `json:"modified"`
	ModifiedUser          string      `json:"modified_user"`
	Native                bool        `json:"native"`
	Notes                 string      `json:"notes"`
	OnHold                bool        `json:"on_hold"`
	Parent                string      `json:"parent"`
	PlainText             string      `json:"plain_text"`
	Priority              int32       `json:"priority"`
	Private               bool        `json:"private"`
	PubURL                string      `json:"pub_url"`
	Published             int64       `json:"published"`
	PushNotify            bool        `json:"push_notify"`
	ReadDuration          int         `json:"read_duration"`
	RelatedArticles       interface{} `json:"related_articles"`
	RelatedArticlesKeys   []int       `json:"related_articles_keys"`
	Sensitive             bool        `json:"sensitive"`
	ShareID               string      `json:"share_id"`
	ShowAuthorCard        bool        `json:"show_author_card"`
	Slug                  string      `json:"slug"`
	SlugCustom            string      `json:"slug_custom"`
	SlugOld               []string    `json:"slug_old"`
	SocialShareText       string      `json:"social_share_text"`
	Source                string      `json:"source"`
	SourceID              string      `json:"source_id"`
	Sponsor               interface{} `json:"sponsor"`
	SponsorKeys           []int       `json:"sponsor_keys"`
	Status                string      `json:"status"`
	Style                 string      `json:"style"`
	Summary               string      `json:"summary"`
	Syndicate             bool        `json:"syndicate"`
	SyndicateStatus       string      `json:"syndicate_status"`
	Synopsis              string      `json:"synopsis"`
	SynopsisCustom        string      `json:"synopsis_custom"`
	Tags                  []string    `json:"tags"`
	Title                 string      `json:"title"`
	Title2                string      `json:"title2"`
	Title2Text            string      `json:"title2_text"`
	Title3                string      `json:"title3"`
	Title3Text            string      `json:"title3_text"`
	TitleCustom           string      `json:"title_custom"`
	TitleListing          string      `json:"title_listing"`
	TitleListingText      string      `json:"title_listing_text"`
	TitleSectionText      string      `json:"title_section_text"`
	TitleText             string      `json:"title_text"`
	Updated               int64       `json:"updated"`
	Weight                int         `json:"weight"`

	// Companies             []string    `json:"companies"`

	Section struct {
		ID   string
		Link string
		Name string
	}

	Sections []Section

	Subsection struct {
		ID   string
		Link string
		Name string
	}

	Publication struct {
		ID      string
		Link    string
		Name    string
		Primary bool
	}

	Widgets []struct {
		Id   string
		Type string
		Data map[string]any
	}
}

type Section struct {
	Primary      bool
	Publication  string
	Section      string
	Subsection   string
	Thirdsection string
}

type Author struct {
	// Active   bool        `json:"active"`
	// Bio      string      `json:"bio"`
	// Category string      `json:"category"`
	// Email    interface{} `json:"email"`
	// Image    interface{} `json:"image"`
	// Key      int64       `json:"key"`
	// Name     string      `json:"name"`
	// PublicationKey string      `json:"publication_key"`
	// Slug           string      `json:"slug"`
	// Tel            interface{} `json:"tel"`
	// Title          string      `json:"title"`
	//
	//	Social struct {
	//		Facebook  interface{} `json:"facebook"`
	//		Instagram interface{} `json:"instagram"`
	//		Linkedin  interface{} `json:"linkedin"`
	//		Twitter   interface{} `json:"twitter"`
	//	}
}

type Image struct {
	Author      string   `json:"author"`
	Average     string   `json:"average"`
	BlobKey     string   `json:"blob_key"`
	BlobPath    string   `json:"blob_path"`
	Blur        string   `json:"blur"`
	ContentType string   `json:"content_type"`
	Description string   `json:"description"`
	Filename    string   `json:"filename"`
	Filepath    string   `json:"filepath"`
	FocalX      int32    `json:"focal_x"`
	FocalY      int32    `json:"focal_y"`
	Width       int32    `json:"width"`
	Height      int32    `json:"height"`
	Keywords    []string `json:"keywords"`
	Palette     []string `json:"palette"`
	Title       string   `json:"title"`
	Link        string   `json:"link"`
	Mode        string   `json:"mode"`
}

//////////// Widgets

type OEmbedData struct {
	AuthorName      string `json:"author_name"`
	AuthorURL       string `json:"author_url"`
	HTML            string `json:"html"`
	ProviderName    string `json:"provider_name"`
	ProviderURL     string `json:"provider_url"`
	ThumbnailURL    string `json:"thumbnail_url"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	URL             string `json:"url"`
	Version         string `json:"version"`
	Title           string `json:"title"`
	CacheAge        int    `json:"cache_age"`
	ThumbnailHeight int    `json:"thumbnail_height"`
	ThumbnailWidth  int    `json:"thumbnail_width"`
	Type            string `json:"type"`
}

type FacebookVideoWidget struct {
	URL  string
	Meta OEmbedData

	Listing struct {
		Aggregate bool   `json:"aggregate"`
		Author    string `json:"author"`
		Image     Image  `json:"image"`
		Title     string `json:"title"`
	}
}

type InstagramWidget struct {
	URL  string
	Meta OEmbedData
}

type TwitterWidget struct {
	URL  string
	Meta OEmbedData
}

type TextWidget struct {
	Clean string
	HTML  string
	Text  string
}

type InfoblockWidget struct {
	Description string
	Float       string
	Title       string
}

type QuoteWidget struct {
	Cite  string
	Float string
	Quote string
}

type GalleryWidget struct {
	Images Image
	Style  string
}

type ImageWidget struct {
	ID    string
	Image Image
}

type SoundCloudWidget struct {
	Autoplay bool
	Height   int
	Style    string
	URL      string
	Meta     OEmbedData
}

type IonoWidget struct {
	OEmbedData
}

type YoutubeWidget struct {
	ID       string `json:"id"`
	Pid      string `json:"pid"`
	ShareURL string `json:"share_url"`
	URL      string `json:"url"`

	Listing struct {
		Aggregate   bool   `json:"aggregate"`
		Author      string `json:"author"`
		Description string `json:"description"`
		Title       string `json:"title"`
		Image       Image  `json:"image"`
	}

	Meta struct {
		Channel              string `json:"channel"`
		Description          string `json:"description"`
		Published            string `json:"published"`
		Thumbnail            string `json:"thumbnail"`
		ThumbnailSmall       string `json:"thumbnail_small"`
		ThumbnailSmallRetina string `json:"thumbnail_small_retina"`
		Title                string `json:"title"`
	}
}

type GiphyWidget struct {
	Gif struct {
		AnalyticsResponsePayload string `json:"analytics_response_payload"`
		BitlyGifURL              string `json:"bitly_gif_url"`
		BitlyURL                 string `json:"bitly_url"`
		ContentURL               string `json:"content_url"`
		EmbedURL                 string `json:"embed_url"`
		ID                       string `json:"id"`
		Images                   struct {
			Four80WStill struct {
				Height string `json:"height"`
				Size   string `json:"size"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"480w_still"`
			Downsized struct {
				Height string `json:"height"`
				Size   string `json:"size"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"downsized"`
			DownsizedLarge struct {
				Height string `json:"height"`
				Size   string `json:"size"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"downsized_large"`
			DownsizedMedium struct {
				Height string `json:"height"`
				Size   string `json:"size"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"downsized_medium"`
			DownsizedSmall struct {
				Height  string `json:"height"`
				Mp4     string `json:"mp4"`
				Mp4Size string `json:"mp4_size"`
				Width   string `json:"width"`
			} `json:"downsized_small"`
			DownsizedStill struct {
				Height string `json:"height"`
				Size   string `json:"size"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"downsized_still"`
			FixedHeight struct {
				Height   string `json:"height"`
				Mp4      string `json:"mp4"`
				Mp4Size  string `json:"mp4_size"`
				Size     string `json:"size"`
				URL      string `json:"url"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
				Width    string `json:"width"`
			} `json:"fixed_height"`
			FixedHeightDownsampled struct {
				Height   string `json:"height"`
				Size     string `json:"size"`
				URL      string `json:"url"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
				Width    string `json:"width"`
			} `json:"fixed_height_downsampled"`
			FixedHeightSmall struct {
				Height   string `json:"height"`
				Mp4      string `json:"mp4"`
				Mp4Size  string `json:"mp4_size"`
				Size     string `json:"size"`
				URL      string `json:"url"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
				Width    string `json:"width"`
			} `json:"fixed_height_small"`
			FixedHeightSmallStill struct {
				Height string `json:"height"`
				Size   string `json:"size"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"fixed_height_small_still"`
			FixedHeightStill struct {
				Height string `json:"height"`
				Size   string `json:"size"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"fixed_height_still"`
			FixedWidth struct {
				Height   string `json:"height"`
				Mp4      string `json:"mp4"`
				Mp4Size  string `json:"mp4_size"`
				Size     string `json:"size"`
				URL      string `json:"url"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
				Width    string `json:"width"`
			} `json:"fixed_width"`
			FixedWidthDownsampled struct {
				Height   string `json:"height"`
				Size     string `json:"size"`
				URL      string `json:"url"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
				Width    string `json:"width"`
			} `json:"fixed_width_downsampled"`
			FixedWidthSmall struct {
				Height   string `json:"height"`
				Mp4      string `json:"mp4"`
				Mp4Size  string `json:"mp4_size"`
				Size     string `json:"size"`
				URL      string `json:"url"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
				Width    string `json:"width"`
			} `json:"fixed_width_small"`
			FixedWidthSmallStill struct {
				Height string `json:"height"`
				Size   string `json:"size"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"fixed_width_small_still"`
			FixedWidthStill struct {
				Height string `json:"height"`
				Size   string `json:"size"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"fixed_width_still"`
			Looping struct {
				Mp4     string `json:"mp4"`
				Mp4Size string `json:"mp4_size"`
			} `json:"looping"`
			Original struct {
				Frames   string `json:"frames"`
				Hash     string `json:"hash"`
				Height   string `json:"height"`
				Mp4      string `json:"mp4"`
				Mp4Size  string `json:"mp4_size"`
				Size     string `json:"size"`
				URL      string `json:"url"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
				Width    string `json:"width"`
			} `json:"original"`
			OriginalMp4 struct {
				Height  string `json:"height"`
				Mp4     string `json:"mp4"`
				Mp4Size string `json:"mp4_size"`
				Width   string `json:"width"`
			} `json:"original_mp4"`
			OriginalStill struct {
				Height string `json:"height"`
				Size   string `json:"size"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"original_still"`
			Preview struct {
				Height  string `json:"height"`
				Mp4     string `json:"mp4"`
				Mp4Size string `json:"mp4_size"`
				Width   string `json:"width"`
			} `json:"preview"`
			PreviewGif struct {
				Height string `json:"height"`
				Size   string `json:"size"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"preview_gif"`
			PreviewWebp struct {
				Height string `json:"height"`
				Size   string `json:"size"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"preview_webp"`
		}
		ImportDatetime   string `json:"import_datetime"`
		IsSticker        int    `json:"is_sticker"`
		Rating           string `json:"rating"`
		Slug             string `json:"slug"`
		Source           string `json:"source"`
		SourcePostURL    string `json:"source_post_url"`
		SourceTld        string `json:"source_tld"`
		Title            string `json:"title"`
		TrendingDatetime string `json:"trending_datetime"`
		Type             string `json:"type"`
		URL              string `json:"url"`

		User struct {
			AvatarURL    string `json:"avatar_url"`
			BannerImage  string `json:"banner_image"`
			BannerURL    string `json:"banner_url"`
			Description  string `json:"description"`
			DisplayName  string `json:"display_name"`
			InstagramURL string `json:"instagram_url"`
			IsVerified   bool   `json:"is_verified"`
			ProfileURL   string `json:"profile_url"`
			Username     string `json:"username"`
			WebsiteURL   string `json:"website_url"`
		}

		Analytics struct {
			Onclick struct {
				URL string `json:"url"`
			}
			Onload struct {
				URL string `json:"url"`
			}
			Onsent struct {
				URL string `json:"url"`
			}
		}

		Username string `json:"username"`
	}
}

type JwPlayerWidget struct {
	ID  string `json:"id"`
	URL string `json:"url"`

	Listing struct {
		Aggregate   bool   `json:"aggregate"`
		Author      string `json:"author"`
		Description string `json:"description"`
		Image       Image  `json:"image"`
		Title       string `json:"title"`
	}

	Meta struct {
		Description     string `json:"description"`
		HTML            string `json:"html"`
		Thumbnail       string `json:"thumbnail"`
		ThumbnailHeight int    `json:"thumbnail_height"`
		ThumbnailWidth  int    `json:"thumbnail_width"`
		Title           string `json:"title"`
		Type            string `json:"type"`
		URL             string `json:"url"`
		Version         string `json:"version"`
	}
}

type HtmlWidget struct {
	HTML string
}

type ChartblocksWidget struct {
	ID  string `json:"id"`
	URL string `json:"url"`

	Meta struct {
		ChartName   string    `json:"chart_name"`
		Created     time.Time `json:"created"`
		IsPublic    bool      `json:"isPublic"`
		UpdatedAt   time.Time `json:"updated_at"`
		VersionHash string    `json:"version_hash"`

		Creator struct {
			AvatarURL string `json:"avatarUrl"`
			ID        string `json:"id"`
			Nickname  string `json:"nickname"`

			Account struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
		}

		Images struct {
			Eps struct {
				Enabled   bool   `json:"enabled"`
				Extension string `json:"extension"`
				Name      string `json:"name"`
				URL       string `json:"url"`
			}
			Pdf struct {
				Enabled   bool   `json:"enabled"`
				Extension string `json:"extension"`
				Name      string `json:"name"`
				URL       string `json:"url"`
			}
			Png struct {
				Enabled   bool   `json:"enabled"`
				Extension string `json:"extension"`
				Name      string `json:"name"`
				URL       string `json:"url"`
			}
			Ps struct {
				Enabled   bool   `json:"enabled"`
				Extension string `json:"extension"`
				Name      string `json:"name"`
				URL       string `json:"url"`
			}
			Svg struct {
				Enabled   bool   `json:"enabled"`
				Extension string `json:"extension"`
				Name      string `json:"name"`
				URL       string `json:"url"`
			}
		}
	}
}

type PollDaddyWidget struct {
	ID   string
	URL  string
	Meta OEmbedData
}

type AccordianWidget struct {
	WidgetTitle string

	Accordions []struct {
		Text      string
		Title     string
		Accordion string
	}
}

type ArticleList struct {
	ArticleIds   []int64 `json:"article_ids"`
	ReadMoreLink string  `json:"read_more_link"`

	Query struct {
		Author      string `json:"author"`
		ContentType string `json:"contentType"`
		DateFrom    string `json:"dateFrom"`
		DateTo      string `json:"dateTo"`
		Group       string `json:"group"`
		Limit       int    `json:"limit"`
		Offset      int    `json:"offset"`
		Page        int    `json:"page"`
		Publication string `json:"publication"`
		Query       string `json:"query"`
		Sponsor     string `json:"sponsor"`
		Tag         string `json:"tag"`
	}
}

type IssuWidget struct {
	URL  string
	Meta OEmbedData
}

type ScribdWidget struct {
	ID   string
	URL  string
	Meta OEmbedData
}

type InfogramWidget struct {
	HTML string
	Meta OEmbedData
}

type GoogleMapWidget struct {
	Address     string
	Coordinates string
	Routes      interface{}
	Type        string
	Zoom        int
}
