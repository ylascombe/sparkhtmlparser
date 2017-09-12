package httpclient

import (
	"io/ioutil"
	"net/http"
)

func RequestPage(url string, login string, pass string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("spark", "spark")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	return s, nil
}

func IsRequestable(url string, login string, pass string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("spark", "spark")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	return (err == nil)
}
