package bottle

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type NodeGetter interface {
	NodeGet(bottle, key string) ([]byte, error)
}

type NodeRouter interface {
	NodeRoute(key string) (node NodeGetter, ok bool)
}

type httpGetter struct {
	addr string
}

func (h *httpGetter) NodeGet(bottle, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v",
		h.addr,
		url.QueryEscape(bottle),
		url.QueryEscape(key),
	)

	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("return status: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("response body: %v", err)
	}

	return bytes, err
}
