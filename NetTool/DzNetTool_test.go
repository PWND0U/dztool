package NetTool

import (
	"math/big"
	"net"
	"strings"
	"testing"
)

func TestIpv4ToLong(t *testing.T) {
	result := Ipv4ToLong("192.168.1.1")
	if result != 3232235777 {
		t.Errorf("Ipv4ToLong(\"192.168.1.1\") = %d, 期望 3232235777", result)
	}

	result = Ipv4ToLong("0.0.0.0")
	if result != 0 {
		t.Errorf("Ipv4ToLong(\"0.0.0.0\") = %d, 期望 0", result)
	}

	result = Ipv4ToLong("255.255.255.255")
	if result != 4294967295 {
		t.Errorf("Ipv4ToLong(\"255.255.255.255\") = %d, 期望 4294967295", result)
	}

	result = Ipv4ToLong("invalid")
	if result != 0 {
		t.Errorf("Ipv4ToLong(\"invalid\") = %d, 期望 0", result)
	}

	result = Ipv4ToLong("999.999.999.999")
	if result != 0 {
		t.Errorf("Ipv4ToLong(\"999.999.999.999\") = %d, 期望 0", result)
	}
}

func TestLongToIpv4(t *testing.T) {
	result := LongToIpv4(3232235777)
	if result != "192.168.1.1" {
		t.Errorf("LongToIpv4(3232235777) = %s, 期望 \"192.168.1.1\"", result)
	}

	result = LongToIpv4(0)
	if result != "0.0.0.0" {
		t.Errorf("LongToIpv4(0) = %s, 期望 \"0.0.0.0\"", result)
	}
}

func TestIpv4ToLongAndBack(t *testing.T) {
	testIPs := []string{"192.168.1.1", "10.0.0.1", "172.16.0.1", "127.0.0.1", "8.8.8.8", "255.255.255.255", "0.0.0.0"}
	for _, ip := range testIPs {
		n := Ipv4ToLong(ip)
		back := LongToIpv4(n)
		if back != ip {
			t.Errorf("IPv4→Long→IPv4 往返不一致: 原始=%s, 转换后=%s", ip, back)
		}
	}
}

func TestIpv6ToBigInt(t *testing.T) {
	result, err := Ipv6ToBigInt("::1")
	if err != nil {
		t.Fatalf("Ipv6ToBigInt(\"::1\") 返回错误: %v", err)
	}
	expected := big.NewInt(1)
	if result.Cmp(expected) != 0 {
		t.Errorf("Ipv6ToBigInt(\"::1\") = %v, 期望 %v", result, expected)
	}

	result, err = Ipv6ToBigInt("fe80::1")
	if err != nil {
		t.Fatalf("Ipv6ToBigInt(\"fe80::1\") 返回错误: %v", err)
	}
	if result.BitLen() == 0 {
		t.Errorf("Ipv6ToBigInt(\"fe80::1\") 结果为零值")
	}

	_, err = Ipv6ToBigInt("invalid")
	if err == nil {
		t.Errorf("Ipv6ToBigInt(\"invalid\") 期望返回错误，但返回 nil")
	}
}

func TestBigIntToIpv6(t *testing.T) {
	result := BigIntToIpv6(big.NewInt(1))
	if result != "::1" {
		t.Errorf("BigIntToIpv6(1) = %s, 期望 \"::1\"", result)
	}

	b, _ := Ipv6ToBigInt("fe80::1")
	back := BigIntToIpv6(b)
	parsedBack, err := Ipv6ToBigInt(back)
	if err != nil {
		t.Fatalf("BigIntToIpv6 转回后无法解析: %v", err)
	}
	if parsedBack.Cmp(b) != 0 {
		t.Errorf("BigInt→IPv6→BigInt 往返不一致: 原始=%v, 转换后=%v", b, parsedBack)
	}
}

func TestIsValidPort(t *testing.T) {
	if IsValidPort(-1) {
		t.Errorf("IsValidPort(-1) = true, 期望 false")
	}
	if !IsValidPort(0) {
		t.Errorf("IsValidPort(0) = false, 期望 true")
	}
	if !IsValidPort(80) {
		t.Errorf("IsValidPort(80) = false, 期望 true")
	}
	if !IsValidPort(65535) {
		t.Errorf("IsValidPort(65535) = false, 期望 true")
	}
	if IsValidPort(65536) {
		t.Errorf("IsValidPort(65536) = true, 期望 false")
	}
}

func TestIsUsableLocalPort(t *testing.T) {
	port := 58000 + Ipv4ToLong("127.0.0.1")%7000
	if !IsUsableLocalPort(int(port)) {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("无法创建监听器: %v", err)
		}
		l.Close()
		t.Errorf("IsUsableLocalPort(%d) = false, 但系统可以监听随机端口", port)
	}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("无法创建监听器: %v", err)
	}
	usedPort := l.Addr().(*net.TCPAddr).Port
	if IsUsableLocalPort(usedPort) {
		t.Errorf("IsUsableLocalPort(%d) = true, 但该端口已被占用", usedPort)
	}
	l.Close()
}

func TestGetUsableLocalPort(t *testing.T) {
	port := GetUsableLocalPort()
	if port <= 0 {
		t.Errorf("GetUsableLocalPort() = %d, 期望 > 0", port)
	}
	if !IsValidPort(port) {
		t.Errorf("GetUsableLocalPort() = %d, 不在有效端口范围内", port)
	}
}

func TestGetUsableLocalPortBetween(t *testing.T) {
	port, err := GetUsableLocalPortBetween(30000, 40000)
	if err != nil {
		t.Fatalf("GetUsableLocalPortBetween(30000, 40000) 返回错误: %v", err)
	}
	if port < 30000 || port > 40000 {
		t.Errorf("GetUsableLocalPortBetween(30000, 40000) = %d, 不在范围内", port)
	}
}

func TestGetUsableLocalPorts(t *testing.T) {
	ports, err := GetUsableLocalPorts(3, 40000, 50000)
	if err != nil {
		t.Fatalf("GetUsableLocalPorts(3, 40000, 50000) 返回错误: %v", err)
	}
	if len(ports) != 3 {
		t.Errorf("GetUsableLocalPorts 返回 %d 个端口, 期望 3", len(ports))
	}
	for _, p := range ports {
		if p < 40000 || p > 50000 {
			t.Errorf("端口 %d 不在 40000~50000 范围内", p)
		}
	}
}

func TestIsInnerIP(t *testing.T) {
	innerIPs := []string{"10.0.0.1", "172.16.0.1", "192.168.1.1", "127.0.0.1"}
	for _, ip := range innerIPs {
		if !IsInnerIP(ip) {
			t.Errorf("IsInnerIP(\"%s\") = false, 期望 true", ip)
		}
	}

	outerIPs := []string{"8.8.8.8", "1.2.3.4"}
	for _, ip := range outerIPs {
		if IsInnerIP(ip) {
			t.Errorf("IsInnerIP(\"%s\") = true, 期望 false", ip)
		}
	}
}

func TestIsInRange(t *testing.T) {
	if !IsInRange("192.168.1.1", "192.168.1.0/24") {
		t.Errorf("IsInRange(\"192.168.1.1\", \"192.168.1.0/24\") = false, 期望 true")
	}
	if IsInRange("192.168.1.1", "10.0.0.0/8") {
		t.Errorf("IsInRange(\"192.168.1.1\", \"10.0.0.0/8\") = true, 期望 false")
	}
	if !IsInRange("10.5.3.2", "10.0.0.0/8") {
		t.Errorf("IsInRange(\"10.5.3.2\", \"10.0.0.0/8\") = false, 期望 true")
	}
	if IsInRange("invalid", "10.0.0.0/8") {
		t.Errorf("IsInRange(\"invalid\", \"10.0.0.0/8\") = true, 期望 false")
	}
	if IsInRange("10.0.0.1", "invalid") {
		t.Errorf("IsInRange(\"10.0.0.1\", \"invalid\") = true, 期望 false")
	}
}

func TestHideIpPart(t *testing.T) {
	result := HideIpPart("192.168.1.1")
	if result != "192.168.1.*" {
		t.Errorf("HideIpPart(\"192.168.1.1\") = %s, 期望 \"192.168.1.*\"", result)
	}

	result = HideIpPart("10.0.0.5")
	if result != "10.0.0.*" {
		t.Errorf("HideIpPart(\"10.0.0.5\") = %s, 期望 \"10.0.0.*\"", result)
	}
}

func TestHideIpPartFromLong(t *testing.T) {
	result := HideIpPartFromLong(3232235777)
	if result != "192.168.1.*" {
		t.Errorf("HideIpPartFromLong(3232235777) = %s, 期望 \"192.168.1.*\"", result)
	}
}

func TestBuildTCPAddr(t *testing.T) {
	addr, err := BuildTCPAddr("", 8080)
	if err != nil {
		t.Fatalf("BuildTCPAddr(\"\", 8080) 返回错误: %v", err)
	}
	if addr.IP.String() != "127.0.0.1" {
		t.Errorf("空 host 时 IP = %s, 期望 127.0.0.1", addr.IP.String())
	}
	if addr.Port != 8080 {
		t.Errorf("空 host 时 Port = %d, 期望 8080", addr.Port)
	}

	addr, err = BuildTCPAddr("192.168.1.1:3306", 8080)
	if err != nil {
		t.Fatalf("BuildTCPAddr(\"192.168.1.1:3306\", 8080) 返回错误: %v", err)
	}
	if addr.IP.String() != "192.168.1.1" {
		t.Errorf("带端口的 host IP = %s, 期望 192.168.1.1", addr.IP.String())
	}
	if addr.Port != 3306 {
		t.Errorf("带端口的 host Port = %d, 期望 3306", addr.Port)
	}

	addr, err = BuildTCPAddr("192.168.1.1", 9090)
	if err != nil {
		t.Fatalf("BuildTCPAddr(\"192.168.1.1\", 9090) 返回错误: %v", err)
	}
	if addr.IP.String() != "192.168.1.1" {
		t.Errorf("纯 host IP = %s, 期望 192.168.1.1", addr.IP.String())
	}
	if addr.Port != 9090 {
		t.Errorf("纯 host Port = %d, 期望 9090", addr.Port)
	}
}

func TestCreateTCPAddr(t *testing.T) {
	addr, err := CreateTCPAddr("127.0.0.1", 8080)
	if err != nil {
		t.Fatalf("CreateTCPAddr(\"127.0.0.1\", 8080) 返回错误: %v", err)
	}
	if addr.IP.String() != "127.0.0.1" {
		t.Errorf("IP = %s, 期望 127.0.0.1", addr.IP.String())
	}
	if addr.Port != 8080 {
		t.Errorf("Port = %d, 期望 8080", addr.Port)
	}
}

func TestGetIpByHost(t *testing.T) {
	result := GetIpByHost("localhost")
	if result == "localhost" {
		t.Errorf("GetIpByHost(\"localhost\") 未解析出 IP，返回了原字符串")
	}

	result = GetIpByHost("this.domain.does.not.exist.invalid")
	if result != "this.domain.does.not.exist.invalid" {
		t.Errorf("GetIpByHost 无效域名应返回原字符串, 实际返回 %s", result)
	}
}

func TestGetNetworkInterfaces(t *testing.T) {
	ifaces, err := GetNetworkInterfaces()
	if err != nil {
		t.Fatalf("GetNetworkInterfaces() 返回错误: %v", err)
	}
	if len(ifaces) == 0 {
		t.Errorf("GetNetworkInterfaces() 返回空列表, 期望至少有一个网卡")
	}
}

func TestLocalIpv4s(t *testing.T) {
	ips, err := LocalIpv4s()
	if err != nil {
		t.Fatalf("LocalIpv4s() 返回错误: %v", err)
	}
	for _, ip := range ips {
		parsed := net.ParseIP(ip)
		if parsed == nil || parsed.To4() == nil {
			t.Errorf("LocalIpv4s() 返回非 IPv4 地址: %s", ip)
		}
	}
}

func TestLocalIps(t *testing.T) {
	ips, err := LocalIps()
	if err != nil {
		t.Fatalf("LocalIps() 返回错误: %v", err)
	}
	for _, ip := range ips {
		if net.ParseIP(ip) == nil {
			t.Errorf("LocalIps() 返回无效 IP: %s", ip)
		}
	}
}

func TestGetLocalhostStr(t *testing.T) {
	s := GetLocalhostStr()
	if s == "" {
		t.Logf("GetLocalhostStr() 返回空字符串，本机可能无可用 IP")
	}
}

func TestGetLocalhost(t *testing.T) {
	ip, err := GetLocalhost()
	if err != nil {
		t.Logf("GetLocalhost() 返回错误: %v，本机可能无可用 IP", err)
		return
	}
	if ip == nil {
		t.Errorf("GetLocalhost() 返回 nil IP")
	}
}

func TestLocalAddressList(t *testing.T) {
	addrs, err := LocalAddressList(nil)
	if err != nil {
		t.Fatalf("LocalAddressList(nil) 返回错误: %v", err)
	}
	if len(addrs) == 0 {
		t.Logf("LocalAddressList(nil) 返回空列表")
	}
}

func TestToIpList(t *testing.T) {
	addrs, err := LocalAddressList(nil)
	if err != nil {
		t.Fatalf("LocalAddressList(nil) 返回错误: %v", err)
	}
	ips := ToIpList(addrs)
	if len(ips) != len(addrs) {
		t.Errorf("ToIpList 长度 = %d, 期望 %d", len(ips), len(addrs))
	}
	for _, ip := range ips {
		if net.ParseIP(ip) == nil {
			t.Errorf("ToIpList 返回无效 IP: %s", ip)
		}
	}
}

func TestGetLocalMacAddress(t *testing.T) {
	mac, err := GetLocalMacAddress()
	if err != nil {
		t.Logf("GetLocalMacAddress() 返回错误: %v", err)
		return
	}
	if mac == "" {
		t.Errorf("GetLocalMacAddress() 返回空字符串")
	}
	if !strings.Contains(mac, "-") {
		t.Errorf("GetLocalMacAddress() 返回格式异常: %s, 期望包含 \"-\" 分隔符", mac)
	}
}

func TestGetMacAddressWithSeparator(t *testing.T) {
	ifaces, err := GetNetworkInterfaces()
	if err != nil {
		t.Fatalf("GetNetworkInterfaces() 返回错误: %v", err)
	}
	var ifaceWithMac *net.Interface
	for i := range ifaces {
		if len(ifaces[i].HardwareAddr) > 0 && ifaces[i].Flags&net.FlagLoopback == 0 {
			ifaceWithMac = &ifaces[i]
			break
		}
	}
	if ifaceWithMac == nil {
		t.Logf("未找到有 MAC 地址的非回环网卡，跳过 TestGetMacAddressWithSeparator")
		return
	}

	mac, err := GetMacAddressWithSeparator(ifaceWithMac, ":")
	if err != nil {
		t.Fatalf("GetMacAddressWithSeparator 返回错误: %v", err)
	}
	if !strings.Contains(mac, ":") {
		t.Errorf("GetMacAddressWithSeparator 使用 \":\" 分隔符, 结果: %s 不包含 \":\"", mac)
	}
	if strings.Contains(mac, "-") {
		t.Errorf("GetMacAddressWithSeparator 使用 \":\" 分隔符, 结果不应包含 \"-\": %s", mac)
	}
}

func TestGetLocalHostName(t *testing.T) {
	hostname := GetLocalHostName()
	if hostname == "" {
		t.Errorf("GetLocalHostName() 返回空字符串, 期望非空")
	}
}

func TestToAbsoluteUrl(t *testing.T) {
	result := ToAbsoluteUrl("http://example.com/base/", "page.html")
	if result != "http://example.com/base/page.html" {
		t.Errorf("ToAbsoluteUrl = %s, 期望 \"http://example.com/base/page.html\"", result)
	}

	result = ToAbsoluteUrl("http://example.com/base/", "/other")
	if result != "http://example.com/other" {
		t.Errorf("ToAbsoluteUrl 绝对路径 = %s, 期望 \"http://example.com/other\"", result)
	}

	result = ToAbsoluteUrl("http://example.com/", "http://other.com/page")
	if result != "http://other.com/page" {
		t.Errorf("ToAbsoluteUrl 完整 URL = %s, 期望 \"http://other.com/page\"", result)
	}
}

func TestIsUnknown(t *testing.T) {
	unknownValues := []string{"", "unknown", "UNKNOWN", "null", "NULL", "Unknown", "Null"}
	for _, v := range unknownValues {
		if !IsUnknown(v) {
			t.Errorf("IsUnknown(\"%s\") = false, 期望 true", v)
		}
	}

	knownValues := []string{"192.168.1.1", "hello", "0.0.0.0"}
	for _, v := range knownValues {
		if IsUnknown(v) {
			t.Errorf("IsUnknown(\"%s\") = true, 期望 false", v)
		}
	}
}

func TestGetMultistageReverseProxyIp(t *testing.T) {
	result := GetMultistageReverseProxyIp("192.168.1.1, 10.0.0.1")
	if result != "10.0.0.1" {
		t.Errorf("GetMultistageReverseProxyIp(\"192.168.1.1, 10.0.0.1\") = %s, 期望 \"10.0.0.1\"", result)
	}

	result = GetMultistageReverseProxyIp("unknown")
	if result != "" {
		t.Errorf("GetMultistageReverseProxyIp(\"unknown\") = %s, 期望 \"\"", result)
	}

	result = GetMultistageReverseProxyIp("")
	if result != "" {
		t.Errorf("GetMultistageReverseProxyIp(\"\") = %s, 期望 \"\"", result)
	}

	result = GetMultistageReverseProxyIp("10.0.0.1")
	if result != "10.0.0.1" {
		t.Errorf("GetMultistageReverseProxyIp(\"10.0.0.1\") = %s, 期望 \"10.0.0.1\"", result)
	}
}

func TestParseCookies(t *testing.T) {
	cookies := ParseCookies("name=value; name2=value2")
	if len(cookies) != 2 {
		t.Fatalf("ParseCookies 返回 %d 个 Cookie, 期望 2", len(cookies))
	}
	if cookies[0].Name != "name" || cookies[0].Value != "value" {
		t.Errorf("第一个 Cookie = %s=%s, 期望 name=value", cookies[0].Name, cookies[0].Value)
	}
	if cookies[1].Name != "name2" || cookies[1].Value != "value2" {
		t.Errorf("第二个 Cookie = %s=%s, 期望 name2=value2", cookies[1].Name, cookies[1].Value)
	}
}

func TestIdnToASCII(t *testing.T) {
	result, err := IdnToASCII("中国.com")
	t.Logf("IdnToASCII(\"中国.com\") = %s, 期望包含 Punycode 前缀 \"xn--\"", result)
	if err != nil {
		t.Fatalf("IdnToASCII(\"中国.com\") 返回错误: %v", err)
	}
	if !strings.Contains(result, "xn--") {
		t.Errorf("IdnToASCII(\"中国.com\") = %s, 期望包含 Punycode 前缀 \"xn--\"", result)
	}

	result, err = IdnToASCII("example.com")
	if err != nil {
		t.Fatalf("IdnToASCII(\"example.com\") 返回错误: %v", err)
	}
	if result != "example.com" {
		t.Errorf("IdnToASCII(\"example.com\") = %s, 期望 \"example.com\"", result)
	}
}
