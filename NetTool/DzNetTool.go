package NetTool

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/idna"
)

const (
	LocalIP      = "127.0.0.1"
	PortRangeMin = 1024
	PortRangeMax = 65535
)

// Ipv4ToLong 将 IPv4 地址转换为 uint32
func Ipv4ToLong(ip string) uint32 {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return 0
	}
	netIP = netIP.To4()
	if netIP == nil {
		return 0
	}
	return uint32(netIP[0])<<24 | uint32(netIP[1])<<16 | uint32(netIP[2])<<8 | uint32(netIP[3])
}

// LongToIpv4 将 uint32 转换为 IPv4 地址
func LongToIpv4(n uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

// Ipv6ToBigInt 将 IPv6 地址转换为 big.Int
func Ipv6ToBigInt(ip string) (*big.Int, error) {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return nil, errors.New("无效的 IPv6 地址")
	}
	netIP = netIP.To16()
	if netIP == nil {
		return nil, errors.New("无效的 IPv6 地址")
	}
	return new(big.Int).SetBytes(netIP), nil
}

// BigIntToIpv6 将 big.Int 转换为 IPv6 地址
func BigIntToIpv6(n *big.Int) string {
	b := n.Bytes()
	if len(b) < 16 {
		padded := make([]byte, 16)
		copy(padded[16-len(b):], b)
		b = padded
	}
	return net.IP(b).String()
}

// IsValidPort 判断端口号是否有效（0~65535）
func IsValidPort(port int) bool {
	return port >= 0 && port <= 65535
}

// IsUsableLocalPort 检测本地端口是否可用
func IsUsableLocalPort(port int) bool {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	l.Close()
	return true
}

// GetUsableLocalPort 在 1024~65535 范围内随机查找可用端口
func GetUsableLocalPort() int {
	return GetUsableLocalPortRange(PortRangeMin)
}

// GetUsableLocalPortRange 在 minPort~65535 范围内随机查找可用端口
func GetUsableLocalPortRange(minPort int) int {
	maxAttempts := PortRangeMax - minPort
	for i := 0; i < maxAttempts; i++ {
		port := minPort + rand.Intn(PortRangeMax-minPort+1)
		if IsUsableLocalPort(port) {
			return port
		}
	}
	return -1
}

// GetUsableLocalPortBetween 在指定范围内查找可用端口
func GetUsableLocalPortBetween(minPort, maxPort int) (int, error) {
	if minPort < 0 {
		minPort = 0
	}
	if maxPort > PortRangeMax {
		maxPort = PortRangeMax
	}
	maxAttempts := maxPort - minPort
	for i := 0; i < maxAttempts; i++ {
		port := minPort + rand.Intn(maxPort-minPort+1)
		if IsUsableLocalPort(port) {
			return port, nil
		}
	}
	return -1, errors.New("在指定范围内未找到可用端口")
}

// GetUsableLocalPorts 查找多个可用端口
func GetUsableLocalPorts(numRequested, minPort, maxPort int) ([]int, error) {
	var ports []int
	for i := 0; i < numRequested; i++ {
		port, err := GetUsableLocalPortBetween(minPort, maxPort)
		if err != nil {
			return ports, err
		}
		ports = append(ports, port)
	}
	return ports, nil
}

// IsInnerIP 判断是否为内网 IP
func IsInnerIP(ip string) bool {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return false
	}
	// 回环地址 127.0.0.0/8
	if netIP.IsLoopback() {
		return true
	}
	// A 类 10.0.0.0/8
	_, ipNet, _ := net.ParseCIDR("10.0.0.0/8")
	if ipNet.Contains(netIP) {
		return true
	}
	// B 类 172.16.0.0/12
	_, ipNet, _ = net.ParseCIDR("172.16.0.0/12")
	if ipNet.Contains(netIP) {
		return true
	}
	// C 类 192.168.0.0/16
	_, ipNet, _ = net.ParseCIDR("192.168.0.0/16")
	if ipNet.Contains(netIP) {
		return true
	}
	return false
}

// IsInRange 判断 IP 是否在 CIDR 范围内
func IsInRange(ip string, cidr string) bool {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return false
	}
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	return ipNet.Contains(netIP)
}

// HideIpPart 隐藏 IP 最后一部分为 *
func HideIpPart(ip string) string {
	if strings.Contains(ip, ":") {
		// IPv6 隐藏最后一段
		parts := strings.Split(ip, ":")
		if len(parts) > 1 {
			parts[len(parts)-1] = "*"
			return strings.Join(parts, ":")
		}
		return ip
	}
	// IPv4 隐藏最后一部分
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		parts[3] = "*"
		return strings.Join(parts, ".")
	}
	return ip
}

// HideIpPartFromLong 先将 uint32 转为 IPv4 再隐藏最后一部分
func HideIpPartFromLong(n uint32) string {
	return HideIpPart(LongToIpv4(n))
}

// BuildTCPAddr 构建 TCPAddr，host 为空时使用 127.0.0.1
func BuildTCPAddr(host string, defaultPort int) (*net.TCPAddr, error) {
	if host == "" {
		host = LocalIP
	}
	var addr string
	if strings.Contains(host, ":") {
		h, port, err := net.SplitHostPort(host)
		if err == nil && port != "" {
			addr = fmt.Sprintf("%s:%s", h, port)
		} else {
			addr = fmt.Sprintf("[%s]:%d", host, defaultPort)
		}
	} else {
		addr = fmt.Sprintf("%s:%d", host, defaultPort)
	}
	return net.ResolveTCPAddr("tcp", addr)
}

// CreateTCPAddr 创建 TCPAddr
func CreateTCPAddr(host string, port int) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
}

// GetIpByHost 域名解析，失败返回原域名
func GetIpByHost(hostName string) string {
	addrs, err := net.LookupHost(hostName)
	if err != nil || len(addrs) == 0 {
		return hostName
	}
	return addrs[0]
}

// GetNetworkInterfaces 获取所有网卡
func GetNetworkInterfaces() ([]net.Interface, error) {
	return net.Interfaces()
}

// GetNetworkInterface 获取指定名称的网卡
func GetNetworkInterface(name string) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Name == name {
			return &iface, nil
		}
	}
	return nil, fmt.Errorf("未找到网卡: %s", name)
}

// LocalIpv4s 获取本机 IPv4 列表
func LocalIpv4s() ([]string, error) {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
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
			if ip != nil {
				ips = append(ips, ip.String())
			}
		}
	}
	return ips, nil
}

// LocalIpv6s 获取本机 IPv6 列表
func LocalIpv6s() ([]string, error) {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
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
			if ip.To4() != nil {
				continue
			}
			if ip.To16() != nil {
				ips = append(ips, ip.String())
			}
		}
	}
	return ips, nil
}

// LocalIps 获取本机所有 IP 列表
func LocalIps() ([]string, error) {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
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
			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}

// GetLocalhostStr 获取第一个非回环 IP 字符串
func GetLocalhostStr() string {
	ip, err := GetLocalhost()
	if err != nil {
		return ""
	}
	return ip.String()
}

// GetLocalhost 获取第一个非回环 IP，优先非 siteLocal 的 IPv4
func GetLocalhost() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var siteLocalIP net.IP
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
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
			ip4 := ip.To4()
			if ip4 == nil {
				continue
			}
			if !ip4.IsPrivate() {
				return ip4, nil
			}
			if siteLocalIP == nil {
				siteLocalIP = ip4
			}
		}
	}
	if siteLocalIP != nil {
		return siteLocalIP, nil
	}
	// 使用 UDP 连接获取本机出口 IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// LocalAddressList 过滤本地地址
func LocalAddressList(filter func(net.Addr) bool) ([]net.Addr, error) {
	var result []net.Addr
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if filter == nil || filter(addr) {
				result = append(result, addr)
			}
		}
	}
	return result, nil
}

// ToIpList 将地址列表转为 IP 字符串列表
func ToIpList(addrs []net.Addr) []string {
	var ips []string
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip != nil {
			ips = append(ips, ip.String())
		}
	}
	return ips
}

// GetLocalMacAddress 获取本机第一个非回环网卡的 MAC 地址
func GetLocalMacAddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		if len(iface.HardwareAddr) > 0 {
			return GetMacAddress(&iface)
		}
	}
	return "", errors.New("未找到非回环网卡")
}

// GetMacAddress 获取指定网卡的 MAC 地址，使用 "-" 分隔
func GetMacAddress(iface *net.Interface) (string, error) {
	return GetMacAddressWithSeparator(iface, "-")
}

// GetMacAddressWithSeparator 获取指定网卡的 MAC 地址，自定义分隔符
func GetMacAddressWithSeparator(iface *net.Interface, separator string) (string, error) {
	if len(iface.HardwareAddr) == 0 {
		return "", errors.New("网卡无 MAC 地址")
	}
	var parts []string
	for _, b := range iface.HardwareAddr {
		parts = append(parts, fmt.Sprintf("%02X", b))
	}
	return strings.Join(parts, separator), nil
}

// GetLocalHostName 获取本机主机名
func GetLocalHostName() string {
	hostname, _ := os.Hostname()
	return hostname
}

// NetCat 通过 TCP 发送数据
func NetCat(host string, port int, data []byte) error {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(data)
	return err
}

// ToAbsoluteUrl 将相对 URL 转为绝对 URL
func ToAbsoluteUrl(basePath string, relativePath string) string {
	base, err := url.Parse(basePath)
	if err != nil {
		return relativePath
	}
	relative, err := url.Parse(relativePath)
	if err != nil {
		return relativePath
	}
	return base.ResolveReference(relative).String()
}

// IdnToASCII 将 Unicode 域名转为 Punycode
func IdnToASCII(unicode string) (string, error) {
	return idna.Lookup.ToASCII(unicode)
}

// GetMultistageReverseProxyIp 从逗号分隔的 IP 列表中取第一个非 unknown 的 IP
func GetMultistageReverseProxyIp(ip string) string {
	if IsUnknown(ip) {
		return ""
	}
	parts := strings.Split(ip, ",")
	for i := len(parts) - 1; i >= 0; i-- {
		p := strings.TrimSpace(parts[i])
		if !IsUnknown(p) {
			return p
		}
	}
	return ""
}

// IsUnknown 检测是否为未知字符串
func IsUnknown(checkString string) bool {
	if checkString == "" {
		return true
	}
	lower := strings.ToLower(checkString)
	return lower == "unknown" || lower == "null"
}

// Ping 检测主机是否可达，默认 3000ms 超时
func Ping(ip string) bool {
	return PingWithTimeout(ip, 3000)
}

// PingWithTimeout 检测主机是否可达，指定超时（毫秒）
func PingWithTimeout(ip string, timeoutMs int) bool {
	timeoutStr := fmt.Sprintf("%d", timeoutMs)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", timeoutStr, ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", timeoutStr, ip)
	}
	err := cmd.Run()
	return err == nil
}

// ParseCookies 解析 Cookie 字符串
func ParseCookies(cookieStr string) []*http.Cookie {
	header := http.Header{}
	header.Add("Cookie", cookieStr)
	request := http.Request{Header: header}
	return request.Cookies()
}

// IsOpen 检测远程端口是否开放
func IsOpen(host string, port int, timeout time.Duration) bool {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
