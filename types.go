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
	MetricsPort int `yaml:"metricsPort"`
	MetricsUri string `yaml:"metricsUri"`
	MasterUrl string `yaml:"masterUrl"`
	ReadTimeoutSeconds int `yaml:"readTimeoutSeconds"`
	WriteTimeoutSeconds int `yaml:"writeTimeoutSeconds"`
	InCluster bool `yaml:"inCluster"`
	Requests []Request `yaml:"requests"`
}

func (c Config) getRequest(method, uri string) (bool, Request) {
	for _, v := range c.Requests {
		if v.Method == method && v.URI == uri {
			return true, v
		}
	}
	return false, Request{}
}

type Response struct {
	TargetCount int `json:"targetCount"`
	SuccessCount int `json:"successCount"`
}