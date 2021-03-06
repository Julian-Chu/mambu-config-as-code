package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	client "github.com/Julian-Chu/MambuConfigurationAPI/configurationClient/rest"
)

const KeyMambuBaseURL = "mambu_base_url"
const KeyMambuApiKey = "mambu_apikey"

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			KeyMambuBaseURL: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MAMBU_BASE_URL", nil),
			},
			KeyMambuApiKey: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("MAMBU_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"mambu_custom_fields": dataSourceCustomFields(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	mambuBaseURL := d.Get(KeyMambuBaseURL).(string)
	mambuApiKey := d.Get(KeyMambuApiKey).(string)

	var diags diag.Diagnostics

	if mambuBaseURL != "" && mambuApiKey != "" {
		c := client.NewClient(mambuBaseURL, mambuApiKey)
		return c, diags
	}

	if mambuBaseURL == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "provider base url is required",
			Detail:   "provider base url is required",
		})
	}
	if mambuApiKey == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "provider api key is required",
			Detail:   "provider api key is required",
		})
	}

	return nil, diags
}
