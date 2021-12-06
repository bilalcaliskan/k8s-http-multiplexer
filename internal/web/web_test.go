package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestRunWebServer(t *testing.T) {
	var (
		conn net.Conn
		err  error
	)

	defer func() {
		err := conn.Close()
		assert.Nil(t, err)
	}()

	go func() {
		err := RunWebServer()
		assert.Nil(t, err)
	}()

	for {
		time.Sleep(1 * time.Second)
		conn, _ = net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", config.Port), 10*time.Second)
		if conn != nil {
			break
		}
	}

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/ping", config.Port))
	assert.Nil(t, err)

	_, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
}
