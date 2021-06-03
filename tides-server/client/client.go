package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"tides-server/pkg/models"
	"time"
)

const (
	defaultCheckDur = 10 * time.Second
)

func main() {
	f, err := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	var dbinfo string
	dbinfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	ticker := time.NewTicker(defaultCheckDur)
	var db *gorm.DB
	waiting:
	for {
		select {
		case <-ticker.C:
			db, err = gorm.Open(postgres.Open(dbinfo), &gorm.Config{})
			//defer db.Close()
			if err != nil {
				log.Println(err)
				continue
			} else {
				break waiting
			}
		}
	}
	//db, err := gorm.Open(postgres.Open(dbinfo), &gorm.Config{})
	//defer db.Close()
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
	log.Println("DB connection success")
	ipaddr := get_internal()
	log.Println(ipaddr)
	var VAPP models.Vapp
	db.Preload("VMs").Where("id = ?", os.Args[6]).First(&VAPP)
	checking:
	for {
		select {
		case <-ticker.C:
			if VAPP.Status == "Creating" {
				log.Println("VAPP is still under Creating status")
				continue
			} else if VAPP.Status == "Running"{
				log.Println("Now VAPP is under status " + VAPP.Status)
				break checking
			} else {
				log.Println("Now VAPP is under status " + VAPP.Status)
				continue
			}
		}
	}
	lineList := "party_list=("
	lineIPList := "party_ip_list=("
	lineServing := "serving_ip_list=("
	count := 0
	for _, vm := range VAPP.VMs {
		if(vm.Name != "Deploy") {
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
		log.Println(err)
		return
	}
	log.Println("Here I'm going to run bash script")
	cmd := exec.Command("/bin/sh", "/root/test.sh")
	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return
	}
	err = cmd.Wait()
	if err != nil {
		log.Println(err)
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
