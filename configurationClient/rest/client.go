package rest

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type MambuConfigClient struct {
	client  http.Client
	baseURL *url.URL
	apikey  string
}

const acceptHeader = "application/vnd.mambu.v2+yaml"

func NewClient(mambuURL string, apikey string) *MambuConfigClient {
	baseURL, err := url.Parse(mambuURL)
	if err != nil {
		log.Fatalf("fail to parse url, err: %s", err)
	}
	return &MambuConfigClient{baseURL: baseURL, apikey: apikey, client: http.Client{Timeout: 1 * time.Minute}}
}

func (c MambuConfigClient) resolvePath(path string, query url.Values) (string, error) {
	relativeURL, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	relativeURL.RawQuery = query.Encode()

	absoluteURL := c.baseURL.ResolveReference(relativeURL)
	return absoluteURL.String(), nil
}

func (c MambuConfigClient) sendRequest(method string, path string, body io.Reader, query url.Values) ([]byte, error) {
	client := c.client
	absolutePath, err := c.resolvePath(path, query)
	fmt.Println(absolutePath)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(method, absolutePath, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", acceptHeader)
	request.Header.Add("ApiKey", c.apikey)
	resp, err := client.Do(request)
	fmt.Printf("%+v\n", request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(bs))
	return bs, nil
}
