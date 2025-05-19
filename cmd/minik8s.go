package main

import (
	"fmt"
	"mini-k8s/api"
	"mini-k8s/controller"
	"mini-k8s/runtime"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadPodSpecFromFile(filename string) (*api.PodSpec, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var spec api.PodSpec
	err = yaml.Unmarshal(data, &spec)
	return &spec, err
}

func main() {
	fmt.Println("Minik8s 启动成功！欢迎使用迷你容器编排系统。")

	rt := runtime.NewDockerRuntime()
	pc := controller.NewPodController(rt)

	// 示例：加载 pod yaml 并创建 pod
	spec, err := LoadPodSpecFromFile("config/pod-example.yaml")
	if err != nil {
		fmt.Println("加载 pod 配置失败:", err)
		return
	}
	if err := pc.CreatePod(*spec); err != nil {
		fmt.Println("创建 pod 失败:", err)
		return
	}
	fmt.Println("Pod 创建成功：", spec.Name)

	// 查询 pod 状态
	status, err := rt.GetContainerStatus(spec.Name)
	if err != nil {
		fmt.Println("查询 pod 状态失败:", err)
	} else {
		fmt.Printf("Pod %s 状态: %s\n", spec.Name, status)
	}
}
