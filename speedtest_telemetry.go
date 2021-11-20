package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
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
	password = PASSWORD
	dbname   = DB_NAME
)

func InsertResultIntoDB(db *sql.DB, data Speedtest) {
	sqlStatement := `
			INSERT INTO speedtests (id, created_at, type, jitter, latency, download_bandwidth, download_bytes, download_elapsed,
				upload_bandwidth, upload_bytes, upload_elapsed, packet_loss, isp, internal_id, interface_name, mac_address, is_vpn,
				external_ip, server_id, server_name, server_location, server_country, server_host, server_port, server_ip, speedtest_url)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)`

	_, err := db.Exec(sqlStatement, data.Result.Id, data.Timestamp, data.Type, data.Ping.Jitter, data.Ping.Latency,
		data.Download.Bandwidth, data.Download.Bytes, data.Download.Elapsed, data.Upload.Bandwidth, data.Upload.Bytes,
		data.Upload.Elapsed, data.PacketLoss, data.Isp, data.Interface.InternalIp, data.Interface.Name, EncryptStringToMD5(data.Interface.MacAddr),
		data.Interface.IsVpn, data.Interface.ExternalIp, data.Server.Id, data.Server.Name, data.Server.Location, data.Server.Country,
		data.Server.Host, data.Server.Port, data.Server.Ip, data.Result.Url)
	if err != nil {
		log.Fatalf("Failed to insert data into the database: %s", err.Error())
	}

	return
}

func EncryptStringToMD5(mac_address string) string {
	stringToBytes := []byte(mac_address)
	hashedString := md5.Sum(stringToBytes)
	return hex.EncodeToString(hashedString[:]) // convert byte array to slice since it's impossible to stringify a byte array
}

func RunSpeedtest(path string) []byte {
	output, err := exec.Command(path+"/speedtest", "-f", "json-pretty").Output()
	if err != nil {
		log.Fatalf("Failed to run a command: %s", err.Error())
	}
	return output
}

func scheduleSpeedtest() {
	panic("Not implemented")
}

func main() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
	}

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(file)

	for {
		fmt.Println("Running a speedtest at ", time.Now().String())

		var output []byte
		output = RunSpeedtest(path)

		var data Speedtest
		err = json.Unmarshal(output, &data)
		if err != nil {
			log.Fatalf("Failed to unmarshal json: %s", err.Error())
		}

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatalf("Failed to establish a database connection: %s", err.Error())
		}

		defer db.Close()

		InsertResultIntoDB(db, data)

		fmt.Println("Successfully inserted the speedtest results into the db")
		fmt.Println("Sleeping for 15 minutes")

		time.Sleep(900 * time.Second)
	}

}
