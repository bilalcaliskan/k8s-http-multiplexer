package metrics

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestRunMetricsServer(t *testing.T) {
	errChan := make(chan error, 1)

	go func() {
		router := mux.NewRouter()
		err := RunMetricsServer(router)
		errChan <- err
	}()

	select {
	case c := <-errChan:
		t.Error(c)
	case <-time.After(10 * time.Second):
		_, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", config.MetricsPort, config.MetricsUri))
		if err != nil {
			log.Fatalln(err)
		}

		t.Log("success")
		return
	}
}
