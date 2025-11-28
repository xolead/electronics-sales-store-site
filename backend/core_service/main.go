package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"core-service/pkg/handlers"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем конкретные origin вместо *
		allowedOrigins := []string{
			"http://localhost:3000",
			"https://localhost:3000",
			"http://127.0.0.1:3000",
		}

		origin := r.Header.Get("Origin")
		allowed := false

		for _, o := range allowedOrigins {
			if o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().
			Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, x-amz-acl, x-amz-meta-*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, ETag")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	router := mux.NewRouter()
	router.Use(CORS)
	router.HandleFunc("/product", handlers.CreateProduct).Methods("POST")
	router.HandleFunc("/product", handlers.ReadAllProduct).Methods("GET")
	router.HandleFunc("/product/{id}", handlers.ReadProduct).Methods("GET")
	router.HandleFunc("/product/{id}", handlers.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/product/change", handlers.ChangeCountProduct).Methods("PUT")
	log.Println("Сервер запущен")

	log.Fatal(http.ListenAndServe(":8082", router))

}
