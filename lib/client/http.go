package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Options struct {
	Method         string // TODO: Use only Options in request(), no extra params
	Path           string
	Headers        http.Header
	Client         *http.Client
	Data           any
	AuthCookieName string
	AccessToken    string
}

type Error map[string]any

func (e Error) Error() string {

	if v, ok := e["error"].(string); ok {
		return v
	}
	if v, ok := e["message"].(string); ok {
		return v
	}
	return ""
}

// Gets the inner message if available, otherwise returns
// entire error.
func (e Error) inner() map[string]any {
	err := e.Error()

	/* If error is empty then return the entire error */
	if err == "" {
		return e
	}

	/* Return formatted inner error */
	return map[string]any{
		"message": e.Error(),
	}
}

func (e Error) Write(w http.ResponseWriter) {
	b, _ := json.Marshal(e.inner())
	w.Write(b)
}

type httpResponse[T any] struct {
	*http.Response
	dst  *T
	opts Options
}

// ReadAll reads the entirety of the body.
func (s *httpResponse[T]) ReadAll() ([]byte, error) {

	defer s.Body.Close()

	/*
	 * Read the response body
	 */
	body, err := io.ReadAll(s.Body)

	if err != nil {
		return nil, fmt.Errorf(
			"could not read response [%s]: %w", s.opts.Path, err)
	}

	/*
	 * Handle response error
	 */
	if s.StatusCode != http.StatusOK {
		return nil, s.getResponseError(body)
	}

	return body, nil
}

// GetErrorResponse gets the response error if any.
func (s *httpResponse[T]) getResponseError(body []byte) error {

	var responseError Error
	if err := json.Unmarshal(body, &responseError); err != nil {
		return errors.New(string(body))
	}

	return responseError
}

func (s *httpResponse[T]) do() (*httpResponse[T], error) {

	path := s.opts.Path

	/*
	 * Encode the payload to json
	 */
	b, err := json.Marshal(s.opts.Data)
	if err != nil {
		return nil, fmt.Errorf(
			"could not marshal request data [%s]: %w", path, err)
	}

	method := s.opts.Method
	if method == "" {
		method = http.MethodPost
	}

	/*
	 * Create the request object
	 */
	req, err := http.NewRequest(method, path, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf(
			"could not create request [%s]: %w", path, err)
	}

	/*
	 * Set the headers
	 */
	if s.opts.Headers != nil {
		req.Header = s.opts.Headers
	}

	if s.opts.AccessToken != "" {

		if s.opts.AuthCookieName == "" {
			s.opts.AuthCookieName = "_cosmos_auth"
		}

		req.AddCookie(&http.Cookie{
			Name:  s.opts.AuthCookieName,
			Value: s.opts.AccessToken})
	}

	req.Header.Add("content-type", "application/json")

	/*
	 * Perform the request
	 */
	client := s.opts.Client
	if client == nil {
		client = &http.Client{Timeout: time.Second * 60}
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"could not send request [%s]: %w", path, err)
	}

	s.Response = res

	return s, nil
}

func (s *httpResponse[T]) Json() (*T, error) {

	body, err := s.ReadAll()
	if err != nil {
		return nil, err
	}

	/*
	 * Decode the response into dst
	 */
	err = json.Unmarshal(body, s.dst)

	if err != nil {
		return nil, fmt.Errorf("could not decode json [%s]: %w: body=%s",
			s.opts.Path, err, body)
	}

	return s.dst, nil
}

func NewHttpRequest(opts Options) (*httpResponse[any], error) {
	r := &httpResponse[any]{opts: opts}
	return r.do()
}

func Json[T any]( /*dst *T, */ opts Options) (*T, error) {

	r := &httpResponse[T]{opts: opts, dst: new(T)}
	if _, err := r.do(); err != nil {
		return nil, err
	}

	return r.Json()
}
