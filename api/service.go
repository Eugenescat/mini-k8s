package api

// ServiceSpec 定义 Service 的期望状态
// kind, name, selector, ports

type ServicePort struct {
	Port       int `yaml:"port" json:"port"`
	TargetPort int `yaml:"targetPort" json:"targetPort"`
}

type ServiceSpec struct {
	Kind     string            `yaml:"kind" json:"kind"`
	Name     string            `yaml:"name" json:"name"`
	Selector map[string]string `yaml:"selector" json:"selector"`
	Ports    []ServicePort     `yaml:"ports" json:"ports"`
}
