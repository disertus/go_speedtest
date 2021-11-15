package main

import (
	"fmt"
	"encoding/json"
	"os"
	"os/exec"
	"time"
	"log"
)

type Speedtest struct {
	Type      string
	Timestamp time.Time
	Ping      struct {
		Jitter  float64
		Latency float64
	} `json:"ping"`
	Download struct {
		Bandwidth int
		Bytes     int
		Elapsed   int
	} `json:"download"`
	Upload struct {
		Bandwidth int
		Bytes     int
		Elapsed   int
	} `json:"upload"`
	PacketLoss float64
	Isp        string
	Interface  struct {
		InternalIp string
		Name       string
		MacAddr    string
		IsVpn      bool
		ExternalIp string
	} `json:"interface"`
	Server struct {
		Id       int
		Name     string
		Location string
		Country  string
		Host     string
		Port     int
		Ip       string
	} `json:"server"`
	Result struct {
		Id  string
		Url string
	} `json:"result"`
} 



func main() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	output, err := exec.Command(path+"/speedtest", "-f", "json-pretty").Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(output))

	var data Speedtest
	err = json.Unmarshal(output, &data)
	if err != nil {
		log.Fatalf("Error occurred while unmarshalling: %s", err.Error())
	}

	// fmt.Println(data.Download.Bandwidth * 8 / 1000000)
	// fmt.Println(data.Upload.Bandwidth * 8 / 1000000)
	// fmt.Println(data.Result.Url)

	// time.Sleep(20 * time.Second)
}
