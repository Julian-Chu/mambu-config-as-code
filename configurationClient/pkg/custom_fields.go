package pkg

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// GetCustomFields represents GET /api/configuration/customfields.yaml
func (c MambuConfigClient) GetCustomFields() (*CustomFieldsConfig, error) {
	bs, err := c.sendRequest(http.MethodGet, "api/configuration/customfields.yaml", nil, nil)
	if err != nil {
		return nil, err
	}
	var customFieldsResp CustomFieldsConfig
	err = yaml.Unmarshal(bs, &customFieldsResp)

	return &customFieldsResp, nil
}

type CustomFieldsConfig struct {
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
