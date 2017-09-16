package httpclient

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"errors"
	"htmlparser/analyser"
)

const URL_SEPARATOR = ";"

func RequestSparkDashboard(url string) (*string, error) {

	login := os.Getenv("SPARK_LOGIN")
	pass := os.Getenv("SPARK_PASSWORD")

	return RequestPage(url, login, pass)
}

func RequestPage(url string, login string, pass string) (*string, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(login, pass)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	return &s, nil
}

func SparkDashboardIsRequestable() bool {
	url := os.Getenv("SPARK_DASHBOARD_URL")
	login := os.Getenv("SPARK_LOGIN")
	pass := os.Getenv("SPARK_PASSWORD")

	return IsRequestable(url ,login, pass)
}

func IsRequestable(url string, login string, pass string) bool {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(login, pass)

	_, err = client.Do(req)
	return err == nil
}

func GetUrl(index int) (string, error) {

	urls := os.Getenv("SPARK_DASHBOARD_URL")

	if strings.Contains(urls, URL_SEPARATOR) {
		urlsArray := strings.Split(urls, URL_SEPARATOR)

		if index > len(urlsArray) {
			return "", errors.New("Invalid index")
		}
		return urlsArray[index], nil
	} else if index == 0 {
		return urls, nil
	} else {
		return "", errors.New("Invalid index")
	}
}

func getUrlNumber() int {
	urls := os.Getenv("SPARK_DASHBOARD_URL")

	if strings.Contains(urls, URL_SEPARATOR) {
		return len(strings.Split(urls, URL_SEPARATOR))
	} else {
		return 1
	}
}

func RequestActiveSparkMasterContent() (*string, error) {
	return GetActiveSparkMasterContent(RequestSparkDashboard)
}

// This function take a function in parameters just to be easily testable, the real one is RequestActiveSparkMasterContent
func GetActiveSparkMasterContent(funcRequestSparkDashboard func(string) (*string, error)) (*string, error) {

	nb := getUrlNumber()

	if nb == 0 {
		return nil, errors.New("No url found in SPARK_DASHBOARD_URL env var. Please chez it !")
	}

	for i:=0; i<nb; i++ {

		url, err := GetUrl(i)
		content, err := funcRequestSparkDashboard(url)

		if err != nil {
			return nil, err
		}

		if analyser.IsActiveSparkMaster(*content) {
			return content, nil
		}
	}
	return nil, errors.New("None of spark master urls given in parameter are ACTIVE")
}
