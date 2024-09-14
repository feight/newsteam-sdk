package lib

import (
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func Json[T any](path, accessToken string) (*T, error) {

	client := resty.New()

	response, err := client.R().
		SetCookie(&http.Cookie{
			Name:  "_cosmos_auth",
			Value: accessToken,
		}).
		Get(path)

	if err != nil {
		return nil, err
	}

	dst := new(T)

	return dst, json.Unmarshal(response.Body(), dst)
}
