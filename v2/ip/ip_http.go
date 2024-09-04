/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package ip

import (
	"net"
	"net/http"
	"strings"
)

// GetClientIP 获取客户端的 IP 地址
func GetClientIP(r *http.Request) string {
	// 优先获取 X-Forwarded-For 头中的 IP 地址
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			// 检查 IP 是否有效且非私有
			if IsValidPublicIP(ip) {
				return ip
			}
		}
	}

	// 其次获取 X-Real-IP 头中的 IP 地址
	xri := r.Header.Get("X-Real-IP")
	if xri != "" && IsValidPublicIP(xri) {
		return xri
	}

	// 最后使用 RemoteAddr 作为备用方案
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && IsValidPublicIP(ip) {
		return ip
	}

	return ""
}

// IsValidPublicIP 验证一个 IP 地址是否是有效的公网 IP 和ip_address文件里面的检测是否私有地址是不同的
func IsValidPublicIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	// 检查是否是私有IP或环回地址
	if parsedIP.IsPrivate() || parsedIP.IsLoopback() || parsedIP.IsUnspecified() {
		return false
	}
	return true
}
