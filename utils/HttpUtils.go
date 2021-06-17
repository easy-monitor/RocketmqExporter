package utils

import (
	"io/ioutil"
	"net/http"
)

func HttpUrl(method string, url string, cookie string) []byte {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil
	}
	if cookie != "" {
		req.AddCookie(&http.Cookie{
			Name:  "JSESSIONID",
			Value: cookie,
		})
	}
	//req.Header.Add("accept", "application/json")
	//req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil
	} else {
		return body
	}
}
