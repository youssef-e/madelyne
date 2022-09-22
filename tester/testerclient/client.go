package testerclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Requester interface {
	Make(r Request) (Response, error)
}

type Request struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    io.Reader
}

type Response struct {
	StatusCode  int
	Body        io.ReadCloser
	ContentType string
	Headers     map[string][]string
}

type Client struct {
	httpClient *http.Client
	baseUrl    string
}

func New(baseUrl string) Client {
	return Client{
		httpClient: &http.Client{},
		baseUrl:    baseUrl,
	}
}

func (c Client) Make(r Request) (Response, error) {
	u, err := encodeUrl(r.Url)
	if err != nil {
		return Response{}, err
	}

	request, err := http.NewRequest(r.Method, c.baseUrl+u, r.Body)
	if err != nil {
		return Response{}, err
	}

	for key, value := range r.Headers {
		request.Header.Add(key, value)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return Response{}, err
	}

	return Response{
		StatusCode:  response.StatusCode,
		Body:        response.Body,
		ContentType: response.Header.Get("Content-Type"),
		Headers:     response.Header,
	}, nil
}

func (c Client) Get(url string, headers map[string]string) (Response, error) {
	return c.Make(Request{
		Method:  "GET",
		Url:     url,
		Headers: headers,
	})
}

func (c Client) Post(url string, body io.Reader, headers map[string]string) (Response, error) {
	return c.Make(Request{
		Method:  "POST",
		Url:     url,
		Headers: headers,
		Body:    body,
	})
}

func (c Client) Put(url string, body io.Reader, headers map[string]string) (Response, error) {
	return c.Make(Request{
		Method:  "PUT",
		Url:     url,
		Headers: headers,
		Body:    body,
	})
}

func (c Client) Patch(url string, body io.Reader, headers map[string]string) (Response, error) {
	return c.Make(Request{
		Method:  "PATCH",
		Url:     url,
		Headers: headers,
		Body:    body,
	})
}

func (c Client) Delete(url string, headers map[string]string) (Response, error) {
	return c.Make(Request{
		Method:  "DELETE",
		Url:     url,
		Headers: headers,
	})
}

func encodeUrl(u string) (string, error) {
	p, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("Can't parse url : %w", err)
	}
	p.RawQuery = p.Query().Encode()
	return p.String(), nil
}
