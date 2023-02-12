package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Zabbix struct {
	Url      string
	User     string
	Password string
	Token    string `json:"result"`
	Client   http.Client
}

type ZabbixResult struct {
	Result []struct {
		Hostid string `json:"hostid"`
		Host   string `json:"host"`
		Name   string `json:"name"`
		Status string `json:"status"`
		Groups []struct {
			Groupid string `json:"groupid"`
			Name    string `json:"name"`
		} `json:"groups"`
		Interfaces []struct {
			IP string `json:"ip"`
		} `json:"interfaces"`
	} `json:"result"`
}

func SetZabbixUrl(url string) string {
	return url + "/api_jsonrpc.php"
}

func (z *Zabbix) NewSession() error {
	jsonZabbixLogin := fmt.Sprintf(`
	{
		"jsonrpc": "2.0",
		"method": "user.login",
		"params": {
			"user": "%s",
			"password": "%s"
		},
		"id": 1,
		"auth": null
	}`, z.User, z.Password)

	jsonData := bytes.NewBuffer([]byte(jsonZabbixLogin))

	req, err := http.NewRequest(http.MethodPost, z.Url, jsonData)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.Client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &z)

	if err != nil {
		return err
	}

	return nil
}

func (z *Zabbix) Logout() error {
	jsonLogout := fmt.Sprintf(`
	{
		"jsonrpc": "2.0",
		"method": "user.logout",
		"params": [],
		"id": 1,
		"auth": "%s" 
	}`, z.Token)

	jsonData := bytes.NewBuffer([]byte(jsonLogout))

	req, err := http.NewRequest("POST", z.Url, jsonData)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.Client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (z *Zabbix) GetHosts() (ZabbixResult, error) {
	zbxResult := ZabbixResult{}

	// "filter": {
	// 	"host": ["v220210856809159963", "v220200756809124041", "v2202301180000214316"]
	// }

	jsonGetHosts := fmt.Sprintf(`
	{
		"jsonrpc": "2.0",
		"method": "host.get",
		"params": {
			"output": ["host","name","available","status"],
			"selectInterfaces":["ip"],
			"selectGroups": ["name"]			
		},		
		"id": 1,
		"auth": "%s" 
	}`, z.Token)

	jsonData := bytes.NewBuffer([]byte(jsonGetHosts))

	req, err := http.NewRequest(http.MethodPost, z.Url, jsonData)
	if err != nil {
		return zbxResult, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.Client.Do(req)

	if err != nil {
		return zbxResult, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return zbxResult, err
	}

	json.Unmarshal(body, &zbxResult)

	//fmt.Println(zbxResult)

	return zbxResult, nil

}
