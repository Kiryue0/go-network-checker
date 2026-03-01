package network

import (
	"errors"
	"net"

	"github.com/Kiryue0/go-network-checker/internal/model"
)

func GetInterfaces() ([]model.InterfaceInfo, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.New("failed to get network interfaces")
	}

	var results []model.InterfaceInfo

	for _, iface := range ifaces {
		info := model.InterfaceInfo{}

		info.Name = iface.Name
		info.MTU = iface.MTU
		info.MACAddress = iface.HardwareAddr.String()
		info.IsUp = iface.Flags&net.FlagUp != 0

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		var ipAddr string
		for _, a := range addrs {
			ip, _, err := net.ParseCIDR(a.String())
			if err != nil {
				continue
			}
			if ip.To4() != nil {
				ipAddr = ip.String()
				break
			}
		}
		info.IPAddress = ipAddr

		results = append(results, info)
	}

	return results, nil
}
