package test

import (
	"mini-k8s/api"
	"mini-k8s/controller"
	"testing"
)

func TestServiceControllerEndpoints(t *testing.T) {
	sc := controller.NewServiceController()
	svc := api.ServiceSpec{
		Name:     "test-svc",
		Selector: map[string]string{"app": "nginx"},
	}
	err := sc.CreateService(svc)
	if err != nil {
		t.Fatalf("CreateService failed: %v", err)
	}

	pods := []api.PodSpec{
		{Name: "nginx"},
		{Name: "nginx2"},
		{Name: "other"},
	}

	sc.UpdateEndpoints(pods)
	eps := sc.GetEndpoints("test-svc")
	if len(eps) != 1 || eps[0] != "nginx" {
		t.Errorf("expected endpoints [nginx], got %v", eps)
	}

	// 修改 selector，测试不同匹配
	sc.SetService("test-svc", api.ServiceSpec{
		Name:     "test-svc",
		Selector: map[string]string{"app": "nginx2"},
	})
	sc.UpdateEndpoints(pods)
	eps = sc.GetEndpoints("test-svc")
	if len(eps) != 1 || eps[0] != "nginx2" {
		t.Errorf("expected endpoints [nginx2], got %v", eps)
	}
}
