package main

import (
	"encoding/json"
	"go-ansible-inventory/inventory"
	"go-ansible-inventory/zabbix"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	hosts := inventory.ListHosts{}

	ansibleEnv, _ := os.LookupEnv("ENV")

	if ansibleEnv != "" {
		envHosts := os.Getenv(ansibleEnv)

		hosts.AddFromEnv(envHosts)

	}

	zbx_server := os.Getenv("ZBX_SERVERS")

	if zbx_server != "" {
		zbx_user := os.Getenv("ZBX_USER")
		zbx_password := os.Getenv("ZBX_PASSWORD")
		zbx := zabbix.Zabbix{
			Url:      zabbix.SetZabbixUrl(zbx_server),
			User:     zbx_user,
			Password: zbx_password}

		err := hosts.AddFromZabbix(zbx)
		if err != nil {
			log.Fatal(err)
		}
	}
	ansibleInventory := hosts.CreateInventory()
	js, err := json.MarshalIndent(ansibleInventory, "", "  ")

	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(js)
}
