package inventory

import (
	"fmt"
	"go-ansible-inventory/zabbix"
	"strings"
)

type Inventory map[string]interface{}

func (inv *Inventory) NewInventory() Inventory {
	return make(Inventory)
}

func (inv Inventory) EnvInv(s string) error {
	var listServerName []string
	all := make(map[string]interface{})
	hostVars := make(map[string]interface{})

	ips := strings.Split(s, ",")
	for i, ip := range ips {
		if ip != "" {
			serverName := fmt.Sprintf("envHost%d", i)
			listServerName = append(listServerName, serverName)
			all["EnvServers"] = listServerName
			hostVars[serverName] = map[string]string{"ansible_host": strings.TrimSpace(ip)}
		}
	}

	inv["_meta"] = hostVars
	inv["all"] = all

	return nil
}

func (inv Inventory) ZabbixInventory(zbx zabbix.Zabbix) error {
	hostVars := make(map[string]interface{}, 0)
	all := make(map[string]interface{}, 0)
	ungrouped := make(map[string]interface{}, 0)
	children := make([]string, 0)
	allHosts := make([]string, 0)
	groups := make(map[string][]interface{}, 0)

	defer zbx.Logout()

	zbx.NewSession()
	hosts, err := zbx.GetHosts()

	if err != nil {
		return err
	}

	for _, host := range hosts.Result {

		hostVars[host.Name] = map[string]string{"ansible_host": host.Interfaces[0].IP}

		allHosts = append(allHosts, strings.ReplaceAll(host.Name, " ", ""))

		for _, group := range host.Groups {
			groupName := strings.ReplaceAll(group.Name, " ", "")
			groups[groupName] = append(groups[groupName], host.Name)
		}

	}

	for groupName, hostName := range groups {
		inv[groupName] = hostName
		children = append(children, groupName)
	}

	all["hosts"] = allHosts
	all["children"] = children
	all["ungrouped"] = ungrouped

	inv["_meta"] = map[string]interface{}{"hostvars": hostVars}
	inv["all"] = all

	return nil

}
