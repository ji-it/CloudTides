package controller

import (
	"encoding/json"
	"io/ioutil"
	"tides-server/pkg/config"
	"tides-server/pkg/models"

	"github.com/robfig/cron/v3"
)

func InitController() {
	cronjobs = map[uint]*cron.Cron{}
	// Query usage every 5 mins
	c := cron.New()
	c.AddFunc(schedule, func() {
		InitJob()
	})
	c.Start()
	// Clean up past usage table every week
	cl := cron.New()
	cl.AddFunc(cleanSchedule, func() {
		InitCleanUp()
	})
	cl.Start()
}

func InitJob() {
	db := config.GetDB()
	var resources []*models.Resource

	db.Where("activated = ? AND is_active = ? AND monitored = ?", true, true, true).Find(&resources)
	for _, res := range resources {
		_, ok := cronjobs[res.ID]
		if !ok {
			var vcd models.Vcd
			db.Where("resource_id = ?", res.ID).First(&vcd)
			newVcdConfig := VcdConfig{
				User:     res.Username,
				Password: res.Password,
				Org:      vcd.Organization,
				Href:     res.HostAddress,
				VDC:      res.Datacenter,
			}
			filename := "./cloudtides-" + res.Datacenter + ".json"
			file, _ := json.MarshalIndent(newVcdConfig, "", "")
			ioutil.WriteFile(filename, file, 0644)
			c := cron.New()
			c.AddFunc(schedule, func() {
				RunJob(filename)
			})
			c.Start()
			cronjobs[res.ID] = c
		}
	}

	db.Where("activated = ? AND is_active = ? AND monitored = ?", true, true, false).Find(&resources)

	for _, res := range resources {
		if res.PlatformType == models.ResourcePlatformTypeVcd {
			var vcd models.Vcd
			db.Where("resource_id = ?", res.ID).First(&vcd)
			newVcdConfig := VcdConfig{
				User:     res.Username,
				Password: res.Password,
				Org:      vcd.Organization,
				Href:     res.HostAddress,
				VDC:      res.Datacenter,
			}
			filename := "./cloudtides-" + res.Datacenter + ".json"
			file, _ := json.MarshalIndent(newVcdConfig, "", "")
			ioutil.WriteFile(filename, file, 0644)

			c := cron.New()
			c.AddFunc(schedule, func() {
				RunJob(filename)
			})
			c.Start()
			cronjobs[res.ID] = c
			res.Monitored = true
			db.Save(&res)
		}
	}
}

func InitCleanUp() {
	db := config.GetDB()
	db.Unscoped().Delete(&models.ResourcePastUsage{})
}

func RemoveJob(ResID uint) {
	c, ok := cronjobs[ResID]
	if ok {
		c.Stop()
	}
}
