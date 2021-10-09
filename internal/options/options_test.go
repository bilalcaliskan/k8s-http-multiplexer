package options

import "testing"

func TestGetK8sHttpMultiplexerOptions(t *testing.T) {
	t.Log("fetching default options.K8sHttpMultiplexerOptions")
	opts := GetK8sHttpMultiplexerOptions()
	t.Logf("fetched default options.K8sHttpMultiplexerOptions, %v\n", opts)
}
