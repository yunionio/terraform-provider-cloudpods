Terraform Yunion Provider
=========================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.10 (to build the provider plugin)
-   [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports):
    ```
    go get golang.org/x/tools/cmd/goimports
    ```

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/yunion.io/x/terraform-provider-yunion`

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/yunion.io/x/terraform-provider-yunion
$ make dev
```

Using the provider
----------------------
## Fill in for each provider

Example terraform configuration:

```
provider "yunion" {
    auth_url = "http://10.168.222.251:5000/v3"
    username = "sysadmin"
    password = "${sysadmin_password}"
    domain = "Default"
    project = "system"
}

# Create a web server
resource "yunion_server" "web" {
    name = "testkvm"
    ncpu = 1
    vmem = "1g"
    image_id = "89d91f52-19e3-4697-a875-b6d961ffd9a8"
}
```