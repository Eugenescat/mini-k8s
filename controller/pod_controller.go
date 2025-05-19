package controller

import (
	"mini-k8s/api"
	"mini-k8s/runtime"
)

type PodController struct {
	runtime runtime.Runtime
	pods    map[string]api.PodSpec // 内存存储，后续可持久化
}

func NewPodController(rt runtime.Runtime) *PodController {
	return &PodController{
		runtime: rt,
		pods:    make(map[string]api.PodSpec),
	}
}

func (pc *PodController) CreatePod(spec api.PodSpec) error {
	if err := pc.runtime.CreateContainer(spec); err != nil {
		return err
	}
	pc.pods[spec.Name] = spec
	return nil
}

func (pc *PodController) DeletePod(name string) error {
	if err := pc.runtime.DeleteContainer(name); err != nil {
		return err
	}
	delete(pc.pods, name)
	return nil
}

func (pc *PodController) GetPod(name string) (*api.PodSpec, error) {
	spec, ok := pc.pods[name]
	if !ok {
		return nil, nil
	}
	return &spec, nil
}

func (pc *PodController) ListPods() ([]api.PodSpec, error) {
	pods := make([]api.PodSpec, 0, len(pc.pods))
	for _, pod := range pc.pods {
		pods = append(pods, pod)
	}
	return pods, nil
}
