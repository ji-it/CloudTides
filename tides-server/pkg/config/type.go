package config

import (
	"context"
	"fmt"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	"log"
	"log/syslog"
	"net/url"
	"tides-server/pkg/models"
)

const (
	defaultLoggingRemoteProtocol string          = "tcp"
	defaultLoggingFlag           int             = log.LstdFlags
	defaultLoggingPriority       syslog.Priority = syslog.LOG_INFO
	defaultLoggingTag            string          = "CloudTides-Server"
	defaultPort                  string          = "80"
)

// ContextRoot
var (
	ContextRoot = context.Background()
)

// Config consists fields to setup the CloudTides server
type Config struct {
	ServerIP         string `json:"serverIP"`
	ServerPort       string `json:"serverPort"`
	PostgresHost     string `json:"postgresHost"`
	PostgresPort     string `json:"postgresPort"`
	PostgresUser     string `json:"postgresUser"`
	PostgresPassword string `json:"postgresPassword"`
	PostgresDB       string `json:"postgresDB"`
	SecretKey        string `json:"secretKey"`
	AdminUser        string `json:"adminUser"`
	AdminPassword    string `json:"adminPassword"`
}

// VcdConfig is json configuration for vCD resource
type VcdConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Org      string `json:"org"`
	Href     string `json:"href"`
	VDC      string `json:"vdc"`
	Token    string `json:"token"`
	Insecure bool   `json:"insecure"`
}

// Client creates a vCD client
func (c *VcdConfig) Client() (*govcd.VCDClient, error) {
	u, err := url.ParseRequestURI(c.Href)
	if err != nil {
		return nil, fmt.Errorf("unable to pass url: %s", err)
	}

	vcdClient := govcd.NewVCDClient(*u, c.Insecure)
	if c.Token != "" {
		_ = vcdClient.SetToken(c.Org, govcd.AuthorizationHeader, c.Token)
	} else {
		_, err := vcdClient.GetAuthResponse(c.User, c.Password, c.Org)
		if err != nil {
			return nil, fmt.Errorf("unable to authenticate: %s", err)
		}
		// fmt.Printf("Token: %s\n", resp.Header[govcd.AuthorizationHeader])
	}
	return vcdClient, nil
}

func GetVcdConfig(vapp *models.Vapp) *VcdConfig {
	db := GetDB()
	var res models.Resource
	var vcd models.Vcd
	if db.Where("id = ?", vapp.ResourceID).First(&res).RowsAffected == 0 {
		return nil
	}
	if db.Where("resource_id = ?", res.ID).First(&vcd).RowsAffected == 0{
		return nil
	}
	conf := VcdConfig{
		Href: res.HostAddress,
		Password: res.Password,
		User: res.Username,
		Org: vcd.Organization,
		VDC: res.Datacenter,
	}
	return &conf
}
