package yunion

import (
	"github.com/hashicorp/terraform/helper/schema"

	"bytes"
	"fmt"

	"github.com/hashicorp/terraform/helper/hashcode"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules"
)

func dataSourceYunionServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceYunionServersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"name": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disks": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hypervisor": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataResourceIdHash(ids []string) string {
	var buf bytes.Buffer

	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func dataSourceYunionServersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SYunionClient)

	s := client.getSession("v2")

	params := jsonutils.NewDict()

	offset := 0
	limit := 100

	var allServers []jsonutils.JSONObject
	for {
		params.Set("offset", jsonutils.NewInt(int64(offset)))
		params.Set("limit", jsonutils.NewInt(int64(limit)))

		results, err := modules.Servers.List(s, params)
		if err != nil {
			return err
		}

		allServers = append(allServers, results.Data...)

		offset += len(results.Data)
		if offset >= results.Total {
			break
		}
	}

	ids := make([]string, len(allServers))
	mappings := make([]map[string]interface{}, len(allServers))
	for i := 0; i < len(allServers); i += 1 {
		id, mapping := json2mapping(allServers[i])
		ids = append(ids, id)
		mappings = append(mappings, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("servers", mappings); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

func json2mapping(jsonData jsonutils.JSONObject) (string, map[string]interface{}) {
	mapping := jsonData.Interface().(map[string]interface{})
	return mapping["id"].(string), mapping
}
