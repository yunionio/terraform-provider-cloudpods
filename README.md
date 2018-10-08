Terraform Yunion Provider
=========================

- Website: https://www.terraform.io


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
