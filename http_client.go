package goexchange

import (
	"strings"

	// "errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// 请求类型
const (
	HTTP_GET    string = "GET"
	HTTP_POST   string = "POST"
	HTTP_DELETE string = "DELETE"
)

func NewHttpRequest(client *http.Client, method string, reqUrl string, postData string, headers map[string]string) HttpClientResponse {
	// logger.Log.Debugf("[%s] request url: %s", method, reqUrl)

	startTime := time.Now().UnixNano() / 1e6
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

	if resp.StatusCode != 200 {
		returnData.Code = resp.StatusCode
		returnData.Msg = HttpRequestError.Msg
		returnData.Error = fmt.Sprintf("HttpStatusCode:%d ,Desc:%s", resp.StatusCode, string(bodyData))
		return returnData
	}
	returnData.Data = bodyData
	return returnData
}

func HttpGet(client *http.Client, reqUrl string) HttpClientResponse {
	respData := NewHttpRequest(client, "GET", reqUrl, "", map[string]string{})
	return respData
}

func HttpGetWithHeader(client *http.Client, reqUrl string, headers map[string]string) HttpClientResponse {
	respData := NewHttpRequest(client, "GET", reqUrl, "", headers)
	return respData
}

func HttpPost(client *http.Client, reqUrl string, postData string) HttpClientResponse {
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	return NewHttpRequest(client, "POST", reqUrl, postData, headers)
}


func HttpPostWithJson(client *http.Client, reqUrl string, postData string, headers map[string]string) HttpClientResponse {
	headers["Content-Type"] = "application/json; charset=UTF-8"
	return NewHttpRequest(client, "POST", reqUrl, postData, headers)
}

//
// func HttpPostWithFormUrlEncoded(client *http.Client, reqUrl string, postData url.Values) ([]byte, [2]int64,error) {
// 	headers := map[string]string{
// 		"Content-Type": "application/x-www-form-urlencoded"}
// 	return NewHttpRequest(client, "POST", reqUrl, postData.Encode(), headers)
// }
//
// func HttpDelete(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte,[2]int64, error) {
// 	if headers == nil {
// 		headers = map[string]string{}
// 	}
// 	headers["Content-Type"] = "application/x-www-form-urlencoded"
// 	return NewHttpRequest(client, "DELETE", reqUrl, postData.Encode(), headers)
// }
//
// func HttpDeleteWithForm(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, [2]int64, error) {
// 	if headers == nil {
// 		headers = map[string]string{}
// 	}
// 	headers["Content-Type"] = "application/x-www-form-urlencoded"
// 	return NewHttpRequest(client, "DELETE", reqUrl, postData.Encode(), headers)
// }
//
// func HttpPut(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, [2]int64, error) {
// 	if headers == nil {
// 		headers = map[string]string{}
// 	}
// 	headers["Content-Type"] = "application/x-www-form-urlencoded"
// 	return NewHttpRequest(client, "PUT", reqUrl, postData.Encode(), headers)
// }