package newsteam

import (
	"encoding/base64"
	"io"
	"net/http"

	"buf.build/gen/go/dgroux/newsteam/protocolbuffers/go/admin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

/*
 * UploadImage
 */
func UploadImage(reader func() io.ReadCloser) *admin.Image {

	id := "/image/" + uuid.New().String()

	register(id, func(w http.ResponseWriter, r *http.Request) {

		body := reader()
		defer body.Close()

		if _, err := io.Copy(w, body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return &admin.Image{Id: id}
}

/*
 * UploadImageFromUrl
 */
func UploadImageFromUrl(url string) *admin.Image {

	return UploadImage(func() io.ReadCloser {

		// fmt.Printf("Uploading image `%s`...\n", url)

		response, err := http.Get(url)
		if err != nil {
			panic(errors.Wrapf(err, "could not get image [%s]", url))
		}

		return response.Body
	})
}

/*
 * UploadImageFromUrl_BasicAuth
 */
func UploadImageFromUrl_BasicAuth(url string, username string, password string) *admin.Image {

	return UploadImage(func() io.ReadCloser {

		// fmt.Printf("Uploading image `%s`...\n", url)

		credentials := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}

		req.Header.Set("authorization", "basic "+credentials)

		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		return response.Body
	})
}
