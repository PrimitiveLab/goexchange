package goexchange

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func NewHttpRequest(client *http.Client, method string, reqUrl string, postData string, headers map[string]string) ([]byte, error) {
	// logger.Log.Debugf("[%s] request url: %s", method, reqUrl)

	req, _ := http.NewRequest(method, reqUrl, strings.NewReader(postData))
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HttpStatusCode:%d ,Desc:%s", resp.StatusCode, string(bodyData)))
	}

	return bodyData, nil
}

func HttpGet(client *http.Client, reqUrl string, headers map[string]string) (map[string]interface{}, error) {
	respData, err := NewHttpRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		log.Println(string(respData))
		return nil, err
	}
	return bodyDataMap, nil
}

func HttpPostForm(client *http.Client, reqUrl string, postData url.Values) ([]byte, error) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded"}
	return NewHttpRequest(client, "POST", reqUrl, postData.Encode(), headers)
}

func HttpPostWithJson(client *http.Client, reqUrl string, postData url.Values) ([]byte, error) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded"}
	return NewHttpRequest(client, "POST", reqUrl, postData.Encode(), headers)
}

func HttpPostWithFormUrlEncoded(client *http.Client, reqUrl string, postData url.Values) ([]byte, error) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded"}
	return NewHttpRequest(client, "POST", reqUrl, postData.Encode(), headers)
}

func HttpDelete(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return NewHttpRequest(client, "DELETE", reqUrl, postData.Encode(), headers)
}

func HttpDeleteWithForm(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return NewHttpRequest(client, "DELETE", reqUrl, postData.Encode(), headers)
}

func HttpPut(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return NewHttpRequest(client, "PUT", reqUrl, postData.Encode(), headers)
}