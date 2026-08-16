package main

import (
	"log"
	"net/http"
	"os"

	"nginx-manager-service/internal/proxy"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /health",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			w.WriteHeader(
				http.StatusOK,
			)

			w.Write(
				[]byte(`{"status":"ok"}`),
			)
		},
	)

	mux.HandleFunc(
		"POST /api/v1/proxies",
		proxy.CreateHandler,
	)

	mux.HandleFunc(
		"GET /api/v1/proxies",
		proxy.ListHandler,
	)

	mux.HandleFunc(
		"DELETE /api/v1/proxies/{domain}",
		proxy.DeleteHandler,
	)

	addr := os.Getenv("LISTEN_ADDR")

	if addr == "" {
		addr = "127.0.0.1:9001"
	}

	log.Printf(
		"nginx-manager-service listening on %s",
		addr,
	)

	if err := http.ListenAndServe(
		addr,
		mux,
	); err != nil {

		log.Fatal(err)
	}
}