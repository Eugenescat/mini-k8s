package api

// HorizontalPodAutoscalerSpec 定义 HPA 的期望状态
// kind, name, target, minReplicas, maxReplicas, metrics

type HPAMetric struct {
	Type   string  `yaml:"type" json:"type"` // cpu, memory, etc.
	Target float64 `yaml:"target" json:"target"`
}

type HorizontalPodAutoscalerSpec struct {
	Kind        string      `yaml:"kind" json:"kind"`
	Name        string      `yaml:"name" json:"name"`
	Target      string      `yaml:"target" json:"target"`
	MinReplicas int         `yaml:"minReplicas" json:"minReplicas"`
	MaxReplicas int         `yaml:"maxReplicas" json:"maxReplicas"`
	Metrics     []HPAMetric `yaml:"metrics" json:"metrics"`
}
