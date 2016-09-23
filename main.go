package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"sort"

	flag "github.com/docker/docker/pkg/mflag"
)

var (
	subnet, subnetmask  string
	listExistingSubnets bool
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
	flag.StringVar(&subnet, []string{"s", "-subnet"}, "172.55.0.0", "Subnet. Default: 172.55.0.0")
	flag.StringVar(&subnetmask, []string{"sm", "-subnetmask"}, "/29", "Subnetmask. Default: /29")
	flag.BoolVar(&listExistingSubnets, []string{"ls", "-listExistingSubnets"}, false, "List Existing Subnets.")
	flag.Parse()

	subnetsInUse := subnetsInUse()

	if listExistingSubnets {
		for _, subnet := range subnetsInUse {
			fmt.Println(subnet)
		}
		os.Exit(0)
	}

	if len(subnetsInUse) > 0 {
		lastUsedNet := subnetsInUse[len(subnetsInUse)-1]
		fmt.Println(nextRangeFirstIP(lastUsedNet).String() + subnetmask)
	} else {
		fmt.Println(subnet + subnetmask)
	}
}
