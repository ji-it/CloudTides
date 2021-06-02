package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"tides-server/pkg/models"
	"time"
)

const (
	defaultCheckDur = 30 * time.Second
)

func main() {
	var dbinfo string
	dbinfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	db, err := gorm.Open(postgres.Open(dbinfo), &gorm.Config{})
	//defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Project{})
	db.AutoMigrate(&models.Template{})
	db.AutoMigrate(&models.VMachine{})
	db.AutoMigrate(&models.VMTemp{})
	db.AutoMigrate(&models.Policy{})
	db.AutoMigrate(&models.VcdPolicy{})
	db.AutoMigrate(&models.Resource{})
	db.AutoMigrate(&models.Vsphere{})
	db.AutoMigrate(&models.Vcd{})
	db.AutoMigrate(&models.VM{})
	db.AutoMigrate(&models.ResourceUsage{})
	db.AutoMigrate(&models.ResourcePastUsage{})
	db.AutoMigrate(&models.VMUsage{})
	db.AutoMigrate(&models.Vendor{})
	db.AutoMigrate(&models.Vapp{})
	fmt.Println("DB connection success")
	ipaddr := get_internal()
	fmt.Println(ipaddr)
	var VM models.VMachine
	var VAPP models.Vapp
	ticker := time.NewTicker(defaultCheckDur)
	checking:
	for {
		select {
		case <-ticker.C:
			db.Where("ip_address = ?", ipaddr).First(&VM)
			db.Where("id = ?", VM.VappID).First(&VAPP)
			if VAPP.Status == "Creating" {
				continue
			} else {
				break checking
			}
		}
	}
	db.Where("ip_address = ?", ipaddr).First(&VM)
	db.Preload("VMs").Where("id = ?", VM.VappID).First(&VAPP)
	lineList := "party_list=("
	lineIPList := "party_ip_list=("
	lineServing := "serving_ip_list=("
	count := 0
	for _, vm := range VAPP.VMs {
		if(vm.IPAddress != ipaddr) {
			lineList += (strconv.Itoa(10000 - count) + " ")
			lineIPList += (vm.IPAddress + " ")
			lineServing += (vm.IPAddress + " ")
			count++
		}
	}
	lineList += ")"
	lineIPList += ")"
	lineServing += ")"
	input, err := ioutil.ReadFile("/root/docker-deploy/parties.conf")
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, "party_list") {
			lines[i] = lineList
		} else if strings.Contains(line, "party_ip_list") {
			lines[i] = lineIPList
		} else if strings.Contains(line, "serving_ip_list") {
			lines[i] = lineServing
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("/root/docker-deploy/parties.conf", []byte(output), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func get_internal() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops:" + err.Error())
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
