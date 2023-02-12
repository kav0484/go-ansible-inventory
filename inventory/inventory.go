package inventory

import (
	"fmt"
	"go-ansible-inventory/zabbix"
	"strings"
)

type ListHosts struct {
	Hosts []Host
}

type Host struct {
	Name   string
	Groups []string
	IP     string
}

func (lh *ListHosts) AddFromEnv(s string) error {
	host := Host{}

	ips := strings.Split(s, ",")
	for i, ip := range ips {
		host.Groups = nil
		if ip != "" {
			host.Name = fmt.Sprintf("envHost%d", i)
			host.Groups = append(host.Groups, "EnvServers")
			host.IP = strings.TrimSpace(ip)
		}
		lh.Hosts = append(lh.Hosts, host)
	}

	return nil
}

func (lh *ListHosts) AddFromZabbix(zbx zabbix.Zabbix) error {
	host := Host{}

	zbx.NewSession()
	defer zbx.Logout()
	hosts, err := zbx.GetHosts()

	if err != nil {
		return err
	}

	for _, zbxHost := range hosts.Result {
		host.Groups = nil
		host.Name = strings.ReplaceAll(zbxHost.Name, " ", "")
		host.IP = zbxHost.Interfaces[0].IP

		for _, zbxGroup := range zbxHost.Groups {
			groupName := strings.ReplaceAll(zbxGroup.Name, " ", "")
			host.Groups = append(host.Groups, groupName)
		}
		lh.Hosts = append(lh.Hosts, host)
	}

	return nil
}

func (lh *ListHosts) CreateInventory() map[string]interface{} {
	inventory := make(map[string]interface{})

	hostVars := make(map[string]interface{}, 0)
	all := make(map[string]interface{}, 0)
	ungrouped := make(map[string]interface{}, 0)
	children := make([]string, 0)
	allHosts := make([]string, 0)
	groups := make(map[string][]string, 0)

	for _, host := range lh.Hosts {
		hostVars[host.Name] = map[string]string{"ansible_host": host.IP}
		allHosts = append(allHosts, host.Name)

		for _, group := range host.Groups {
			groups[group] = append(groups[group], host.Name)
		}

	}

	inventory["_meta"] = map[string]interface{}{"hostvars": hostVars}

	for groupName, hosts := range groups {
		inventory[groupName] = map[string][]string{"hosts": hosts}
		children = append(children, groupName)
	}

	all["hosts"] = allHosts
	all["children"] = children
	inventory["all"] = all

	inventory["ungrouped"] = ungrouped

	return inventory
}
