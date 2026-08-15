package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	handler "kos-api/api"
	"kos-api/internal/config"
	"kos-api/internal/repository/postgres"
	"kos-api/internal/service"
)

func main() {
	appConfig, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	nodeRepository, err := postgres.NewNodeRepository(appConfig.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nodeRepository.Close()

	nodeService := service.NewNodeService(nodeRepository)
	mux := http.NewServeMux()
	mux.HandleFunc("/api", handler.Handler)
	mux.Handle("/api/nodes", handler.NewNodesHandler(nodeService, appConfig.DebugErrors))
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
