package goexchange

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

// signing a message  using: hmac sha256 + base64
func HmacSha256Base64Signer(message string, secretKey string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(message))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// signing a message  using: hmac sha256
func HmacSha256Signer(message string, secretKey string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(message))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", mac.Sum(nil)), nil
}

// signing a message  using: hmac sha512
func HmacSha512Signer(message string, secretKey string) (string, error) {
	mac := hmac.New(sha512.New, []byte(secretKey))
	_, err := mac.Write([]byte(message))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", mac.Sum(nil)), nil
}

// md5 sign
func Md5Signer(message string) string {
	hash := md5.Sum([]byte(message))
	return fmt.Sprintf("%x", hash)
}

// Sha512Signer sign
func Sha512Signer(message string) string {
	hash := sha512.Sum512([]byte(message))
	return fmt.Sprintf("%x", hash)
}

// signing a message  using: hmac sha256 + base64
func Base64Signer(message string) string {
	return base64.StdEncoding.EncodeToString([]byte(message))
}

// Get a iso time eg: 2018-03-16T18:02:48.284Z
func IsoTime() string {
	utcTime := time.Now().UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}

// Get a iso time eg: 2018-03-16T18:02:48.284Z
func GetNowUtcTime() string {
	utcTime := time.Now().UTC()
	return utcTime.Format("2006-01-02T15:04:05")
}

// GetNowMillisecond Get current mill second timestamp
// eg: 1521221737376
func GetNowMillisecond() int64 {
	return time.Now().UnixNano() / 1000000
}

// GetNowMillisecond Get current mill second timestamp
// eg: 1521221737376
func GetNowMillisecondStr() string {
	return strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
}

// GetNowMicrosecondStr Get current mill second timestamp
// eg: 1521221737376000
func GetNowMicrosecondStr() string {
	return strconv.FormatInt(time.Now().UnixNano()/1000, 10)
}

// getNowTimestamp Get current second timestamp
// eg: 1521221737
func GetNowTimestamp() int64 {
	return time.Now().Unix()
}

// GetNowTimestampStr Get current second timestamp
// eg: 1521221737
func GetNowTimestampStr() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// build http get request params, and order return string: eg: aa=111&bb=222&cc=333
func BuildParams(params map[string]string) string {
	urlParams := url.Values{}
	for k := range params {
		urlParams.Add(k, (params)[k])
	}
	return urlParams.Encode()
}

// LoadConfig load exchange config
func LoadConfig(exchange string) (map[string]interface{}, error) {
	file, err := os.Open("../config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	result := map[string]interface{}{}
	decoder.Decode(&result)

	// retData := map[string]string{}
	if _, ok := result[exchange]; !ok {
		return nil, errors.New("exchange do not exist api config")
	}
	return result[exchange].(map[string]interface{}), nil
}
