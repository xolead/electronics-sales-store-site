package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	cloudstorage "electronic/pkg/cloud_storage"
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
	// Настройка CORS для S3 при каждом запуске
	s3Config := cloudstorage.LoadConfig()

	// Настраиваем CORS каждый раз при запуске
	err := cloudstorage.SetupS3CORS(s3Config)
	if err != nil {
		log.Printf("⚠️ Failed to setup S3 CORS: %v", err)
	} else {
		log.Println("✅ S3 CORS configured on startup")
	}

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
