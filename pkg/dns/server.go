package dns

import (
	"fmt"
	"net"
	"strings"

	"mini-k8s/controller"

	"github.com/miekg/dns"
)

// Server 实现 DNS 服务器
type Server struct {
	controller *controller.DNSController
	server     *dns.Server
}

// NewServer 创建一个新的 DNS 服务器
func NewServer(controller *controller.DNSController, port int) *Server {
	s := &Server{
		controller: controller,
	}

	// 创建 DNS 服务器
	dns.HandleFunc(".", s.handleDNSRequest)
	s.server = &dns.Server{
		Addr: fmt.Sprintf(":%d", port),
		Net:  "udp",
	}

	return s
}

// Start 启动 DNS 服务器
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Stop 停止 DNS 服务器
func (s *Server) Stop() error {
	return s.server.Shutdown()
}

// handleDNSRequest 处理 DNS 查询请求
func (s *Server) handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	if len(r.Question) == 0 {
		m.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(m)
		return
	}

	question := r.Question[0]
	host := strings.TrimSuffix(question.Name, ".")

	// 处理 A 记录查询
	if question.Qtype == dns.TypeA {
		// 检查是否是集群内部域名
		if strings.HasSuffix(host, ".cluster.local") {
			fmt.Printf("[DNS DEBUG] host: %s\n", host)
			// 查找服务
			_, port, err := s.controller.GetServiceForHostAndPath(host, "/")
			if err != nil {
				m.SetRcode(r, dns.RcodeNameError)
				w.WriteMsg(m)
				return
			}

			// 返回服务的 IP 地址
			// 注意：这里需要实现服务 IP 的分配和管理
			// 这里使用一个示例 IP，实际应该根据服务端口分配
			rr := &dns.A{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    300,
				},
				A: net.ParseIP(fmt.Sprintf("10.0.0.%d", port)), // 使用端口号作为 IP 的最后一段
			}
			m.Answer = append(m.Answer, rr)
		} else {
			// 对于非集群域名，返回 NXDOMAIN
			m.SetRcode(r, dns.RcodeNameError)
		}
	} else {
		// 对于其他类型的查询，返回 NXDOMAIN
		m.SetRcode(r, dns.RcodeNameError)
	}

	w.WriteMsg(m)
}
