package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"sort"
)

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func nextRangeFirstIP(cidr string) net.IP {
	var nextRangeFirstIP net.IP

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatal(err)
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		nextRangeFirstIP = ip
	}
	return nextRangeFirstIP
}

func subnetsInUse() []string {
	var subnets []string
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if addrs, err := inter.Addrs(); err == nil {
			for _, addr := range addrs {
				if ok, _ := regexp.MatchString("br-*", inter.Name); ok {
					subnets = append(subnets, addr.String())
				}
			}
		}
	}
	sort.Strings(subnets)
	return subnets
}

func main() {
	subnetsInUse := subnetsInUse()
	if len(subnetsInUse) > 0 {
		lastUsedNet := subnetsInUse[len(subnetsInUse)-1]
		fmt.Println(nextRangeFirstIP(lastUsedNet))
	} else {
		fmt.Println("172.55.0.1")
	}
}
