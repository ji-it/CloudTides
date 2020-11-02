package controller

import (
	"encoding/json"
	"io/ioutil"
	"tides-server/pkg/config"
	"tides-server/pkg/models"

	"github.com/robfig/cron/v3"
)

func InitController() {
	c := cron.New()
	c.AddFunc("*/2 * * * *", func() {
		InitializeJob()
	})
	c.Start()
}

func InitializeJob() {
	db := config.GetDB()
	var resources []*models.Resource

	db.Where("monitored = ?", false).Find(&resources)

	for _, res := range resources {
		if res.IsActive && res.PlatformType == models.ResourcePlatformTypeVcd {
			var vcd models.Vcd
			db.Where("resource_id = ?", res.ID).First(&vcd)
			newVcdConfig := VcdConfig{
				User:     res.Username,
				Password: res.Password,
				Org:      vcd.Organization,
				Href:     res.HostAddress,
				VDC:      res.Datacenter,
			}
			filename := "../pkg/controller/cloudtides-" + res.Datacenter + ".json"
			file, _ := json.MarshalIndent(newVcdConfig, "", "")
			ioutil.WriteFile(filename, file, 0644)

			c := cron.New()
			c.AddFunc("*/2 * * * *", func() {
				RunJob(filename)
			})
			c.Start()
			cronjobs[res.ID] = c
			res.Monitored = true
			db.Save(&res)
		} else if !res.IsActive {
			c, ok := cronjobs[res.ID]
			if ok {
				c.Stop()
			}
		}
	}
}
