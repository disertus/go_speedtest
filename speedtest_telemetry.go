package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/exec"
	"time"
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

const (
	host     = "localhost"
	port     = 5432
	user     = USER
	password = PASS
	dbname   = DBNAME
)

func insertIntoDB(db *sql.DB, data Speedtest) {
	sqlStatement := `
			INSERT INTO speedtests (id, created_at, type, jitter, latency, download_bandwidth, download_bytes, download_elapsed,
				upload_bandwidth, upload_bytes, upload_elapsed, packet_loss, isp, internal_id, interface_name, mac_address, is_vpn,
				external_ip, server_id, server_name, server_location, server_country, server_host, server_port, server_ip, speedtest_url)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)`

	_, err := db.Exec(sqlStatement, data.Result.Id, data.Timestamp, data.Type, data.Ping.Jitter, data.Ping.Latency,
		data.Download.Bandwidth, data.Download.Bytes, data.Download.Elapsed, data.Upload.Bandwidth, data.Upload.Bytes,
		data.Upload.Elapsed, data.PacketLoss, data.Isp, data.Interface.InternalIp, data.Interface.Name, data.Interface.MacAddr,
		data.Interface.IsVpn, data.Interface.ExternalIp, data.Server.Id, data.Server.Name, data.Server.Location, data.Server.Country,
		data.Server.Host, data.Server.Port, data.Server.Ip, data.Result.Url)
	if err != nil {
		panic(err)
	}

	return
}

func runSpeedtest(path string) []byte {
	output, err := exec.Command(path+"/speedtest", "-f", "json-pretty").Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	return output
}

func scheduleSpeedtest() {
	panic("Not implemented")
}

func main() {
	for {
		path, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}

		var output []byte
		output = runSpeedtest(path)

		fmt.Println(string(output))

		var data Speedtest
		err = json.Unmarshal(output, &data)
		if err != nil {
			log.Fatalf("Error occurred while unmarshalling: %s", err.Error())
		}

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}

		defer db.Close()

		insertIntoDB(db, data)

		if err != nil {
			panic(err)
		}

		fmt.Println("Successfully connected!")
		fmt.Println("Sleeping for 20 secs")
		time.Sleep(1800 * time.Second)
	}

}
