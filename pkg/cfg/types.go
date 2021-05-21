package cfg

// Config is the representation of parsed config file
type Config struct {
	Port                int       `yaml:"port"`
	MetricsPort         int       `yaml:"metricsPort"`
	MetricsUri          string    `yaml:"metricsUri"`
	MasterUrl           string    `yaml:"masterUrl"`
	ReadTimeoutSeconds  int       `yaml:"readTimeoutSeconds"`
	WriteTimeoutSeconds int       `yaml:"writeTimeoutSeconds"`
	InCluster           bool      `yaml:"inCluster"`
	Requests            []Request `yaml:"requests"`
}

// GetRequest method returns the specific Request struct from Config
func (c Config) GetRequest(method, uri string) (bool, Request) {
	for _, v := range c.Requests {
		if v.Method == method && v.URI == uri {
			return true, v
		}
	}
	return false, Request{}
}

// Request is the representation of incoming request to application
type Request struct {
	Method               string   `yaml:"method"`
	URI                  string   `yaml:"uri"`
	Label                string   `yaml:"label"`
	TargetPort           int32    `yaml:"targetPort"`
	BasicAuth            bool     `yaml:"basicAuth"`
	Headers              []Header `yaml:"headers"`
	Username             string   `yaml:"username,omitempty"`
	Password             string   `yaml:"password,omitempty"`
	ReturnResponseBody   bool     `yaml:"returnResponseBody,omitempty"`
	RequestBody          string   `yaml:"requestBody,omitempty"`
	ExpectedResponseCode int      `yaml:"expectedResponseCode"`
}

// Header is the representation of required header keys and values in Request
type Header struct {
	Key   string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
}

// Response is the representation of the json response to multiplexed requests
type Response struct {
	TargetCount  int `json:"targetCount"`
	SuccessCount int `json:"successCount"`
}
