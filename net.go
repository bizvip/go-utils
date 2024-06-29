/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package goutils

import (
	"net"
)

type NetUtils struct{}

func NewNetUtils() *NetUtils {
	return &NetUtils{}
}

func (n *NetUtils) GetLocalIP() string {
	addresses, _ := net.InterfaceAddrs()
	for _, addr := range addresses {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil && !n.IsPrivateIP(ip.IP) {
				return ip.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func (n *NetUtils) IsPrivateIP(ip net.IP) bool {
	privateIPBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, cidr := range privateIPBlocks {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}
