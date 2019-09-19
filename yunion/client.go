package yunion

import (
	"context"

	"yunion.io/x/log"

	"yunion.io/x/onecloud/pkg/mcclient"
)

type SYunionConfig struct {
	AuthUrl           string
	Username          string
	Password          string
	Domain            string
	Project           string
	ProjectDomain     string
	Region            string
	EndpointType      string
	DefaultApiVersion string
	Timeout           int
	Insecure          bool
	Debug             bool
}

type SYunionClient struct {
	client *mcclient.Client
	token  mcclient.TokenCredential
	conf   SYunionConfig
}

func NewYunionClient(conf SYunionConfig) (*SYunionClient, error) {
	cli := mcclient.NewClient(conf.AuthUrl, conf.Timeout, conf.Debug, conf.Insecure, "", "")
	token, err := cli.Authenticate(conf.Username, conf.Password, conf.Domain, conf.Project, conf.ProjectDomain)
	if err != nil {
		log.Errorf("authenticate fail %s", err)
		return nil, err
	}
	client := SYunionClient{
		client: cli,
		token:  token,
		conf:   conf,
	}
	return &client, nil
}

func (cli *SYunionClient) refreshToken() error {
	token, err := cli.client.Authenticate(cli.conf.Username, cli.conf.Password, cli.conf.Domain, cli.conf.Project, cli.conf.ProjectDomain)
	if err != nil {
		return err
	}
	cli.token = token
	return nil
}

func (cli *SYunionClient) getSession(apiVersion string) *mcclient.ClientSession {
	if !cli.token.IsValid() {
		err := cli.refreshToken()
		if err != nil {
			return nil
		}
	}
	if len(apiVersion) == 0 {
		apiVersion = cli.conf.DefaultApiVersion
	}
	return cli.client.NewSession(context.Background(), cli.conf.Region, "", cli.conf.EndpointType, cli.token, apiVersion)
}
