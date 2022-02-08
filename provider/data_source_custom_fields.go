package provider

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	client "github.com/Julian-Chu/MambuConfigurationAPI/configurationClient/rest"
)

func dataSourceCustomFields() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomFieldsRead,
		Schema: map[string]*schema.Schema{
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
										Type:     schema.TypeMap,
										Computed: true,
									},
									"display_settings": &schema.Schema{
										Type:     schema.TypeMap,
										Computed: true,
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
										Type:     schema.TypeList,
										MaxItems: 1,
										Required: true,
										Elem: &schema.Resource{Schema: map[string]*schema.Schema{
											"all_users": {
												Type:     schema.TypeBool,
												Computed: true,
											},
											"roles": {
												Type:     schema.TypeList,
												Computed: true,
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										}},
									},
									"edit_rights": &schema.Schema{
										Type:     schema.TypeList,
										MaxItems: 1,
										Required: true,
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
										//Set: schema.HashResource(schemaEditRights()),
									},
								},
							},
						},
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

	customFieldSets := flattenCustomFieldSets(&customFieldsResponse.CustomFieldSets)
	bs, _ = json.Marshal(customFieldSets)
	if err := d.Set("custom_field_sets", customFieldSets); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("custom_field_sets")
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  string(bs),
		Detail:   "fields",
	})
	return diags
}

func flattenCustomFieldSets(sets *[]client.CustomFieldSet) interface{} {
	if sets != nil {
		fieldSets := make([]interface{}, len(*sets), len(*sets))

		for i, item := range *sets {
			cf := make(map[string]interface{})
			cf["id"] = item.ID
			cf["name"] = item.Name
			cf["description"] = item.Description
			cf["type"] = item.Type
			cf["available_for"] = item.AvailableFor
			// todo custom fields
			cf["custom_fields"] = flattenCustomFields(&item.CustomFields)

			fieldSets[i] = cf
		}
		return fieldSets
	}
	return make([]interface{}, 0)
}

func flattenCustomFields(fields *[]client.CustomField) interface{} {
	if fields == nil {
		return make([]interface{}, 0)
	}
	customFields := make([]interface{}, len(*fields), len(*fields))
	for i, field := range *fields {
		cf := make(map[string]interface{})
		cf["id"] = field.ID
		cf["type"] = field.Type
		cf["state"] = field.State
		cf["validation_rules"] = map[string]interface{}{
			"unique": strconv.FormatBool(field.ValidationRules.Unique),
		}
		cf["display_settings"] = map[string]interface{}{
			"displayName": field.DisplaySettings.DisplayName,
			"description": field.DisplaySettings.Description,
			"fieldSize":   field.DisplaySettings.FieldSize,
		}
		cf["usage"] = flattenUsage(&field.Usage)
		cf["view_rights"] = flattenRights(&field.ViewRights)
		cf["edit_rights"] = flattenRights(&field.EditRights)

		customFields[i] = cf
	}
	return customFields
}

func flattenRights(rs *client.Rights) interface{} {
	if rs == nil {
		return []interface{}{}
	}

	right := make(map[string]interface{})
	right["all_users"] = rs.AllUsers
	right["roles"] = rs.Roles
	return []interface{}{right}
}

func flattenUsage(items *[]client.Usage) interface{} {
	if items == nil {
		return make([]interface{}, 0)
	}
	usages := make([]interface{}, len(*items), len(*items))
	for i, item := range *items {
		usage := make(map[string]interface{})
		usage["id"] = item.ID
		usage["required"] = item.Required
		usage["default"] = item.Default

		usages[i] = usage
	}
	return usages
}
