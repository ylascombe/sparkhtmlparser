package httpclient

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUrlIndex0(t *testing.T) {

	// arrange
	os.Setenv("SPARK_DASHBOARD_URL", "http://server1.com:8080/,http://server2.com:8080/")

	// act
	res, err := GetUrl(0)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, "http://server1.com:8080/", res)
}

func TestGetUrlIndex1(t *testing.T) {

	// arrange
	os.Setenv("SPARK_DASHBOARD_URL", "http://server1.com:8080/,http://server2.com:8080/")

	// act
	res, err := GetUrl(1)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, "http://server2.com:8080/", res)
}

func TestGetUrlOnlyOneUrl(t *testing.T) {

	// arrange
	os.Setenv("SPARK_DASHBOARD_URL", "http://server1.com:8080/")

	// act
	res, err := GetUrl(0)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, "http://server1.com:8080/", res)
}

func TestGetUrlInvalidIndex(t *testing.T) {

	// arrange
	os.Setenv("SPARK_DASHBOARD_URL", "http://server1.com:8080/")

	// act
	res, err := GetUrl(1)

	// assert
	assert.NotNil(t, err)
	assert.Equal(t, "", res)
}

func TestGetUrlNumberWhen1(t *testing.T) {
	// arrange
	os.Setenv("SPARK_DASHBOARD_URL", "http://server1.com:8080/")

	// act
	res := getUrlNumber()

	// assert
	assert.Equal(t, 1, res)
}

func TestGetUrlNumberWhen2(t *testing.T) {
	// arrange
	os.Setenv("SPARK_DASHBOARD_URL", "http://server1.com:8080/,http://server2.com:8080/")

	// act
	res := getUrlNumber()

	// assert
	assert.Equal(t, 2, res)
}

func TestGetActiveSparkMasterContent(t *testing.T) {
	// arrange
	url1 := "http://server1.com:8080/"
	url2 := "http://server2.com:8080/"
	os.Setenv("SPARK_DASHBOARD_URL", fmt.Sprintf("%s,%s", url1, url2))

	htmlContentWithActiveStatus := "<html><body><ul><li><strong>Status:</strong> ALIVE</li></ul></body></html>"

	fakeRequestSparkDashboardFunc := func(url string) (*string, error) {
		if url == url1 {
			res := "<html><body><ul><li><strong>Status:</strong> STAND BY</li></ul></body></html>"
			return &res, nil
		} else if url == url2 {
			res := htmlContentWithActiveStatus
			return &res, nil
		} else {
			return nil, errors.New("Not expected case")
		}
	}

	// act
	res, err := GetActiveSparkMasterContent(fakeRequestSparkDashboardFunc)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, htmlContentWithActiveStatus, *res)
}
