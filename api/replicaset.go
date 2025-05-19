package api

// ReplicaSetSpec 定义 ReplicaSet 的期望状态
// kind, name, replicas, selector, template

type ReplicaSetSpec struct {
	Kind     string            `yaml:"kind" json:"kind"`
	Name     string            `yaml:"name" json:"name"`
	Replicas int               `yaml:"replicas" json:"replicas"`
	Selector map[string]string `yaml:"selector" json:"selector"`
	Template PodSpec           `yaml:"template" json:"template"`
}
