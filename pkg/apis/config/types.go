package config

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkTrafficArgs holds arguments used to configure NetworkTraffic plugin.
type NetworkTrafficArgs struct {
	metav1.TypeMeta

	// Address of the Prometheus Server
	Address string `json:"prometheusAddress"`
	// NetworkInterface to be monitored, assume that nodes OS is homogeneous
	NetworkInterface string `json:"networkInterface"`
	// TimeRangeInMinutes used to aggregate the network metrics
	TimeRangeInMinutes int64 `json:"timeRangeInMinutes"`
}
