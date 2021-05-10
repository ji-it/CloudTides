package controller

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
	"tides-server/pkg/config"
	"tides-server/pkg/models"
)

// VcdConfig is json configuration for vCD resource
/*type VcdConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Org      string `json:"org"`
	Href     string `json:"href"`
	VDC      string `json:"vdc"`
	Token    string `json:"token"`
	Insecure bool   `json:"insecure"`
}*/

// Policy schema
type Policy struct {
	CPU float64 `json:"cpu"`
	RAM float64 `json:"ram"`
}

type VMMonitor struct {
	Task *cron.Cron
	VMID uint
	conf *config.VcdConfig
	Lock sync.Mutex
}

func NewVMMonitor(ID uint, Config *config.VcdConfig) *VMMonitor {
	c := cron.New()
	temp := VMMonitor{
		VMID: ID,
		Task: c,
		conf: Config,
		Lock: sync.Mutex{},
	}
	c.AddFunc(schedule, temp.CheckStatus)
	return &temp
}

type VAppMonitor struct {
	ID uint
	Lock sync.Mutex
}

func NewVAppMonitor(ID uint) *VAppMonitor {
	temp := VAppMonitor{
		ID: ID,
		Lock: sync.Mutex{},
	}
	return &temp
}

func StartVMMonitor(ID uint) bool {
	monitor, ok := VMMonitors.Load(ID)
	if ok {
		monitor.Task.Start()
	}
	return ok
}

func (vm *VMMonitor) CheckStatus() {
	fmt.Println("The monitor is working well")
	db := config.GetDB()
	var vmachine = models.VMachine{}
	if db.Where("id = ?", vm.VMID).First(&vmachine).RowsAffected == 0 {
		return
	}
	var vapp = models.Vapp{}
	if db.Where("id = ?", vmachine.VappID).First(&vapp).RowsAffected == 0 {
		return
	}
	var resource = models.Resource{}
	if db.Where("id = ?", vapp.ResourceID).First(&resource).RowsAffected == 0 {
		return
	}
	var vcd models.Vcd
	if db.Where("resource_id = ?", resource.ID).First(&vcd).RowsAffected == 0{
		return
	}
	client, err := vm.conf.Client()
	if err != nil {
		fmt.Println(err)
		return
	}
	org, err := client.GetOrgByName(vm.conf.Org)
	if err != nil {
		fmt.Println(err)
		return
	}
	vdc, err := org.GetVDCByName(vm.conf.VDC, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	Vapp, err := vdc.GetVAppByName(vapp.Name, true)
	VM, err := Vapp.GetVMByName(vmachine.Name, true)
	status, err := VM.GetStatus()
	if err != nil {
		fmt.Println(err)
		return
	}
	vmachine.Status = status
	fmt.Println("The status of the vm is: " + status)
	db.Save(&vmachine)
}

var (
	letters  = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	cronjobs map[uint]*cron.Cron
	VMMonitors *MoniterStore
	VappMonitors *MoniterStore
)

type MoniterStore struct {
	m sync.Map
}

func NewMonitorStore() *MoniterStore {
	ms := MoniterStore{
		m: sync.Map{},
	}
	return &ms
}

func (ms *MoniterStore) Store(key uint, value *VMMonitor) {
	ms.m.Store(key, value)
}

func (ms *MoniterStore) StoreVapp(key uint, value *VAppMonitor) {
	ms.m.Store(key, value)
}

func (ms *MoniterStore) Load(key uint) (value *VMMonitor, ok bool) {
	v, ok := ms.m.Load(key)
	if v != nil {
		value = v.(*VMMonitor)
	}
	return
}

func (ms *MoniterStore) LoadVapp(key uint) (value *VAppMonitor, ok bool) {
	v, ok := ms.m.Load(key)
	if v != nil {
		value = v.(*VAppMonitor)
	}
	return
}

func (ms *MoniterStore) Delete(key uint) {
	ms.m.Delete(key)
}

func (ms *MoniterStore) Range(f func(key, value interface{}) bool) {
	ms.m.Range(f)
}

const (
	schedule      string = "*/5 * * * *"
	cleanSchedule string = "0 0 * * 0"
)
