package api

// PodSpec 定义 Pod 的期望状态
// 可根据需求扩展字段
// kind, name, image, command, resources, volume, port

type PodSpec struct {
	Kind    string   `yaml:"kind" json:"kind"`
	Name    string   `yaml:"name" json:"name"`
	Image   string   `yaml:"image" json:"image"`
	Command []string `yaml:"command" json:"command"`
	CPU     string   `yaml:"cpu" json:"cpu"`
	Memory  string   `yaml:"memory" json:"memory"`
	Volume  string   `yaml:"volume" json:"volume"`
	Ports   []int    `yaml:"ports" json:"ports"`
}
