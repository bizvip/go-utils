/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package ip

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
)

// GetLocalPrivateIP 获取本机第一个私有环回地址
func GetLocalPrivateIP() string {
	addresses, _ := net.InterfaceAddrs()
	for _, addr := range addresses {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil && !IsPrivateIP(ip.IP) {
				return ip.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

// GetLocalPublicIP 获取本机的公网IP地址
func GetLocalPublicIP() (string, error) {
	// 调用外部服务来获取公网 IP
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 验证响应是否是有效的IP地址
	publicIP := strings.TrimSpace(string(ip))
	if net.ParseIP(publicIP) == nil {
		return "", errors.New("无法解析的公网 IP 地址")
	}

	return publicIP, nil
}

// IsPrivateIP 判断给定的 IP 地址是否是私有IP（包括环回IP和局域网IP）
func IsPrivateIP(ip net.IP) bool {
	// 检查是否是环回 IP
	if ip.IsLoopback() {
		return true
	}

	// 检查是否是私有网络 IP
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
