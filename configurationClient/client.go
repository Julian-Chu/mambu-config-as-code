package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"
)

type MambuConfigClient struct {
	client   http.Client
	mambuURL string
	apikey   string
}

const acceptHeader = "application/vnd.mambu.v2+yaml"

func NewClient(mambuURL string, apikey string) *MambuConfigClient {
	return &MambuConfigClient{mambuURL: mambuURL, apikey: apikey, client: http.Client{Timeout: 1 * time.Minute}}
}

// GetCustomFields represents GET /api/configuration/customfields.yaml
func (c MambuConfigClient) GetCustomFields() (*CustomFieldsResponse, error) {
	client := c.client

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/configuration/customfields.yaml", c.mambuURL), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", acceptHeader)
	request.Header.Add("ApiKey", c.apikey)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var customFieldsResp CustomFieldsResponse
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bs, &customFieldsResp)

	return &customFieldsResp, nil
}

type CustomFieldsResponse struct {
	CustomFieldSets []CustomFieldSet `yaml:"customFieldSets"`
}
type CustomFieldSet struct {
	ID           string        `yaml:"id"`
	Name         string        `yaml:"name"`
	Description  string        `yaml:"description"`
	Type         string        `yaml:"type"`
	AvailableFor string        `yaml:"availableFor"`
	CustomFields []CustomField `yaml:"customFields"`
}

type CustomField struct {
	ID              string `yaml:"id"`
	Type            string `yaml:"type"`
	State           string `yaml:"state"`
	ValidationRules struct {
		Unique bool `yaml:"unique"`
	} `yaml:"validationRules"`
	DisplaySettings DisplaySettings `yaml:"displaySettings"`
	Usage           []Usage         `yaml:"usage"`
	ViewRights      Rights          `yaml:"viewRights"`
	EditRights      Rights          `yaml:"editRights"`
}

type DisplaySettings struct {
	DisplayName string `yaml:"displayName"`
	Description string `yaml:"description"`
	FieldSize   string `yaml:"fieldSize"`
}
type Usage struct {
	ID       string `yaml:"id"`
	Required bool   `yaml:"required"`
	Default  bool   `yaml:"default"`
}

type Rights struct {
	Roles    []interface{} `yaml:"roles"`
	AllUsers bool          `yaml:"allUsers"`
}
