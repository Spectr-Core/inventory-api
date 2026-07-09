package handlers

import (
	"ProductAPI/internal/database"
	"ProductAPI/internal/models"
	"ProductAPI/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Server struct {
	db *database.Database
}

func New(db *database.Database) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) ProdutsHandler(w http.ResponseWriter, req *http.Request) {
	limitStr := req.URL.Query().Get("limit")
	offsetStr := req.URL.Query().Get("offset")
	sortParam := req.URL.Query().Get("sort")
	categoryParam := req.URL.Query().Get("category")

	query := "SELECT * FROM products"
	args := []any{}

	if categoryParam != "" {
		query += " WHERE category = $1"
		args = append(args, categoryParam)
	}

	if sortParam == "price_asc" {
		query += " ORDER BY price ASC"
	} else if sortParam == "price_desc" {
		query += " ORDER BY price DESC"
	} else if sortParam != "" {

	}

	if limitStr != "" || offsetStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}

		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		args = append(args, limit, offset)
	}

	product, ok := s.db.QueryProduct(query, args...)
	if !ok {
		http.Error(w, "error get products", http.StatusInternalServerError)
		return
	}
	utils.WriteAnswer(w, product)
}
func (s *Server) Product(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		id := req.URL.Query().Get("id")
		productid, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		data, ok := s.db.QueryProduct("SELECT id, name, price, category, stock FROM products WHERE id=$1", productid)
		if !ok {
			http.Error(w, "Product Not Found:", http.StatusBadRequest)
			return
		}
		utils.WriteAnswer(w, data)
	case http.MethodPost:
		var Product models.Product
		err := json.NewDecoder(req.Body).Decode(&Product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok := s.db.AddProduct(Product.Name, Product.Price, Product.Category, Product.Stock)
		if ok {
			io.WriteString(w, "ADDED")
		} else {
			http.Error(w, "Not Added:", http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		id := req.URL.Query().Get("id")
		productid, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		ok := s.db.DeleteProduct(productid)
		if ok {
			io.WriteString(w, "DELETE")
		} else {
			http.Error(w, "Not DELETE:", http.StatusBadRequest)
			return
		}
	case http.MethodPut:
		id := req.URL.Query().Get("id")
		productid, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		var product models.Product

		json.NewDecoder(req.Body).Decode(&product)
		ok := s.db.UpdateProduct(productid, product.Name, product.Price, product.Category, product.Stock)
		if !ok {
			http.Error(w, "Not Update:", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Method Not Alloved", http.StatusBadRequest)
		return

	}
}
func (s *Server) ProductStats(w http.ResponseWriter, req *http.Request) {
	route := req.URL.Path
	switch route {
	case "/products/stats":
		data, ok := s.db.StatisticProducts()
		if !ok {
			http.Error(w, "ERROR GET PRODUCT STATS:", http.StatusBadRequest)
			return
		}
		utils.WriteAnswer(w, data)

	case "/products/expensive":
		data, ok := s.db.QueryProduct("SELECT * FROM products ORDER BY price DESC LIMIT 5")
		if !ok {
			http.Error(w, "EROR GET EXSPANIVE PRODUCT:", http.StatusBadRequest)
			return
		}
		utils.WriteAnswer(w, data)
	case "/products/cheap":
		data, ok := s.db.QueryProduct("SELECT * FROM products ORDER BY price ASC LIMIT 5")
		if !ok {
			http.Error(w, "EROR GET POOR PRODUCT:", http.StatusBadRequest)
			return
		}
		utils.WriteAnswer(w, data)
	case "/products/category/stats":
		data, ok := s.db.CategoryCount()
		if !ok {
			http.Error(w, "EROR GET CATEGORY STATS:", http.StatusBadRequest)
			return
		}
		utils.WriteAnswer(w, data)
	default:
		http.Error(w, "Get not Found", http.StatusNotFound)
		return
	}

}
func (s *Server) SearchHendler(w http.ResponseWriter, req *http.Request) {
	route := req.URL.Path
	switch route {
	case "/products/search":
		queryParams := req.URL.Query()

		name := queryParams.Get("name")
		if name == "" {
			http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
			return
		}
		pattern := "%" + name + "%"
		data, ok := s.db.QueryProduct("SELECT * FROM products WHERE name LIKE $1", pattern)
		if !ok {
			http.Error(w, "ERROR GET PRODUCT SEARCH:", http.StatusBadRequest)
			return
		}
		utils.WriteAnswer(w, data)
	case "/products/category":
		queryParams := req.URL.Query()

		name := queryParams.Get("name")
		if name == "" {
			http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
			return
		}

		data, ok := s.db.QueryProduct("SELECT * FROM products WHERE category = $1", name)
		if !ok {
			http.Error(w, "ERROR GET SEARCH CATEGORY:", http.StatusBadRequest)
			return
		}
		utils.WriteAnswer(w, data)
	default:
		http.Error(w, "Get not Found", http.StatusNotFound)
		return

	}

}
func (s *Server) Orders(w http.ResponseWriter, req *http.Request) {
	route := req.URL.Path
	status := req.URL.Query().Get("status")

	switch route {
	case "/orders/users":
		data, ok := s.db.StatsUsersOrders()
		if !ok {
			http.Error(w, "ERROR GET ORDERS USERS:", http.StatusBadRequest)
			return
		}
		utils.WriteAnswer(w, data)
	case "/orders/product":
		data, ok := s.db.GetProductOrder()
		if !ok {
			http.Error(w, "ERROR GET ORDERS:", http.StatusBadRequest)
			return
		}
		utils.WriteAnswer(w, data)
	case "/orders":
		if status != "" {
			query := `
		    SELECT users.name, products.name, orders.quantity, orders.status 
			FROM orders INNER JOIN users 
			ON orders.user_id = users.id 
			INNER JOIN products 
			ON orders.product_id = products.id
			WHERE orders.status = $1;`
			data, ok := s.db.OrdersDb(query, status)
			if !ok {
				http.Error(w, "ERROR GET ALL ORDERS:", http.StatusBadRequest)
				return
			}
			utils.WriteAnswer(w, data)
		} else {
			query := `
		SELECT users.name, products.name, orders.quantity, orders.status 
			FROM orders INNER JOIN users 
			ON orders.user_id = users.id 
			INNER JOIN products 
			ON orders.product_id = products.id; `
			data, ok := s.db.OrdersDb(query)
			if !ok {
				http.Error(w, "ERROR GET ALL ORDERS:", http.StatusBadRequest)
				return
			}
			utils.WriteAnswer(w, data)
		}
	case "/orders/paid":
		query := `
		SELECT users.name, products.name, orders.quantity, orders.status 
		FROM orders INNER JOIN users 
		ON orders.user_id = users.id 
		INNER JOIN products 
		ON orders.product_id = products.id 
		WHERE orders.status = 'paid';`
		data, ok := s.db.OrdersDb(query)
		if !ok {
			http.Error(w, "ERROR GET ALL ORDERS PAID:", http.StatusBadRequest)
			return
		}
		utils.WriteAnswer(w, data)
	case "/orders/user":
		id := req.URL.Query().Get("id")
		productid, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		query := `SELECT users.name AS user_name, products.name AS product_name, orders.quantity, orders.status 
		FROM orders JOIN users ON orders.user_id = users.id 
		JOIN products ON orders.product_id = products.id 
		WHERE orders.user_id = $1
	`
		data, ok := s.db.OrdersDb(query, productid)
		if !ok {
			http.Error(w, "EROR WITH ORDERS BY USER:", http.StatusBadRequest)
			return
		}
		utils.WriteAnswer(w, data)
	case "/orders/product/id":
		id := req.URL.Query().Get("id")
		productid, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		query := `SELECT users.name AS user_name, products.name AS product_name, orders.quantity, orders.status 
		FROM orders JOIN users ON orders.user_id = users.id 
		JOIN products ON orders.product_id = products.id 
		WHERE orders.product_id = $1`
		data, ok := s.db.OrdersDb(query, productid)
		if !ok {
			http.Error(w, "EROR SEACH PRODUCT BY ID:", http.StatusBadRequest)
			return
		}
		utils.WriteAnswer(w, data)
	case "/orders/stats":
		data, ok := s.db.StatsOrders()
		if !ok {
			http.Error(w, "EROR ORDERS STATS:", http.StatusBadRequest)
			return
		}
		utils.WriteAnswer(w, data)
	default:
		http.Error(w, "Get not Found", http.StatusNotFound)
		return

	}
}
