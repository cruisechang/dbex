package dbex

import (
	"errors"
	"net/http"
	"time"
	"strings"
	"io/ioutil"
)

type httpClient struct {
	address   string
	port      string
	url       string
	client *http.Client
	postURI   string
	postQuery string
}

func newClient(address, port string, connectTimeout, handshakeTimeout, requestTimeout int) (*httpClient, error) {
	re := &httpClient{
		address: address,
		port:    port,
		url:     "http://" + address + ":" + port,
	}

	var netTransport = &http.Transport{
		//Dial: (&net.Dialer{
		//	Timeout: connectTimeout * time.Second,
		//}).Dial,
		TLSHandshakeTimeout: time.Duration(handshakeTimeout) * time.Second,
	}

	re.client = &http.Client{
		Timeout:   time.Duration(requestTimeout) * time.Second,
		Transport: netTransport,
	}

	return re, nil
}

func (re *httpClient) URL() string {
	return re.url
}
func (h *httpClient) SetURL(addr, port string) {

	h.url = "http://" + addr + ":" + port
}

func (h *httpClient) Do(req *http.Request) (*http.Response, error) {

	return h.client.Do(req)

}

//SetPostURI
//queryPair is a map , key is query key, value is query value. eg.  ["data"]{"xxxx"} => data=xxxx
//query => data=""&age=33....
func (re *httpClient) SetPostURI(path string, queryPair map[string]string) error {
	if len(queryPair) <= 0 {
		return errors.New("query pair is len=0")
	}
	var query string
	for k, v := range queryPair {
		query = k + "=" + v + "&"
	}

	query = strings.TrimRight(query, "&")

	re.postQuery = query
	re.postURI = re.url + path
	return nil
}

//Get after set post url
func (re *httpClient) PostURI() string {
	return re.postURI
}

func (re *httpClient) PostQuery() string {
	return re.postQuery
}
func (re *httpClient) Post() (string, error) {
	resp, err := re.client.Post(re.postURI, "application/x-www-form-urlencoded", strings.NewReader(re.postQuery))

	if err != nil {
		return "", errors.New("HTTPPost error")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("HTTPPost readAll error")
	}

	return string(body), nil

}
