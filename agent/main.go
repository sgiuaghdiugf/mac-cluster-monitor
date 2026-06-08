package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"mac-cluster-monitor/agent/config"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// HeartbeatData 心跳数据
type HeartbeatData struct {
	MacAddress  string  `json:"mac_address"`
	Hostname    string  `json:"hostname"`
	IPAddress   string  `json:"ip_address"`
	Region      string  `json:"region"`
	MemoryUsed  int64   `json:"memory_used"`
	MemoryTotal int64   `json:"memory_total"`
	CpuUsage    float64 `json:"cpu_usage"`
	CpuCores    int     `json:"cpu_cores"`
}

func main() {
	// 硬编码配置（编译时固定）
	cfg := &config.Config{
		Server: config.ServerConfig{
			URL:               "http://192.168.109.106:8080", // 监控服务器地址
			HeartbeatInterval: 5,                             // 上报间隔（秒）
			Timeout:           10,                            // 请求超时时间（秒）
		},
		Device: config.DeviceConfig{
			Name:      "",   // 设备名称（留空自动获取主机名）
			ID:        "",   // 设备ID（留空自动使用MAC地址）
			Interface: "",   // 网卡接口（留空自动选择）
			Region:    "上海", // 设备所属地区（根据实际情况修改）
		},
		Log: config.LogConfig{
			Level: "info",                           // 日志级别
			File:  "/var/log/mac-monitor-agent.log", // 日志文件路径
		},
	}

	log.Println("Mac Cluster Monitor Agent 启动...")
	log.Printf("服务器地址: %s", cfg.Server.URL)
	log.Printf("上报间隔: %d秒", cfg.Server.HeartbeatInterval)

	// 获取MAC地址和主机名
	macAddress, err := getMacAddress(cfg.Device.Interface)
	if err != nil {
		log.Fatal("获取MAC地址失败:", err)
	}

	hostname := cfg.Device.Name
	if hostname == "" {
		// 优先从系统获取友好名称
		if h, err := getFriendlyHostname(); err == nil && h != "" && !isIPAddress(h) {
			hostname = h
		} else if h, err = os.Hostname(); err == nil {
			hostname = h
		} else {
			hostname = "Mac-Mini-" + macAddress[len(macAddress)-6:]
		}
	}

	log.Printf("设备MAC: %s", macAddress)
	log.Printf("设备名称: %s", hostname)

	// 定时上报心跳
	ticker := time.NewTicker(time.Duration(cfg.Server.HeartbeatInterval) * time.Second)
	defer ticker.Stop()

	// 立即上报一次
	sendHeartbeat(cfg, macAddress, hostname)

	for range ticker.C {
		sendHeartbeat(cfg, macAddress, hostname)
	}
}

// sendHeartbeat 发送心跳
func sendHeartbeat(cfg *config.Config, macAddress, hostname string) {
	// 获取系统信息
	ipAddress := getIPAddress(cfg.Device.Interface)

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("获取内存信息失败: %v", err)
		return
	}

	// 获取CPU使用率（macOS可能不支持，使用默认值）
	var cpuUsage float64
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Printf("获取CPU使用率失败: %v，使用默认值0", err)
		cpuUsage = 0
	} else if len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	// 获取CPU信息
	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Printf("获取CPU信息失败: %v", err)
	}

	cpuCores := runtime.NumCPU()
	if len(cpuInfo) > 0 && cpuInfo[0].Cores > 0 {
		cpuCores = int(cpuInfo[0].Cores)
	}

	data := HeartbeatData{
		MacAddress:  macAddress,
		Hostname:    hostname,
		IPAddress:   ipAddress,
		Region:      cfg.Device.Region,
		MemoryUsed:  int64(memInfo.Used / 1024 / 1024),  // MB
		MemoryTotal: int64(memInfo.Total / 1024 / 1024), // MB
		CpuUsage:    cpuUsage,
		CpuCores:    cpuCores,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("序列化数据失败: %v", err)
		return
	}

	// 发送HTTP请求
	url := cfg.Server.URL + "/api/heartbeat"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("发送心跳失败: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("心跳发送成功 [IP: %s, CPU: %.1f%%, 内存: %d/%d MB]",
			ipAddress, data.CpuUsage, data.MemoryUsed, data.MemoryTotal)
	} else {
		log.Printf("心跳发送失败,状态码: %d", resp.StatusCode)
	}
}

// getMacAddress 获取MAC地址
func getMacAddress(preferredInterface string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// 如果指定了优先接口，先查找该接口
	if preferredInterface != "" {
		for _, iface := range interfaces {
			if iface.Name == preferredInterface && len(iface.HardwareAddr) > 0 {
				return iface.HardwareAddr.String(), nil
			}
		}
	}

	// 在macOS上，优先获取en0接口的MAC地址
	for _, iface := range interfaces {
		if iface.Name == "en0" && len(iface.HardwareAddr) > 0 {
			return iface.HardwareAddr.String(), nil
		}
	}

	// 如果没有en0，找第一个非lo的网卡
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback == 0 && len(iface.HardwareAddr) > 0 {
			return iface.HardwareAddr.String(), nil
		}
	}

	return "", nil
}

// getIPAddress 获取IP地址
func getIPAddress(preferredInterface string) string {
	// 如果指定了优先接口，先获取该接口的IP
	if preferredInterface != "" {
		iface, err := net.InterfaceByName(preferredInterface)
		if err == nil {
			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	// 优先获取en0的IP
	iface, err := net.InterfaceByName("en0")
	if err == nil {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}

	// 获取默认IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// getFriendlyHostname 获取友好的主机名
func getFriendlyHostname() (string, error) {
	// macOS: 尝试从系统配置获取计算机名（更友好）
	if runtime.GOOS == "darwin" {
		// 读取 ComputerName (scutil --get ComputerName)
		name := runCommand("scutil", "--get", "ComputerName")
		if name != "" && !isIPAddress(name) {
			return name, nil
		}

		// 尝试 LocalHostName
		name = runCommand("scutil", "--get", "LocalHostName")
		if name != "" && !isIPAddress(name) {
			return name, nil
		}
	}

	// 回退到 os.Hostname()
	return os.Hostname()
}

// isIPAddress 检查字符串是否是IP地址
func isIPAddress(s string) bool {
	return net.ParseIP(s) != nil
}

// runCommand 执行命令并返回输出
func runCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out.String())
}
