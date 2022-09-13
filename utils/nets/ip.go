/**
 * @Author raven
 * @Description
 * @Date 2022/8/30
 **/
package nets

import (
	"errors"
	"net"
	"strconv"
)

// GetLocalIP ...
// TODO
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("unable to determine local ip")
}

// GetLocalIP ...
func GetLocalMainIP() (string, int, error) {
	// UDP Connect, no handshake
	conn, err := net.Dial("udp", "8.8.8.8:8")
	if err != nil {
		return "", 0, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), localAddr.Port, nil
}

func GetAddrAndPort(addr string) (string, int) {
	ip, portStr, _ := net.SplitHostPort(addr)
	port, _ := strconv.Atoi(portStr)
	return ip, port
}