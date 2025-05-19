// test/service_network_test.go
package test

import (
    "testing"
    "mini-k8s/network"
    "time"
)

func TestTCPProxy(t *testing.T) {
    backends := []string{"127.0.0.1:8081", "127.0.0.1:8082"}
    go network.StartTCPProxy(18080, backends)
    time.Sleep(10 * time.Second) // 等待一会儿，手动用 curl 测试
}