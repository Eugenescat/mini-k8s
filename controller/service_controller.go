package controller

import (
	"mini-k8s/api"
)

type ServiceController struct {
	services  map[string]api.ServiceSpec
	endpoints map[string][]string // serviceName -> podName 列表
}

func NewServiceController() *ServiceController {
	return &ServiceController{
		services:  make(map[string]api.ServiceSpec),
		endpoints: make(map[string][]string),
	}
}

func (sc *ServiceController) GetEndpoints(name string) []string {
    return sc.endpoints[name]
}
func (sc *ServiceController) SetService(name string, spec api.ServiceSpec) {
    sc.services[name] = spec
}

func (sc *ServiceController) CreateService(spec api.ServiceSpec) error {
	sc.services[spec.Name] = spec
	// endpoints 由 UpdateEndpoints 维护
	return nil
}

func (sc *ServiceController) DeleteService(name string) error {
	delete(sc.services, name)
	delete(sc.endpoints, name)
	return nil
}

func (sc *ServiceController) GetService(name string) (*api.ServiceSpec, error) {
	spec, ok := sc.services[name]
	if !ok {
		return nil, nil
	}
	return &spec, nil
}

func (sc *ServiceController) ListServices() ([]api.ServiceSpec, error) {
	list := make([]api.ServiceSpec, 0, len(sc.services))
	for _, svc := range sc.services {
		list = append(list, svc)
	}
	return list, nil
}

// 根据 selector 匹配 pod，维护 endpoints
func (sc *ServiceController) UpdateEndpoints(pods []api.PodSpec) {
	for svcName, svc := range sc.services {
		matched := []string{}
		for _, pod := range pods {
			// 简单实现：selector 只支持 app: name
			if val, ok := svc.Selector["app"]; ok && pod.Name == val {
				matched = append(matched, pod.Name)
			}
		}
		sc.endpoints[svcName] = matched
	}
}
