package model

import (
	"time"
)

type InterfaceInfo struct {
	Name       string `json:"name"`
	IPAddress  string `json:"ip"`
	MACAddress string `json:"mac"`
	IsUp       string `json:"is_up"`
	MTU        int    `json:"mtu"`
}

type PingResult struct {
	Host       string        `json:"host"`
	IsAlive    bool          `json:"is_alive"`
	RTT        time.Duration `json:"rtt"`
	PacketLoss float64       `json:"packet_loss"`
	Timestamp  time.Time     `json:"timestamp"`
}

type PortResult struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	IsOpen       bool          `json:"is_open"`
	Service      string        `json:"service"`
	ResponseTime time.Duration `json:"response_time"`
}

type ScanReport struct {
	ScanDate    time.Time    `json:"scan_date"`
	TotalHost   int          `json:"total_host"`
	AliveHost   int          `json:"alive_host"`
	Results     []PingResult `json:"results"`
	PortResults []PortResult `json:"port_results"`
}
