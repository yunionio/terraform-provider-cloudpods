package yunion

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/utils"

	"strings"

	"yunion.io/x/log"
	"yunion.io/x/onecloud/pkg/mcclient/modules"
	"yunion.io/x/onecloud/pkg/util/httputils"
)

func resourceYunionServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceYunionServerCreate,
		Read:   resourceYunionServerRead,
		Update: resourceYunionServerUpdate,
		Delete: resourceYunionServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ncpu": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vmem": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"hypervisor": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateServerHypervisor,
			},
			"data_disks": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems:      12,
				MinItems:      0,
				PromoteSingle: false,
				Optional:      true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func validateServerHypervisor(v interface{}, k string) (ws []string, errs []error) {
	value := v.(string)
	if utils.IsInStringArray(value, []string{"kvm", "baremetal", "esxi", "aliyun", "azure"}) {
		errs = append(errs, fmt.Errorf("Invalid %s %s", k, value))
	}
	return
}

func resourceYunionServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SYunionClient)

	s := client.getSession("v2")

	params := jsonutils.NewDict()

	params.Add(jsonutils.NewString(d.Get("name").(string)), "name")
	params.Add(jsonutils.NewString(d.Get("vmem").(string)), "vmem_size")
	params.Add(jsonutils.NewString(d.Get("image_id").(string)), "disk.0")
	if v, ok := d.GetOk("data_disks"); ok {
		disks := v.([]string)
		for i, d := range disks {
			params.Add(jsonutils.NewString(d), fmt.Sprintf("disk.%d", i+1))
		}
	}
	if v, ok := d.GetOk("ncpu"); ok {
		params.Add(jsonutils.NewInt(int64(v.(int))), "ncpu_count")
	}
	if v, ok := d.GetOk("hypervisor"); ok {
		params.Add(jsonutils.NewString(v.(string)), "hypervisor")
	}
	if v, ok := d.GetOk("description"); ok {
		params.Add(jsonutils.NewString(v.(string)), "description")
	}

	params.Add(jsonutils.JSONFalse, "disable_delete")

	params.Add(jsonutils.JSONTrue, "auto_start")

	server, err := modules.Servers.Create(s, params)
	if err != nil {
		log.Errorf("fail to create server %s", err)
		return err
	}

	id, _ := server.GetString("id")

	if len(id) == 0 {
		return fmt.Errorf("fail to find id??")
	}

	d.SetId(id)

	wait := 0 * time.Second
	waitInterval := 10 * time.Second
	maxWait := 30 * time.Minute

	for wait < maxWait {
		server, err = modules.Servers.Get(s, id, nil)
		if err != nil {
			log.Errorf("fail to get server %s", err)
			return err
		}

		status, _ := server.GetString("status")
		if strings.HasSuffix(status, "fail") {
			log.Errorf("failed status %s", status)
			return fmt.Errorf("server status fail %s", status)
		}

		if status == "running" {
			return nil
		}

		time.Sleep(waitInterval) // wait 15 seconds and try again

		wait += waitInterval
	}

	log.Errorf("timeout")
	return fmt.Errorf("timeout")
}

type yunionServer struct {
	Id          string
	Name        string
	VmemSize    int // MB
	VcpuCount   int
	Description string
	Hypervisor  string
}

func resourceYunionServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SYunionClient)

	s := client.getSession("v2")

	obj, err := modules.Servers.Get(s, d.Id(), nil)
	if err != nil {
		return err
	}

	srv := yunionServer{}
	err = obj.Unmarshal(&srv)
	if err != nil {
		return err
	}

	d.Set("name", srv.Name)
	d.Set("vmem", srv.VmemSize)
	d.Set("ncpu", srv.VcpuCount)
	d.Set("hypervisor", srv.Hypervisor)
	d.Set("description", srv.Description)
	return nil
}

func resourceYunionServerUpdate(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*SYunionClient)

	// s := client.getSession("v2")

	return nil
}

func isNotFoundError(err error) bool {
	jsonErr, ok := err.(*httputils.JSONClientError)
	if !ok {
		return false
	}
	if jsonErr.Code != 404 {
		return false
	}
	return true
}

func resourceYunionServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SYunionClient)

	s := client.getSession("v2")

	params := jsonutils.NewDict()
	params.Add(jsonutils.JSONTrue, "override_pending_delete")

	_, err := modules.Servers.Delete(s, d.Id(), params)

	if err != nil {
		if isNotFoundError(err) {
			return nil
		} else {
			return err
		}
	}

	waited := 0 * time.Second
	waitInterval := 5 * time.Second
	maxWait := 10 * time.Minute

	for waited < maxWait {
		_, err := modules.Servers.Get(s, d.Id(), nil)
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
		}

		time.Sleep(waitInterval)
		waited += waitInterval
	}
	return fmt.Errorf("delete timeout")
}
