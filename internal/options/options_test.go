package options

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetK8sHttpMultiplexerOptions(t *testing.T) {
	t.Log("fetching default options.K8sHttpMultiplexerOptions")
	opts := GetK8sHttpMultiplexerOptions()
	assert.NotNil(t, opts)
	t.Logf("fetched default options.K8sHttpMultiplexerOptions, %v\n", opts)
}
