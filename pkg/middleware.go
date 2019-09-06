package pkg

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var isDev = os.Getenv("DEV") == "true"

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		tic := time.Now()
		statusRes := statusWriter{res, 200}
		next.ServeHTTP(res, req)
		if statusRes.status >= 400 || isDev {
			ns := float64(time.Now().Sub(tic).Nanoseconds())
			message := fmt.Sprintf(
				"[%s] %s %s (%s ms) %d",
				req.RemoteAddr,
				req.Method,
				req.URL,
				fmt.Sprintf("%.3f", ns/1000000.0),
				statusRes.status,
			)
			fmt.Println(message)
		}
	})
}

func methodFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		method := req.Method
		if !(method == "GET" || method == "POST" || method == "OPTIONS") {
			res.WriteHeader(400)
			res.Write([]byte("Method must be a GET, POST, or OPTIONS"))
			return
		}

		next.ServeHTTP(res, req)
	})
}

func allowCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// origin := req.Header.Get("Origin")
		if os.Getenv("DEV") == "true" {
			res.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			res.Header().Set("Access-Control-Allow-Origin", "*") // TODO - Figure out origins for prod
		}
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if req.Method == "OPTIONS" {
			res.WriteHeader(200)
			return
		}
		next.ServeHTTP(res, req)
	})
}
