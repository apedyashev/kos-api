package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	handler "kos-api/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api", handler.Handler)
	mux.HandleFunc("/api/nodes", handler.NodesHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "KOS API root")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
