package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Service struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ServiceSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `son:"metadata,omitempty"`
	Items           []Service `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AppConfig struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `son:"metadata,omitempty"`
	Spec            ConfigSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AppConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `son:"metadata,omitempty"`
	Items           []AppConfig `json:"items"`
}

type ConfigSpec struct {
	Groups []GroupItem `json:"groups,omitempty"`
}

type GroupItem struct {
	Oidc       string `json:"oidc"`
	Kubernetes string `json:"kubernetes"`
}

type ServiceSpec struct {
	Image             string      `json:"image,omitempty"`
	Port              int32       `json:"port,omitempty`
	SharedVolume      string      `json:"sharedVolume,omitempty"`
	PrivateVolumeSize string      `json:"privateVolumeSize,omitempty"`
	Args              []string    `json:"args,omitempty"`
	Input             []InputItem `json:"input,omitempty"`
}

type InputItem struct {
	Name    string `json:"name"`
	Caption string `json:"caption"`
}

type InputValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
