package main

import (
	"github.com/hashicorp/terraform/plugin"
	"yunion.io/x/terraform-provider-yunion/yunion"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: yunion.Provider})
}
