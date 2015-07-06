package mesoskafka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	url        *url.URL
	httpClient *http.Client
}

func NewClient(hostname string, port int) Client {
	return NewClientForUrl(fmt.Sprintf("https://%s:%d", hostname, port))
}

func NewClientForUrl(rawurl string) Client {
	url, err := url.Parse(rawurl)

	if err != nil {
		panic(err)
	}

	c := Client{
		url:        url,
		httpClient: &http.Client{},
	}

	return c
}

func (c *Client) getFullUrl(apiEndpoint string) string {
	fullUrl, err := c.url.Parse(apiEndpoint)
	if err != nil {
		panic(err)
	}
	return fullUrl.String()
}

func (c *Client) getJson(apiEndpoint string) ([]byte, error) {
	resp, err := c.httpClient.Get(c.getFullUrl(apiEndpoint))
	if err != nil {
		return nil, err
	}
	if statusCodeErr := checkSuccessfullStatusCode(resp); statusCodeErr != nil {
		return nil, statusCodeErr
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

func (c *Client) putJson(apiEndpoint string, json []byte) ([]byte, error) {
	req, err := http.NewRequest("PUT", c.getFullUrl(apiEndpoint), bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if statusCodeErr := checkSuccessfullStatusCode(resp); statusCodeErr != nil {
		return nil, statusCodeErr
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

func (c *Client) deleteJson(apiEndpoint string) ([]byte, error) {
	req, _ := http.NewRequest("DELETE", c.getFullUrl(apiEndpoint), nil)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if statusCodeErr := checkSuccessfullStatusCode(resp); statusCodeErr != nil {
		return nil, statusCodeErr
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

func checkSuccessfullStatusCode(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		defer resp.Body.Close()
		responseBody := buf.String()
		return fmt.Errorf("Returned HTTP status: %s \nReturned HTTP body: %s", resp.Status, responseBody)
	}
	return nil
}

type Status struct {
	FrameworkID string   `json:"frameworkId"`
	Brokers     []Broker `json:"brokers"`
}

type Broker struct {
	Id     string `json:"id"`
	Active bool   `json:"active"`
}

type Brokers struct {
	Brokers []Broker `json:"brokers"`
}

type MutateStatus struct {
	Status string `json:"started"`
}

func (c *Client) ApiBrokersStatus() (*Status, error) {
	body, e := c.getJson("/api/brokers/status")

	if e != nil {
		return nil, e
	}

	var status Status
	e = json.Unmarshal(body, &status)
	if e != nil {
		return nil, e
	}

	return &status, nil
}

func (c *Client) ApiBrokersAdd(BrokerId int) (*Brokers, error) {
	url := fmt.Sprintf("/api/brokers/add?id=%d&mem=256", BrokerId)
	body, e := c.getJson(url)

	if e != nil {
		return nil, e
	}

	var response Brokers
	e = json.Unmarshal(body, &response)
	if e != nil {
		return nil, e
	}

	return &response, nil
}

func (c *Client) ApiBrokersStart(BrokerId int) (*MutateStatus, error) {
	url := fmt.Sprintf("/api/brokers/start?id=%d", BrokerId)
	body, e := c.getJson(url)

	if e != nil {
		return nil, e
	}

	var response MutateStatus
	e = json.Unmarshal(body, &response)
	if e != nil {
		return nil, e
	}

	return &response, nil
}

func (c *Client) ApiBrokersStop(BrokerId int) (*MutateStatus, error) {
	url := fmt.Sprintf("/api/brokers/stop?id=%d", BrokerId)
	body, e := c.getJson(url)

	if e != nil {
		return nil, e
	}

	var response MutateStatus
	e = json.Unmarshal(body, &response)
	if e != nil {
		return nil, e
	}

	return &response, nil
}

func (c *Client) ApiBrokersRemove(BrokerId int) (*MutateStatus, error) {
	url := fmt.Sprintf("/api/brokers/remove?id=%d", BrokerId)
	body, e := c.getJson(url)

	if e != nil {
		return nil, e
	}

	var response MutateStatus
	e = json.Unmarshal(body, &response)
	if e != nil {
		return nil, e
	}

	return &response, nil
}

func (c *Client) ApiBrokersCreate(BrokerIds []int) error {

	for _, brokerId := range BrokerIds {
		_, err := c.ApiBrokersAdd(brokerId)

		if err != nil {
			return err
		}

		_, err = c.ApiBrokersStart(brokerId)

		if err != nil {
			return err
		}

	}

	return nil
}

func (c *Client) ApiBrokersDelete(BrokerIds []int) error {

	for brokerId, _ := range BrokerIds {
		_, err := c.ApiBrokersStop(brokerId)

		if err != nil {
			return err
		}

		_, err = c.ApiBrokersRemove(brokerId)

		if err != nil {
			return err
		}

	}

	return nil
}
