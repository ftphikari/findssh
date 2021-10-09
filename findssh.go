package main

import (
	"fmt"
	"net"
	"sync"
	"time"
	"strings"
)

func main() {
	var wg sync.WaitGroup
	var lock sync.Mutex
	timeout := time.Duration(1 * time.Second)

	ping := func(host string) {
		defer wg.Done()

		c, err := net.DialTimeout("tcp", host+":22", timeout)
		if err == nil {
			c.Close()
			lock.Lock()
			fmt.Println(host)
			lock.Unlock()
		}
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return;
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			addr := strings.Split(ip.String(), ".")

			wg.Add(254)
			for i := 1; i < 255; i++ {
				go ping(fmt.Sprintf("%s.%s.%s.%d", addr[0], addr[1], addr[2], i))
			}
		}
	}

	wg.Wait()
}
