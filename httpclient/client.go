package httpclient

import (
	"io/ioutil"
	"net/http"
	"os"
)

func RequestSparkDashboard() (string, error) {
	url := os.Getenv("SPARK_DASHBOARD_URL")
	login := os.Getenv("SPARK_LOGIN")
	pass := os.Getenv("SPARK_PASSWORD")

	return RequestPage(url, login, pass)
}

func RequestPage(url string, login string, pass string) (string, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(login, pass)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	return s, nil
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
