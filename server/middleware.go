package main

import "net/http"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow requests from our React frontend
		w.Header().Set("Access-Control-Allow-Origin", 
		"https://go-payroll-processor-system.vercel.app/")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Browser sends OPTIONS request first to check if CORS is allowed
		// We just respond with 200 and stop here
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
