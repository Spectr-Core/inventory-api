package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ProductAPI/internal/database"
	"ProductAPI/internal/handlers"
	"ProductAPI/internal/repository"

	"github.com/jackc/pgx/v5"
)

func StartServer(conn *pgx.Conn) {

	db := database.New(conn)

	server := handlers.New(db)

	fmt.Println("http://localhost:9090")
	http.HandleFunc("/products", server.ProdutsHandler)
	http.HandleFunc("/product", server.Product)
	http.HandleFunc("/products/stats", server.ProductStats)
	http.HandleFunc("/products/expensive", repository.Logger(server.ProductStats))
	http.HandleFunc("/products/cheap", server.ProductStats)
	http.HandleFunc("/products/category/stats", server.ProductStats)
	http.HandleFunc("/products/search", server.SearchHendler)
	http.HandleFunc("/products/category", server.SearchHendler)
	http.HandleFunc("/orders/users", server.Orders)
	http.HandleFunc("/orders/product", server.Orders)
	http.HandleFunc("/orders", server.Orders)
	http.HandleFunc("/orders/paid", server.Orders)
	http.HandleFunc("/orders/user", server.Orders)
	http.HandleFunc("/orders/product/id", server.Orders)
	http.HandleFunc("/orders/stats", server.Orders)

	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), nil))
}
