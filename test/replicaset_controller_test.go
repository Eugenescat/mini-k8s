package test

import (
	"mini-k8s/api"
	"mini-k8s/controller"
	"testing"
)

type FakeRuntime struct{}

func (f *FakeRuntime) CreateContainer(spec api.PodSpec) error         { return nil }
func (f *FakeRuntime) DeleteContainer(name string) error              { return nil }
func (f *FakeRuntime) GetContainerStatus(name string) (string, error) { return "Running", nil }

func TestReplicaSetController(t *testing.T) {
	// 1. 初始化 PodController 和 ReplicaSetController
	pc := controller.NewPodController(&FakeRuntime{}) // 用 FakeRuntime 防止 panic
	rc := controller.NewReplicaSetController(pc)

	// 2. 创建 ReplicaSet，期望 2 个副本
	rs := api.ReplicaSetSpec{
		Name:     "nginx-rs",
		Replicas: 2,
		Selector: map[string]string{"app": "nginx"},
		Template: api.PodSpec{
			Name:  "nginx", // 只用于 selector，实际创建时会自动加 -0/-1
			Image: "nginx:latest",
		},
	}
	rc.CreateReplicaSet(rs)

	// 3. 第一次同步副本数
	rc.SyncReplicaSets()
	pods, _ := pc.ListPods()
	if len(pods) != 2 {
		t.Errorf("expected 2 pods after first sync, got %d", len(pods))
	}
	expectedNames := map[string]bool{"nginx-rs-0": true, "nginx-rs-1": true}
	for _, pod := range pods {
		if !expectedNames[pod.Name] {
			t.Errorf("unexpected pod name: %s", pod.Name)
		}
	}

	// 4. 手动删除一个 Pod
	pc.DeletePod("nginx-rs-0")
	rc.SyncReplicaSets()
	pods, _ = pc.ListPods()
	if len(pods) != 2 {
		t.Errorf("expected 2 pods after resync, got %d", len(pods))
	}

	// 5. 修改副本数为 1，测试自动缩容
	rs.Replicas = 1
	rc.CreateReplicaSet(rs)
	rc.SyncReplicaSets()
	pods, _ = pc.ListPods()
	if len(pods) != 1 {
		t.Errorf("expected 1 pod after scaling down, got %d", len(pods))
	}
}
