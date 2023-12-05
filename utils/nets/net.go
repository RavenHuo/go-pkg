package nets

import (
	"fmt"
	"net"
	"strings"
)

// IPFormat 对 IPv4,IPv6 格式化
func IPFormat(ip string) string {
	parseIP := net.ParseIP(ip)
	if parseIP == nil {
		return ip
	}
	if !strings.Contains(ip, ":") {
		return ip
	}
	return iPv6Format(ip)
}

// HostPortFormat 对 host+port 格式化
// host 支持: IPv4, IPv6, domain, hostname
func HostPortFormat(hostport string) string {
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return IPFormat(hostport)
	}

	if !strings.Contains(host, ":") {
		return hostport
	}

	var ip = host
	if hostport[0] == '[' {
		ip = fmt.Sprintf("[%s]", iPv6Format(ip))
	}

	return fmt.Sprintf("%s:%s", ip, port)
}

func iPv6Format(data string) string {
	ip := net.ParseIP(data)
	if ip == nil {
		return ""
	}
	array := []string{"0", "0", "0", "0", "0", "0", "0", "0"}
	hasIPv4 := false
	var ipv4 string
	// 包含IPv4地址
	if strings.Count(data, ".") == 3 {
		array = []string{"0", "0", "0", "0", "0", "0"}
		hasIPv4 = true
	}
	// IPv6地址里面只能包含一个"::",此处分左右两部分处理
	if strings.Contains(data, "::") {
		tmp := strings.Split(data, "::")
		left := strings.Split(tmp[0], ":")
		if len(left) > 0 {
			for i, v := range left {
				v = strings.TrimLeft(v, "0")
				if len(v) == 0 {
					v = "0"
				}
				array[i] = v
			}
		}
		// 右部分会包含IPv4地址,需要额外处理
		right := strings.Split(tmp[1], ":")
		if len(right) == 0 && hasIPv4 {
			ipv4 = tmp[1]
		}
		if len(right) > 0 {
			size := len(right)
			if hasIPv4 {
				size--
				ipv4 = right[size]
			}
			tail := len(array) - 1
			for i := size - 1; i >= 0; i-- {
				v := strings.TrimLeft(right[i], "0")
				if len(v) == 0 {
					v = "0"
				}
				array[tail] = v
				tail--
			}
		}
	} else { // 标准的IPv6格式,未缩写
		tmp := strings.Split(data, ":")
		if hasIPv4 {
			ipv4 = tmp[len(tmp)-1]
			tmp = tmp[:len(tmp)-1]
		}
		for i, v := range tmp {
			v = strings.TrimLeft(v, "0")
			if len(v) == 0 {
				v = "0"
			}
			array[i] = v
		}
	}
	if hasIPv4 {
		array = append(array, ipv4)
	}
	return strings.Join(array[:], ":")
}
