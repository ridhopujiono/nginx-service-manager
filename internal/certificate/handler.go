package certificate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	nginxconfig "nginx-manager-service/internal/nginx"
)

var hostnameRegex = regexp.MustCompile(
	`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`,
)

func RenewHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	domain := strings.ToLower(
		strings.TrimSpace(
			r.PathValue("domain"),
		),
	)

	if domain == "" ||
		len(domain) > 253 ||
		!hostnameRegex.MatchString(domain) {

		http.Error(
			w,
			"invalid domain",
			http.StatusBadRequest,
		)

		return
	}

	if err := nginxconfig.RenewCertificate(
		domain,
	); err != nil {

		http.Error(
			w,
			fmt.Sprintf(
				"certificate renewal failed: %v",
				err,
			),
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusOK,
	)

	_ = json.NewEncoder(w).Encode(
		map[string]any{
			"domain": domain,
			"status": "renewal_check_completed",
		},
	)
}