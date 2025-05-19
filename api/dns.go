package api

// DNSSpec 定义 DNS 的期望状态
// kind, name, host, paths

type DNSPath struct {
	Path    string `yaml:"path" json:"path"`
	Service string `yaml:"service" json:"service"`
	Port    int    `yaml:"port" json:"port"`
}

type DNSSpec struct {
	Kind  string    `yaml:"kind" json:"kind"`
	Name  string    `yaml:"name" json:"name"`
	Host  string    `yaml:"host" json:"host"`
	Paths []DNSPath `yaml:"paths" json:"paths"`
}
