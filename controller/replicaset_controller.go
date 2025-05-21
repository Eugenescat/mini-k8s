package controller

import (
	"fmt"
	"mini-k8s/api"
	"time"
)

type ReplicaSetController struct {
	replicasets   map[string]api.ReplicaSetSpec
	podController *PodController
}

func NewReplicaSetController(pc *PodController) *ReplicaSetController {
	return &ReplicaSetController{
		replicasets:   make(map[string]api.ReplicaSetSpec),
		podController: pc,
	}
}

func (rc *ReplicaSetController) CreateReplicaSet(spec api.ReplicaSetSpec) error {
	rc.replicasets[spec.Name] = spec
	return nil
}

func (rc *ReplicaSetController) DeleteReplicaSet(name string) error {
	delete(rc.replicasets, name)
	return nil
}

// SyncReplicaSets 保证所有 ReplicaSet 的实际 Pod 数量与期望一致
func (rc *ReplicaSetController) SyncReplicaSets() {
	for _, rs := range rc.replicasets {
		// 统计当前属于该 ReplicaSet 的 Pod（通过名字前缀判断）
		pods, _ := rc.podController.ListPods()
		matched := []api.PodSpec{}
		prefix := rs.Name + "-"
		for _, pod := range pods {
			if len(pod.Name) >= len(prefix) && pod.Name[:len(prefix)] == prefix {
				matched = append(matched, pod)
			}
		}
		// 构建现有 Pod 名字集合
		existing := map[string]bool{}
		for _, pod := range matched {
			existing[pod.Name] = true
		}
		// 补全缺失编号的 Pod
		for i := 0; i < rs.Replicas; i++ {
			name := fmt.Sprintf("%s-%d", rs.Name, i)
			if !existing[name] {
				newPod := rs.Template
				newPod.Name = name
				rc.podController.CreatePod(newPod)
			}
		}
		// 删除多余的 Pod
		if len(matched) > rs.Replicas {
			// 多余的编号全部删除
			for _, pod := range matched {
				var idx int
				_, err := fmt.Sscanf(pod.Name, prefix+"%d", &idx)
				if err != nil || idx >= rs.Replicas {
					rc.podController.DeletePod(pod.Name)
				}
			}
		}
	}
}

func (rc *ReplicaSetController) GetReplicaSet(name string) (*api.ReplicaSetSpec, error) {
	spec, ok := rc.replicasets[name]
	if !ok {
		return nil, nil
	}
	return &spec, nil
}

// SyncLoop 启动 ReplicaSet 控制器的同步循环
func (rc *ReplicaSetController) SyncLoop() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		rc.SyncReplicaSets()
	}
}
