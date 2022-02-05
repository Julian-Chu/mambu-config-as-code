package mambu

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	client "github.com/Julian-Chu/MambuConfigurationAPI/configurationClient"
)

func dataSourceCustomFields() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomFieldsRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_field_sets": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						/*
							"description": &schema.Schema{
								Type:     schema.TypeString,
								Computed: true,
							},
							"type": &schema.Schema{
								Type:     schema.TypeString,
								Computed: true,
							},
							"available_for": &schema.Schema{
								Type:     schema.TypeString,
								Computed: true,
							},
							"custom_fields": &schema.Schema{
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"id": &schema.Schema{
											Type:     schema.TypeString,
											Computed: true,
										},

										"type": &schema.Schema{
											Type:     schema.TypeString,
											Computed: true,
										},
										"state": &schema.Schema{
											Type:     schema.TypeString,
											Computed: true,
										},
										"validation_rules": &schema.Schema{
											Type:     schema.TypeSet,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"unique": &schema.Schema{
														Type:     schema.TypeBool,
														Computed: true,
													},
													"validation_pattern": &schema.Schema{
														Type:     schema.TypeString,
														Computed: true,
													},
												},
											},
										},
										"display_settings": &schema.Schema{
											Type:     schema.TypeSet,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"display_name": &schema.Schema{
														Type:     schema.TypeString,
														Computed: true,
													},
													"description": &schema.Schema{
														Type:     schema.TypeString,
														Computed: true,
													},
													"field_size": &schema.Schema{
														Type:     schema.TypeString,
														Computed: true,
													},
												},
											},
										},
										"usage": &schema.Schema{
											Type:     schema.TypeList,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"id": &schema.Schema{
														Type:     schema.TypeString,
														Computed: true,
													},
													"required": &schema.Schema{
														Type:     schema.TypeBool,
														Computed: true,
													},
													"default": &schema.Schema{
														Type:     schema.TypeBool,
														Computed: true,
													},
												},
											},
										},
										"view_rights": &schema.Schema{
											Type:     schema.TypeSet,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"roles": &schema.Schema{
														Type:     schema.TypeList,
														Computed: true,
														Elem: &schema.Schema{
															Type: schema.TypeString,
														},
													},
													"all_users": &schema.Schema{
														Type:     schema.TypeBool,
														Computed: true,
													},
												},
											},
										},
										"edit_rights": &schema.Schema{
											Type:     schema.TypeSet,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"roles": &schema.Schema{
														Type:     schema.TypeList,
														Computed: true,
														Elem: &schema.Schema{
															Type: schema.TypeString,
														},
													},
													"all_users": &schema.Schema{
														Type:     schema.TypeBool,
														Computed: true,
													},
												},
											},
										},
									},
								},
							},*/
					},
				},
			},
		},
	}
}

func dataSourceCustomFieldsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.MambuConfigClient)

	var diags diag.Diagnostics

	customFieldsResponse, err := c.GetCustomFields()
	if err != nil {
		return diag.FromErr(err)
	}
	bs, _ := json.Marshal(customFieldsResponse)
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  string(bs),
		Detail:   "response",
	})

	customFieldSets := flattenCustomFieldsData(&customFieldsResponse.CustomFieldSets)
	bs, _ = json.Marshal(customFieldSets)
	//if err := d.Set("custom_field_sets", interface{}(
	//	[]map[string]interface{}{{
	//		"error": "err",
	//		"id":    "id"}},
	//)); err != nil {
	results := make([]map[string]interface{}, 1)
	results[0] = make(map[string]interface{})
	results[0]["id"] = "test"
	results[0]["name"] = "test"
	if err := d.Set("custom_field_sets", results); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", "t"); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  string(bs),
		Detail:   "fields",
	})
	return diags
}

func flattenCustomFieldsData(sets *[]struct {
	ID           string `yaml:"id"`
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	Type         string `yaml:"type"`
	AvailableFor string `yaml:"availableFor"`
	CustomFields []struct {
		ID              string `yaml:"id"`
		Type            string `yaml:"type"`
		State           string `yaml:"state"`
		ValidationRules struct {
			Unique bool `yaml:"unique"`
		} `yaml:"validationRules"`
		DisplaySettings struct {
			DisplayName string `yaml:"displayName"`
			Description string `yaml:"description"`
			FieldSize   string `yaml:"fieldSize"`
		} `yaml:"displaySettings"`
		Usage []struct {
			ID       string `yaml:"id"`
			Required bool   `yaml:"required"`
			Default  bool   `yaml:"default"`
		} `yaml:"usage"`
		ViewRights struct {
			Roles    []interface{} `yaml:"roles"`
			AllUsers bool          `yaml:"allUsers"`
		} `yaml:"viewRights"`
		EditRights struct {
			Roles    []interface{} `yaml:"roles"`
			AllUsers bool          `yaml:"allUsers"`
		} `yaml:"editRights"`
	} `yaml:"customFields"`
}) interface{} {

	if sets != nil {
		fieldSets := make([]interface{}, len(*sets), len(*sets))

		for i, item := range *sets {
			cf := make(map[string]interface{})
			cf["id"] = item.ID
			//cf["custom_field_sets_id"] = item.ID
			cf["name"] = item.Name
			//cf["custom_field_sets_name"] = item.Name

			fieldSets[i] = cf
		}
		return fieldSets
	}
	return make([]interface{}, 0)
}
