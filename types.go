package main

type TargetPod struct {
	addr string
	label string
}

func (targetPod *TargetPod) Equals(other *TargetPod) bool {
	isAddrEquals := targetPod.addr == other.addr
	isLabelEquals := targetPod.label == other.label
	return isAddrEquals && isLabelEquals
}

type Header struct {
	Key string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type Request struct {
	Method string `yaml:"method"`
	URI string `yaml:"uri"`
	Namespace string `yaml:"namespace"`
	Label string `yaml:"label"`
	TargetPort int32 `yaml:"targetPort"`
	BasicAuth bool `yaml:"basicAuth"`
	Headers []Header `yaml:"headers"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	ReturnResponseBody bool `yaml:"returnResponseBody,omitempty"`
	RequestBody string `yaml:"requestBody,omitempty"`
	ExpectedResponseCode int `yaml:"expectedResponseCode"`
}

type Config struct {
	Port int `yaml:"port"`
	ReadTimeoutSeconds int `yaml:"readTimeoutSeconds"`
	WriteTimeoutSeconds int `yaml:"writeTimeoutSeconds"`
	MasterUrl string `yaml:"masterUrl"`
	InCluster bool `yaml:"inCluster"`
	Requests []Request `yaml:"requests"`
}