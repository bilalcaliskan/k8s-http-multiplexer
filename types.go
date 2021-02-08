package main

type Header struct {
	Key string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type Config struct {
	Port int `yaml:"port"`
	ReadTimeoutSeconds int `yaml:"readTimeoutSeconds"`
	WriteTimeoutSeconds int `yaml:"writeTimeoutSeconds"`
	MasterUrl string `yaml:"masterUrl"`
	InCluster bool `yaml:"inCluster"`
	Requests []struct{
		Method string `yaml:"method"`
		URI string `yaml:"uri"`
		Namespace string `yaml:"namespace,omitempty"`
		Label string `yaml:"label"`
		BasicAuth bool `yaml:"basicAuth"`
		ReturnResponpseBody bool `yaml:"returnResponpseBody,omitempty"`
		Headers []Header `yaml:"headers"`
		Username string `yaml:"username,omitempty"`
		Password string `yaml:"password,omitempty"`
		ReturnResponseBody bool `yaml:"returnResponseBody,omitempty"`
	} `yaml:"requests"`
}