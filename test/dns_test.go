package test

import (
	"context"
	"testing"
	"time"

	"mini-k8s/api"
	"mini-k8s/controller"
	dnsserver "mini-k8s/pkg/dns"

	miekgdns "github.com/miekg/dns"
)

func TestDNSController(t *testing.T) {
	// 创建 DNS 控制器
	dc := controller.NewDNSController()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动控制器
	if err := dc.Start(ctx); err != nil {
		t.Fatalf("Failed to start DNS controller: %v", err)
	}

	// 创建示例服务
	service := &api.ServiceSpec{
		Kind:     "service",
		Name:     "test-service",
		Selector: map[string]string{"app": "test"},
		Ports: []api.ServicePort{
			{
				Port:       80,
				TargetPort: 8080,
			},
		},
	}

	// 更新服务
	dc.UpdateService(service)

	// 创建 DNS 配置
	dnsSpec := &api.DNSSpec{
		Kind: "dns",
		Name: "test-dns",
		Host: "test.cluster.local",
		Paths: []api.DNSPath{
			{
				Path:    "/",
				Service: "test-service",
				Port:    80,
			},
		},
	}

	// 更新 DNS 配置
	if err := dc.UpdateDNS(dnsSpec); err != nil {
		t.Fatalf("Failed to update DNS: %v", err)
	}

	// 创建 DNS 服务器，使用端口 5354
	server := dnsserver.NewServer(dc, 5354)
	go func() {
		if err := server.Start(); err != nil {
			t.Errorf("Failed to start DNS server: %v", err)
		}
	}()
	defer server.Stop()

	// 等待服务器启动
	time.Sleep(time.Second)

	// 创建 DNS 客户端
	c := new(miekgdns.Client)
	m := new(miekgdns.Msg)
	m.SetQuestion("test.cluster.local.", miekgdns.TypeA)
	m.RecursionDesired = true

	// 发送 DNS 查询到新端口
	r, _, err := c.Exchange(m, "127.0.0.1:5354")
	if err != nil {
		t.Fatalf("Failed to send DNS query: %v", err)
	}

	// 验证响应
	if r.Rcode != miekgdns.RcodeSuccess {
		t.Errorf("Expected RcodeSuccess, got %v", r.Rcode)
	}

	if len(r.Answer) == 0 {
		t.Error("Expected at least one answer, got none")
		return
	}

	// 验证 IP 地址
	if a, ok := r.Answer[0].(*miekgdns.A); ok {
		expectedIP := "10.0.0.80" // 基于端口号 80
		if a.A.String() != expectedIP {
			t.Errorf("Expected IP %s, got %s", expectedIP, a.A.String())
		}
	} else {
		t.Error("Expected A record, got different type")
	}
}
