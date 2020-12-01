package goexchange

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "net/url"
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

// Get a iso time eg: 2018-03-16T18:02:48.284Z
func IsoTime() string {
    utcTime := time.Now().UTC()
    iso := utcTime.String()
    isoBytes := []byte(iso)
    iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
    return iso
}


// build http get request params, and order return string: eg: aa=111&bb=222&cc=333
func BuildParams(params map[string]string) string {
    urlParams := url.Values{}
    for k := range params {
        urlParams.Add(k, (params)[k])
    }
    return urlParams.Encode()
}