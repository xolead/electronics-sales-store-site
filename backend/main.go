package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"electronic/pkg/handlers"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range")

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

	log.Fatal(http.ListenAndServe(":8080", router))

}
