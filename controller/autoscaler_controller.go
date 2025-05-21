package controller

import (
	"fmt"
	"mini-k8s/api"
	"mini-k8s/pkg"
	"time"
)

type AutoscalerController struct {
	autoscalers  map[string]api.HorizontalPodAutoscalerSpec
	rsController *ReplicaSetController
	metrics      map[string]float64 // podName -> current metric value
}

func NewAutoscalerController(rc *ReplicaSetController) *AutoscalerController {
	ac := &AutoscalerController{
		autoscalers:  make(map[string]api.HorizontalPodAutoscalerSpec),
		rsController: rc,
		metrics:      make(map[string]float64),
	}
	// 启动后台监控
	go ac.monitorLoop()
	return ac
}

func (ac *AutoscalerController) CreateAutoscaler(spec api.HorizontalPodAutoscalerSpec) error {
	ac.autoscalers[spec.Name] = spec
	return nil
}

func (ac *AutoscalerController) DeleteAutoscaler(name string) error {
	delete(ac.autoscalers, name)
	return nil
}

func (ac *AutoscalerController) GetAutoscaler(name string) (*api.HorizontalPodAutoscalerSpec, error) {
	spec, ok := ac.autoscalers[name]
	if !ok {
		return nil, nil
	}
	return &spec, nil
}

// monitorLoop 定期收集指标并执行扩缩容
func (ac *AutoscalerController) monitorLoop() {
	ticker := time.NewTicker(15 * time.Second)
	for range ticker.C {
		ac.collectMetrics()
		ac.reconcile()
	}
}

// collectMetrics 收集所有 Pod 的资源使用指标
func (ac *AutoscalerController) collectMetrics() {
	// TODO: 实际环境中应该从监控系统获取指标
	// 这里简单模拟：随机生成 0-100 的指标值
	for _, hpa := range ac.autoscalers {
		pods, _ := ac.rsController.podController.ListPods()
		for _, pod := range pods {
			// 只收集目标 ReplicaSet 的 Pod 指标
			if len(pod.Name) >= len(hpa.Target) && pod.Name[:len(hpa.Target)] == hpa.Target {
				// 简单实现：只监控 CPU 使用率
				ac.metrics[pod.Name] = pkg.RandomFloat(0, 100)
			}
		}
	}
}

// reconcile 根据指标执行扩缩容
func (ac *AutoscalerController) reconcile() {
	for name, hpa := range ac.autoscalers {
		// 1. 获取目标 ReplicaSet
		rs, err := ac.rsController.GetReplicaSet(hpa.Target)
		if err != nil || rs == nil {
			fmt.Printf("[HPA] ReplicaSet %s not found\n", hpa.Target)
			continue
		}

		// 2. 计算当前指标平均值
		var totalMetric float64
		var podCount int
		for podName, metric := range ac.metrics {
			if len(podName) >= len(hpa.Target) && podName[:len(hpa.Target)] == hpa.Target {
				totalMetric += metric
				podCount++
			}
		}
		if podCount == 0 {
			continue
		}
		avgMetric := totalMetric / float64(podCount)

		// 3. 根据指标计算期望副本数
		// 简单实现：只考虑 CPU 指标，目标值 50%
		var targetReplicas int
		for _, metric := range hpa.Metrics {
			if metric.Type == "cpu" {
				// 当前使用率超过目标值，需要扩容
				if avgMetric > metric.Target {
					targetReplicas = int(float64(rs.Replicas) * (avgMetric / metric.Target))
				} else {
					// 当前使用率低于目标值，可以缩容
					targetReplicas = int(float64(rs.Replicas) * (avgMetric / metric.Target))
				}
				break
			}
		}

		// 4. 确保副本数在允许范围内
		if targetReplicas < hpa.MinReplicas {
			targetReplicas = hpa.MinReplicas
		}
		if targetReplicas > hpa.MaxReplicas {
			targetReplicas = hpa.MaxReplicas
		}

		// 5. 如果副本数需要调整，更新 ReplicaSet
		if targetReplicas != rs.Replicas {
			fmt.Printf("[HPA] %s: scaling %s from %d to %d (avg metric: %.2f)\n",
				name, hpa.Target, rs.Replicas, targetReplicas, avgMetric)
			rs.Replicas = targetReplicas
			ac.rsController.CreateReplicaSet(*rs)
		}
	}
}
