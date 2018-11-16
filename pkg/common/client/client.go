package client

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)
func newTimeoutClient() *http.Client {
	var transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	var client = &http.Client{
		Timeout:   time.Second * 30,
		Transport: transport,
	}
	return client
}
func getReqWrapper(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}
func doReq(req *http.Request) ([]byte, error) {
	//res, err := newTimeoutClient().Do(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("resp status:" + fmt.Sprint(res.StatusCode))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Download(url string) ([]byte, error) {
	req, err := getReqWrapper("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return doReq(req)
}