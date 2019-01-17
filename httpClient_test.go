package dbex

import (
	"encoding/json"
	"testing"
)

func TestGetURL(t *testing.T) {
	n, _ := newClient("127.0.0.1", "8888", 5, 5, 10)

	t.Logf("GetURL:%s\n", n.URL())

}

func TestRequester_GetPostURL(t *testing.T) {
	n, _ := newClient("www.yahoo.com.tw", "8888", 5, 5, 1)

	c := &struct {
		SessionID string
	}{
		SessionID: "sid",
	}

	js, _ := json.Marshal(c)

	query := map[string]string{}
	query["data"] = string(js)
	n.SetPostURI("login", query)

	t.Logf("GetPostURI:%s\n", n.PostURI())
}
func TestPOST(t *testing.T) {
	n, _ := newClient("www.yahoo.com.tw", "8888", 5, 5, 1)

	c := &struct {
		SessionID string
	}{
		SessionID: "sid",
	}

	js, _ := json.Marshal(c)

	query := map[string]string{}
	query["data"] = string(js)

	n.SetPostURI("login", query)

	res, err := n.Post()
	if err != nil {
		t.Error("post failed")
		return
	}
	t.Logf("return body:%s", res)

}
