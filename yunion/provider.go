package yunion

import (
	"os"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"yunion.io/x/log"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", os.Getenv("OS_AUTH_URL")),
				Description: "Keystone Auth URL",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USERNAME", os.Getenv("OS_USERNAME")),
				Description: "User account",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSWORD", os.Getenv("OS_PASSWORD")),
				Description: "User password",
			},
			"domain": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DOMAIN_NAME", os.Getenv("OS_DOMAIN_NAME")),
				Description: "user domain domain",
			},
			"project": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PROJECT_NAME", os.Getenv("OS_PROJECT_NAME")),
				Description: "project name",
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION_NAME", os.Getenv("OS_REGION_NAME")),
				Description: "name of cloud/region",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5 * time.Minute,
				Description: "request timeout",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"yunion_servers": dataSourceYunionServers(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"yunion_server": resourceYunionServer(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	conf := SYunionConfig{}

	if v, ok := d.GetOk("auth_url"); ok {
		conf.AuthUrl = v.(string)
	}

	if v, ok := d.GetOk("timeout"); ok {
		conf.Timeout = v.(int)
	}

	if v, ok := d.GetOk("username"); ok {
		conf.Username = v.(string)
	}

	if v, ok := d.GetOk("password"); ok {
		conf.Password = v.(string)
	}

	if v, ok := d.GetOk("domain"); ok {
		conf.Domain = v.(string)
	}

	if v, ok := d.GetOk("project"); ok {
		conf.Project = v.(string)
	}

	if v, ok := d.GetOk("region"); ok {
		conf.Region = v.(string)
	}

	conf.Insecure = true
	conf.DefaultApiVersion = ""
	conf.EndpointType = "internal"

	client, err := NewYunionClient(conf)
	if err != nil {
		log.Errorf("Fail to authenticate")
		return nil, err
	}

	return client, nil
}
