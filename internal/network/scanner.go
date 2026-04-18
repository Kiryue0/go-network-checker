package network

import (
	"context"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/Kiryue0/go-network-checker/internal/model"
)

var knownServices = map[int]string{
	22:    "SSH",
	80:    "HTTP",
	443:   "HTTPS",
	53:    "DNS",
	21:    "FTP",
	25:    "SMTP",
	3389:  "RDP",
	23:    "Telnet",
	3306:  "MySQL",
	5432:  "PostgreSQL",
	6379:  "Redis",
	8080:  "HTTP-Alt",
	27017: "MongoDB",
}

func ScanPort(ctx context.Context, host string, port int, timeout time.Duration) model.PortResult {

	start := time.Now()
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return model.PortResult{
			Host:    host,
			Port:    port,
			IsOpen:  false,
			Service: knownServices[port],
		}
	}

	elapsed := time.Since(start)
	defer conn.Close()

	return model.PortResult{
		Host:         host,
		Port:         port,
		IsOpen:       true,
		Service:      knownServices[port],
		ResponseTime: elapsed,
	}

}

func ScanPorts(ctx context.Context, hosts []string, ports []int, timeout time.Duration) []model.PortResult {
	var wg sync.WaitGroup
	results := make(chan model.PortResult, len(hosts)*len(ports))
	sem := make(chan struct{}, 50)
	for _, host := range hosts {
		for _, port := range ports {
			wg.Add(1)
			go func(h string, p int) {
				defer wg.Done()
				select {
				case sem <- struct{}{}:
				case <-ctx.Done():
					return
				}
				result := ScanPort(ctx, h, p, timeout)
				<-sem
				results <- result
			}(host, port)
		}
	}

	wg.Wait()
	close(results)

	var scanResults []model.PortResult
	for result := range results {
		scanResults = append(scanResults, result)
	}

	sort.Slice(scanResults, func(i, j int) bool {
		return scanResults[i].Port < scanResults[j].Port
	})

	return scanResults
}
