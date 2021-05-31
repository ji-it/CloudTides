package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"tides-server/pkg/config"
	"tides-server/pkg/models"
)

func main() {
	config.GetConfig()
	db := config.GetDB()
	ipaddr := get_internal()
	fmt.Println(ipaddr)
	var VM models.VMachine
	var VAPP models.Vapp
	db.Where("ip_address = ?", ipaddr).First(&VM)
	db.Preload("VMs").Where("id = ?", VM.VappID).First(&VAPP)
	lineList := "party_list=("
	lineIPList := "party_ip_list=("
	lineServing := "serving_ip_list=("
	for index, vm := range VAPP.VMs {
		if(vm.IPAddress != ipaddr) {
			lineList += (strconv.Itoa(10000 - index) + " ")
			lineIPList += (vm.IPAddress + " ")
			lineServing += (vm.IPAddress + " ")
		}
	}
	lineList += ")"
	lineIPList += ")"
	lineServing += ")"
	fi, err := os.OpenFile("/root/docker-deploy/parties.conf", os.O_RDWR, os.ModePerm)
	//fi, err := os.OpenFile("/home/frank/cloudTides/test", os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fi.Close()
	reader := bufio.NewReader(fi)
	lineCnt := 0
	seekP := 0
	for {
		bs, _, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Println("Done")
			return
		}
		lineCnt = len(bs) + 1
		if strings.Contains(string(bs), "party_list") {
			delByte := make([]byte, 0)
			for i := 0; i < lineCnt; i++ {
				delByte = append(delByte, 127)
			}

			_, err := fi.WriteAt(delByte, int64(seekP))
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = fi.WriteAt([]byte(lineList + "\n"), int64(seekP))
			if err != nil {
				fmt.Println(err)
				return
			}
			lineCnt = len([]byte((lineList + "\n")))
		} else if strings.Contains(string(bs), "party_ip_list") {
			delByte := make([]byte, 0)
			for i := 0; i < lineCnt; i++ {
				delByte = append(delByte, 127)
			}

			_, err := fi.WriteAt(delByte, int64(seekP))
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = fi.WriteAt([]byte(lineIPList + "\n"), int64(seekP))
			if err != nil {
				fmt.Println(err)
				return
			}
			lineCnt = len([]byte((lineIPList + "\n")))
		} else if strings.Contains(string(bs), "serving_ip_list") {
			delByte := make([]byte, 0)
			for i := 0; i < lineCnt; i++ {
				delByte = append(delByte, 127)
			}

			_, err := fi.WriteAt(delByte, int64(seekP))
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = fi.WriteAt([]byte(lineServing + "\n"), int64(seekP))
			if err != nil {
				fmt.Println(err)
				return
			}
			lineCnt = len([]byte((lineServing + "\n")))
		}

		seekP += lineCnt
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
