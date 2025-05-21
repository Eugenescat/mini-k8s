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

func LoadReplicaSetSpecFromFile(filename string) (*api.ReplicaSetSpec, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var spec api.ReplicaSetSpec
	err = yaml.Unmarshal(data, &spec)
	return &spec, err
}

func LoadHPASpecFromFile(filename string) (*api.HorizontalPodAutoscalerSpec, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var spec api.HorizontalPodAutoscalerSpec
	err = yaml.Unmarshal(data, &spec)
	return &spec, err
}

func main() {
	fmt.Println("Minik8s 启动成功！欢迎使用迷你容器编排系统。")

	// 初始化运行时和控制器
	rt := runtime.NewDockerRuntime()
	pc := controller.NewPodController(rt)
	rc := controller.NewReplicaSetController(pc)
	ac := controller.NewAutoscalerController(rc)

	// 示例：加载并创建 ReplicaSet
	rsSpec, err := LoadReplicaSetSpecFromFile("config/replicaset-example.yaml")
	if err != nil {
		fmt.Println("加载 ReplicaSet 配置失败:", err)
		return
	}
	if err := rc.CreateReplicaSet(*rsSpec); err != nil {
		fmt.Println("创建 ReplicaSet 失败:", err)
		return
	}
	fmt.Println("ReplicaSet 创建成功：", rsSpec.Name)

	// 示例：加载并创建 HPA
	hpaSpec, err := LoadHPASpecFromFile("config/hpa-example.yaml")
	if err != nil {
		fmt.Println("加载 HPA 配置失败:", err)
		return
	}
	if err := ac.CreateAutoscaler(*hpaSpec); err != nil {
		fmt.Println("创建 HPA 失败:", err)
		return
	}
	fmt.Println("HPA 创建成功：", hpaSpec.Name)

	// 启动 ReplicaSet 控制器同步循环
	go rc.SyncLoop()

	// 保持主程序运行
	select {}
}
