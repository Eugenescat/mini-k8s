package runtime

import (
	"fmt"
	"mini-k8s/api"
	"os/exec"
)

type Runtime interface {
	CreateContainer(spec api.PodSpec) error
	DeleteContainer(name string) error
	GetContainerStatus(name string) (string, error)
}

type DockerRuntime struct{}

func NewDockerRuntime() *DockerRuntime {
	return &DockerRuntime{}
}

func (d *DockerRuntime) CreateContainer(spec api.PodSpec) error {
	args := []string{"run", "-d", "--name", spec.Name}
	if spec.CPU != "" {
		args = append(args, "--cpus", spec.CPU)
	}
	if spec.Memory != "" {
		args = append(args, "-m", spec.Memory)
	}
	for _, port := range spec.Ports {
		args = append(args, "-p", fmt.Sprintf("%d:%d", port, port))
	}
	if spec.Volume != "" {
		args = append(args, "-v", spec.Volume+":"+spec.Volume)
	}
	args = append(args, spec.Image)
	args = append(args, spec.Command...)
	cmd := exec.Command("docker", args...)
	return cmd.Run()
}

func (d *DockerRuntime) DeleteContainer(name string) error {
	cmd := exec.Command("docker", "rm", "-f", name)
	return cmd.Run()
}

func (d *DockerRuntime) GetContainerStatus(name string) (string, error) {
	cmd := exec.Command("docker", "inspect", "-f", "{{.State.Status}}", name)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
