# Minik8s Lab

## Project Overview

Minik8s 是一个迷你版的容器编排工具，支持多机、多容器生命周期管理、自动扩缩容、服务发现、负载均衡等核心功能。项目基于 Fminik8s 实现，支持自定义调度、微服务、Serverless 平台集成等扩展功能，旨在帮助同学深入理解云原生与容器编排系统的原理与实现。

---

## 基本功能 Basic Features

### 1. Pod abstraction and lifecycle management✅
- 支持通过 YAML 文件配置和启动 Pod。
- 支持 Pod 的自动启动、终止、状态查询（get pod, describe pod）。
- 支持多容器互访（localhost）。
- 支持配置项包括：
  - kind: pod
  - name: pod 名称
  - image: 镜像及版本
  - command: 容器命令
  - 资源限制（如 1cpu, 128MB 内存）
  - volume: 共享卷
  - port: 容器端口暴露

### 2. Service Abstraction✅
- 支持 Service 发现与负载均衡，暴露端口，支持 selector 选择 Pod。
- 支持通过 YAML 配置 Service，内容包括：
  - kind: service
  - name: service 名称
  - selector: 匹配 pod
  - ports: 暴露端口（port/targetPort）

### 3. ReplicaSet/Deployment✅
- 支持为 Service 配置多个副本（replica），自动监控和恢复 Pod。
- 支持通过 YAML 配置 ReplicaSet/Deployment。

### 4. Auto-scaling✅
- 支持基于 Service 下 Pod 资源使用情况自动扩缩容。
- 支持 CPU、内存、带宽等多种资源类型监控。
- 支持通过 YAML 配置 HPA（HorizontalPodAutoscaler）。

### 5. DNS binding✅
- 支持通过 YAML 配置 Service 的域名和路径，实现集群内服务的 DNS 访问。

### 6. Fault tolerant🤔(Working on)
- 控制面 crash 不影响已运行 Pod。
- 控制面重启后，Service 可重新访问。

### 7. GPU application support
- 支持通过 YAML 配置 GPU 任务，集成 Slurm 平台调度。

### 8. Multi-node minik8s
- 支持多机集群，支持 Node 动态加入/移除。
- 支持调度器调度策略（如 round robin）。
- Service 抽象隐藏 Pod 具体运行位置。

---

## 进阶/自选功能 Advanced/Optional Features

### Microservice
- 支持 Service Mesh（如 Istio）流量管控。
- 支持流量控制、灰度发布、服务升级等。
- 支持自定义/开源 microservice 应用部署。

### Serverless
- 支持 Function 部署，HTTP 触发调用。
- 支持 Serverless Workflow（DAG）。
- 支持自动扩容/缩容（scale-to-0）。

---

## 考核方式
- 阶段性答辩与最终答辩，需提交中期/最终文档。
- 功能实现占 80%，工程实现占 20%。

---

## 工程/组织要求
- 禁止抄袭，组内协作需分工明确。
- 每个功能需独立分支开发，PR 合并。
- 使用 gitee private 仓库，添加助教为协作者。

---

## 参考资料
- [Kubernetes 官方文档](https://kubernetes.io/zh/docs/)
- [Slurm 任务调度平台](https://docs.hpc.sjtu.edu.cn/job/slurm.html)
- [Istio Service Mesh](https://istio.io/)
- [Knative Serverless](https://knative.dev/)
