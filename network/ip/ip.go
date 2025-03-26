package ip

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/goccy/go-json"
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
	privateIPBlocks := []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}

	for _, cidr := range privateIPBlocks {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// GetClientIP 获取客户端的 IP 地址
func GetClientIP(r *http.Request) string {
	// 定义所有可能包含客户端 IP 的 HTTP 头，按优先级排序
	possibleHeaders := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"CF-Connecting-IP", // Cloudflare
		"True-Client-IP",   // Akamai 和其他 CDN
		"X-Client-IP",      // 某些代理服务器
		"Forwarded-For",    // RFC 7239
		"Forwarded",        // RFC 7239 标准
	}

	// 1. 检查所有可能的头，按优先级获取
	for _, header := range possibleHeaders {
		if value := r.Header.Get(header); value != "" {
			// 某些头可能包含多个 IP，以逗号分隔，取第一个非私有 IP
			if header == "X-Forwarded-For" || header == "Forwarded-For" || header == "Forwarded" {
				ips := strings.Split(value, ",")
				for _, ip := range ips {
					ip = strings.TrimSpace(ip)
					// 从 Forwarded 头中提取 for= 参数
					if header == "Forwarded" {
						if forIP := extractForIPFromForwarded(ip); forIP != "" {
							ip = forIP
						}
					}
					if IsValidPublicIP(ip) {
						return ip
					}
				}
			} else if IsValidPublicIP(value) {
				return value
			}
		}
	}

	// 2. 使用 RemoteAddr 作为备用方案
	if r.RemoteAddr != "" {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil && IsValidPublicIP(ip) {
			return ip
		}

		// 处理没有端口的情况
		if err != nil && IsValidPublicIP(r.RemoteAddr) {
			return r.RemoteAddr
		}
	}

	// 3. 最后的备用选项：尝试获取任何有效 IP，即使是私有 IP
	for _, header := range possibleHeaders {
		if value := r.Header.Get(header); value != "" {
			ips := strings.Split(value, ",")
			for _, ip := range ips {
				ip = strings.TrimSpace(ip)
				if parsedIP := net.ParseIP(ip); parsedIP != nil {
					return ip
				}
			}
		}
	}

	// 如果还是没有找到，尝试从 RemoteAddr 获取任何有效 IP
	if r.RemoteAddr != "" {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil && net.ParseIP(ip) != nil {
			return ip
		}
		if err != nil && net.ParseIP(r.RemoteAddr) != nil {
			return r.RemoteAddr
		}
	}

	// 实在获取不到时返回本地环回地址
	return "127.0.0.1"
}

// extractForIPFromForwarded 从 Forwarded 头格式中提取 for 参数的 IP
// 例如从 "for=192.0.2.60;proto=http;by=203.0.113.43" 提取 "192.0.2.60"
func extractForIPFromForwarded(forwarded string) string {
	parts := strings.Split(forwarded, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(strings.ToLower(part), "for=") {
			// 提取 for= 后面的值
			ip := strings.TrimPrefix(strings.ToLower(part), "for=")
			// 处理带引号和方括号（IPv6）的情况
			ip = strings.Trim(ip, "\"[]")
			return ip
		}
	}
	return ""
}

// GetFullClientInfo 获取客户端的完整信息，包括原始 IP 和各种代理头信息
// 对于调试和审计非常有用
func GetFullClientInfo(r *http.Request) map[string]string {
	info := make(map[string]string)

	// 记录 RemoteAddr
	info["RemoteAddr"] = r.RemoteAddr

	// 记录所有可能与 IP 相关的头信息
	relevantHeaders := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"CF-Connecting-IP",
		"True-Client-IP",
		"X-Client-IP",
		"Forwarded-For",
		"Forwarded",
		"Via",
		"X-Forwarded-Proto",
		"X-Forwarded-Host",
	}

	for _, header := range relevantHeaders {
		if value := r.Header.Get(header); value != "" {
			info[header] = value
		}
	}

	// 添加最终确定的客户端 IP
	info["ClientIP"] = GetClientIP(r)

	return info
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

// ToUniqueStr 将IP变成一个唯一字符串
func ToUniqueStr(ipStr string) string {
	// 解析IP地址为 net.IP 类型
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}

	// 获取IPv4部分（如果是IPv6则会返回 nil）
	ip = ip.To4()
	if ip == nil {
		return ""
	}

	// 使用 strings.Replace 去掉所有的 '.' 符号
	return strings.ReplaceAll(ip.String(), ".", "")
}

// GeoIPInfo represents geographic information about an IP address
type GeoIPInfo struct {
	IP           string  `json:"ip"`
	Country      string  `json:"country"`
	CountryCode  string  `json:"country_code"`
	City         string  `json:"city"`
	Region       string  `json:"region"`
	RegionCode   string  `json:"region_code"`
	Timezone     string  `json:"timezone"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Organization string  `json:"organization"`
	ASN          int     `json:"asn"`
}

// 定义错误变量
var (
	ErrFailedToFetchGeoIP = errors.New("failed to fetch geo ip information")
	ErrInvalidIPAddress   = errors.New("invalid ip address")
	ErrUnmarshalGeoIPData = errors.New("failed to unmarshal geo ip data")
)

// GetGeoIPInfo 获取指定IP的地理位置信息
func GetGeoIPInfo(ipAddr string) (*GeoIPInfo, error) {
	// 验证IP地址格式
	if net.ParseIP(ipAddr) == nil {
		return nil, ErrInvalidIPAddress
	}

	// 构建请求URL
	url := "https://get.geojs.io/v1/ip/geo/" + ipAddr + ".json"

	// 发送HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, ErrFailedToFetchGeoIP
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, ErrFailedToFetchGeoIP
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrFailedToFetchGeoIP
	}

	// 解析JSON数据
	var geoInfo GeoIPInfo
	if err := json.Unmarshal(body, &geoInfo); err != nil {
		return nil, ErrUnmarshalGeoIPData
	}

	return &geoInfo, nil
}

// GetMyGeoIPInfo 获取本机公网IP的地理位置信息
func GetMyGeoIPInfo() (*GeoIPInfo, error) {
	// 先获取本机公网IP
	publicIP, err := GetLocalPublicIP()
	if err != nil {
		return nil, err
	}

	// 获取该IP的地理信息
	return GetGeoIPInfo(publicIP)
}
