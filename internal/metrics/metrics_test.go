package metrics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestRunMetricsServer(t *testing.T) {
	var conn net.Conn

	defer func() {
		err := conn.Close()
		assert.Nil(t, err)
	}()

	go func() {
		err := RunMetricsServer()
		assert.Nil(t, err)
	}()

	for {
		time.Sleep(1 * time.Second)
		conn, _ = net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", config.MetricsPort), 10*time.Second)
		if conn != nil {
			break
		}
	}

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d%s", config.MetricsPort, config.MetricsUri))
	assert.Nil(t, err)

	_, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
}
