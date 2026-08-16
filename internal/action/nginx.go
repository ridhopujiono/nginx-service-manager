package action

import (
	"encoding/json"
	"fmt"
	"net/http"

	nginxmanager "nginx-manager-service/internal/nginx"
)

func TestNginxHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := nginxmanager.TestConfig(); err != nil {
		writeJSON(
			w,
			http.StatusUnprocessableEntity,
			map[string]any{
				"success": false,
				"status":  "invalid",
				"message": fmt.Sprintf(
					"Nginx configuration is invalid: %v",
					err,
				),
			},
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		map[string]any{
			"success": true,
			"status":  "valid",
			"message": "Nginx configuration is valid.",
		},
	)
}

func ReloadNginxHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	/*
		Jangan reload config yang invalid.
	*/

	if err := nginxmanager.TestConfig(); err != nil {
		writeJSON(
			w,
			http.StatusUnprocessableEntity,
			map[string]any{
				"success": false,
				"status":  "invalid_config",
				"message": fmt.Sprintf(
					"Nginx configuration is invalid: %v",
					err,
				),
			},
		)

		return
	}

	/*
		Config valid, reload.
	*/

	if err := nginxmanager.Reload(); err != nil {
		writeJSON(
			w,
			http.StatusInternalServerError,
			map[string]any{
				"success": false,
				"status":  "reload_failed",
				"message": fmt.Sprintf(
					"Nginx reload failed: %v",
					err,
				),
			},
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		map[string]any{
			"success": true,
			"status":  "reloaded",
			"message": "Nginx reloaded successfully.",
		},
	)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}