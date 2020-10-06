package middleware

import (
	"net/http"
	"os"
)

// UsefulHeaders - adds useful headers to the response send
func UsefulHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("access-control-allow-origin", os.Getenv("FRONTEND"))
		w.Header().Set("access-control-allow-methods", "*")
		w.Header().Set("access-control-allow-credentials", "true")
		w.Header().Set("access-control-allow-headers", "*")

		if os.Getenv("APP_ENV") == "production" {
			w.Header().Set("Content-Security-Policy", "default-src 'self';base-uri 'self';block-all-mixed-content;font-src 'self' https: data:;frame-ancestors 'self';img-src 'self' data:;object-src 'none';script-src 'self';script-src-attr 'none';style-src 'self' https: 'unsafe-inline';upgrade-insecure-requests")
		}

		w.Header().Set("X-Robots-Tag", "none")
		w.Header().Set("Expect-CT", "max-age=0")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Strict-Transport-Security", "max-age=15552000; includeSubDomains")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-DNS-Prefetch-Control", "off")
		w.Header().Set("X-Download-Options", "noopen")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}
