package api

// GPUSpec 定义 GPU 任务的期望状态
// kind, name, job, resources, slurm

type GPUSpec struct {
	Kind     string `yaml:"kind" json:"kind"`
	Name     string `yaml:"name" json:"name"`
	Job      string `yaml:"job" json:"job"`
	Resource string `yaml:"resource" json:"resource"`
	Slurm    string `yaml:"slurm" json:"slurm"`
}
