package goexchange

import (
	"strings"

	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// http request mode
const (
	HTTP_GET    string = "GET"
	HTTP_POST   string = "POST"
	HTTP_DELETE string = "DELETE"
)

// NewHttpRequest http request
func NewHttpRequest(client *http.Client, method string, reqURL string, postData string, headers map[string]string) HttpClientResponse {
	startTime := time.Now().UnixNano() / 1e6
	req, _ := http.NewRequest(method, reqURL, strings.NewReader(postData))
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	endTime := time.Now().UnixNano() / 1e6
	var returnData HttpClientResponse
	returnData.St = startTime
	returnData.Et = endTime
	if err != nil {
		returnData.Code = HttpClientInternalError.Code
		returnData.Msg = HttpClientInternalError.Msg
		returnData.Error = err.Error()
		return returnData
	}

	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		returnData.Code = HttpClientInternalError.Code
		returnData.Msg = HttpClientInternalError.Msg
		returnData.Error = err.Error()
		return returnData
	}

	fmt.Println(fmt.Sprintf("url:%s, param:%s, status:%d, time:[%d|%d|%d], resp:%s", reqURL, postData, resp.StatusCode, startTime, endTime, endTime-startTime, string(bodyData)))

	if resp.StatusCode != 200 {
		returnData.Code = resp.StatusCode
		returnData.Msg = HttpRequestError.Msg
		returnData.Error = fmt.Sprintf("HttpStatusCode:%d, Desc:%s", resp.StatusCode, string(bodyData))
		return returnData
	}
	returnData.Data = bodyData
	return returnData
}

func HttpGet(client *http.Client, reqURL string) HttpClientResponse {
	respData := NewHttpRequest(client, "GET", reqURL, "", map[string]string{})
	return respData
}

func HttpGetWithHeader(client *http.Client, reqURL string, headers map[string]string) HttpClientResponse {
	respData := NewHttpRequest(client, "GET", reqURL, "", headers)
	return respData
}

func HttpPost(client *http.Client, reqURL string, postData string) HttpClientResponse {
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	return NewHttpRequest(client, "POST", reqURL, postData, headers)
}

func HttpPostWithHeader(client *http.Client, reqURL string, postData string, headers map[string]string) HttpClientResponse {
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return NewHttpRequest(client, "POST", reqURL, postData, headers)
}

func HttpPostWithJson(client *http.Client, reqURL string, postData string, headers map[string]string) HttpClientResponse {
	headers["Content-Type"] = "application/json; charset=UTF-8"
	return NewHttpRequest(client, "POST", reqURL, postData, headers)
}

func HttpDelete(client *http.Client, reqURL string) HttpClientResponse {
	return NewHttpRequest(client, "DELETE", reqURL, "", nil)
}

func HttpDeleteWithHeader(client *http.Client, reqURL string, headers map[string]string) HttpClientResponse {
	headers["Content-Type"] = "application/json; charset=UTF-8"
	return NewHttpRequest(client, "DELETE", reqURL, "", headers)
}
