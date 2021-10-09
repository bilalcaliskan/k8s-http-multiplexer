package web

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestRunWebServer(t *testing.T) {
	errChan := make(chan error, 1)

	go func() {
		router := mux.NewRouter()
		webServer := RunWebServer(router)
		errChan <- webServer
	}()

	for {
		select {
		case c := <-errChan:
			t.Error(c)
		case <-time.After(10 * time.Second):
			_, err := http.Get(fmt.Sprintf("http://localhost:%d/ping", config.Port))
			if err != nil {
				log.Fatalln(err)
			}

			t.Log("success")
			return
		}
	}
}
