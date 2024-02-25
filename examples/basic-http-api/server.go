package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		panic("server port not provided")
	}

	port := os.Args[1]
	instanceId := os.Args[2]

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/users", func(w http.ResponseWriter, r *http.Request) {
		log.Println("received req", r.URL, r.Method, r.RemoteAddr)
		w.Header().Set("content-type", "application/json")
		io.WriteString(w, fmt.Sprintf(`{
			"instance": "%s",
			"users": [{
				"id": 0,
				"name": "user1"
			}, {
				"id": 1,
				"name": "user2"
			}]
		}`, instanceId))
	})

	err := http.ListenAndServe(fmt.Sprintf("localhost:%s", port), mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Println("server could not be started", err)
		os.Exit(1)
	}
}
