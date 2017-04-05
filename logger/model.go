package logger

import "time"

// Message fields
type Message struct {
	Log        string     `json:"log"`
	Stream     string     `json:"stream"`
	Kubernetes Kubernetes `json:"kubernetes"`
	Docker     Docker     `json:"docker"`
	Time       time.Time  `json:"time"`
}

// Kubernetes specific log message fields
type Kubernetes struct {
	Namespace     string            `json:"namespace_name"`
	PodID         string            `json:"pod_id"`
	PodName       string            `json:"pod_name"`
	ContainerName string            `json:"container_name"`
	Labels        map[string]string `json:"labels"`
	Host          string            `json:"host"`
}

// Docker specific log message fields
type Docker struct {
	ContainerID string `json:"container_id"`
}
