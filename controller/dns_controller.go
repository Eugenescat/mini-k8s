package controller

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"mini-k8s/api"
)

// DNSController 负责管理 DNS 配置和服务发现
type DNSController struct {
	// dnsEntries 存储域名到服务的映射
	dnsEntries map[string]*api.DNSSpec
	// serviceCache 存储服务信息
	serviceCache map[string]*api.ServiceSpec
	// mutex 保护并发访问
	mutex sync.RWMutex
	// stopCh 用于停止控制器
	stopCh chan struct{}
}

// NewDNSController 创建一个新的 DNS 控制器
func NewDNSController() *DNSController {
	return &DNSController{
		dnsEntries:   make(map[string]*api.DNSSpec),
		serviceCache: make(map[string]*api.ServiceSpec),
		stopCh:       make(chan struct{}),
	}
}

// Start 启动 DNS 控制器
func (dc *DNSController) Start(ctx context.Context) error {
	go dc.run(ctx)
	return nil
}

// Stop 停止 DNS 控制器
func (dc *DNSController) Stop() {
	close(dc.stopCh)
}

// run 运行 DNS 控制器的主循环
func (dc *DNSController) run(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-dc.stopCh:
			return
		case <-ticker.C:
			dc.reconcile()
		}
	}
}

// reconcile 协调 DNS 配置和实际状态
func (dc *DNSController) reconcile() {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	// 验证所有 DNS 条目的有效性
	for host, dnsSpec := range dc.dnsEntries {
		for _, path := range dnsSpec.Paths {
			// 检查服务是否存在
			if _, exists := dc.serviceCache[path.Service]; !exists {
				fmt.Printf("Warning: Service %s for host %s path %s does not exist\n",
					path.Service, host, path.Path)
			}
		}
	}
}

// UpdateDNS 更新 DNS 配置
func (dc *DNSController) UpdateDNS(dnsSpec *api.DNSSpec) error {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	// 验证 DNS 配置
	if err := dc.validateDNSSpec(dnsSpec); err != nil {
		return err
	}

	// 更新 DNS 条目
	dc.dnsEntries[dnsSpec.Host] = dnsSpec
	return nil
}

// UpdateService 更新服务信息
func (dc *DNSController) UpdateService(serviceSpec *api.ServiceSpec) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	dc.serviceCache[serviceSpec.Name] = serviceSpec
}

// DeleteDNS 删除 DNS 配置
func (dc *DNSController) DeleteDNS(host string) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	delete(dc.dnsEntries, host)
}

// DeleteService 删除服务信息
func (dc *DNSController) DeleteService(serviceName string) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	delete(dc.serviceCache, serviceName)
}

// GetServiceForHostAndPath 根据主机名和路径查找对应的服务
func (dc *DNSController) GetServiceForHostAndPath(host, path string) (*api.ServiceSpec, int, error) {
	host = strings.ToLower(strings.TrimSuffix(host, "."))
	dc.mutex.RLock()
	defer dc.mutex.RUnlock()

	dnsSpec, exists := dc.dnsEntries[host]
	if !exists {
		return nil, 0, fmt.Errorf("host %s not found", host)
	}

	// 查找匹配的路径
	for _, dnsPath := range dnsSpec.Paths {
		if dnsPath.Path == path {
			service, exists := dc.serviceCache[dnsPath.Service]
			if !exists {
				return nil, 0, fmt.Errorf("service %s not found", dnsPath.Service)
			}
			return service, dnsPath.Port, nil
		}
	}

	return nil, 0, fmt.Errorf("no matching path found for host %s and path %s", host, path)
}

// validateDNSSpec 验证 DNS 配置的有效性
func (dc *DNSController) validateDNSSpec(dnsSpec *api.DNSSpec) error {
	if dnsSpec.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if len(dnsSpec.Paths) == 0 {
		return fmt.Errorf("at least one path must be specified")
	}

	// 验证所有路径配置
	for _, path := range dnsSpec.Paths {
		if path.Path == "" {
			return fmt.Errorf("path cannot be empty")
		}
		if path.Service == "" {
			return fmt.Errorf("service cannot be empty for path %s", path.Path)
		}
		if path.Port <= 0 {
			return fmt.Errorf("invalid port number for path %s", path.Path)
		}
	}

	return nil
}
