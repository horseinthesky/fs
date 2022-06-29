package parser

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"
)

func GetDataFromFile(fileName string) (*FlowSpecData, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	filterData := &FlowSpecData{}

	if err := yaml.Unmarshal(content, filterData); err != nil {
		return nil, err
	}

	return filterData, nil
}

func GetDataFromWeb(url string, token string) (*FlowSpecData, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, _ := http.NewRequest("GET", url, nil)
	authHeader := "x-oauth-token " + token
	req.Header.Set("Authorization", authHeader)

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	filterData := &FlowSpecData{}

	if err := yaml.Unmarshal(body, &filterData); err != nil {
		return nil, err
	}

	return filterData, nil
}
