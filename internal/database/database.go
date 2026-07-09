package database

import (
	"ProductAPI/internal/models"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Database {
	return &Database{
		conn: conn,
	}
}

func (db *Database) QueryProduct(sql string, args ...any) ([]models.Product, bool) {
	var products []models.Product
	rows, err := db.conn.Query(context.Background(), sql, args...)
	if err != nil {
		fmt.Println("EROR GetProducts:", err)
		return products, false
	}
	defer rows.Close()
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Category,
			&product.Stock,
		)
		if err != nil {
			fmt.Println("EROR QueryProduct:", err)
			return products, false
		}
		products = append(products, product)
	}
	return products, true

}
func (db *Database) OrdersDb(sql string, args ...any) ([]models.AllOrders, bool) {
	var orders []models.AllOrders

	rows, err := db.conn.Query(context.Background(), sql, args...)
	if err != nil {
		fmt.Println("ERRO GET SEARCH:", err)
		return orders, false
	}
	defer rows.Close()
	for rows.Next() {
		var order models.AllOrders
		err := rows.Scan(
			&order.UserName,
			&order.ProductName,
			&order.Quantity,
			&order.OrderStatus,
		)
		if err != nil {
			fmt.Println("SCAN ERROR:", err)
			return orders, false
		}
		orders = append(orders, order)
	}
	return orders, true
}
func (db *Database) StatisticProducts() (models.ProductStats, bool) {
	var product models.ProductStats

	err := db.conn.QueryRow(
		context.Background(),
		"SELECT COUNT(*), AVG(price), MAX(price), MIN(price) FROM products",
	).Scan(
		&product.Count,
		&product.Avg,
		&product.Max,
		&product.Min,
	)

	if err != nil {
		fmt.Println("ERROR GET STATS:", err)
		return models.ProductStats{}, false
	}

	return product, true
}
func (db *Database) AddProduct(name string, price float64, category string, stock int) bool {
	_, err := db.conn.Exec(context.Background(), "INSERT INTO products (name, price, category, stock) VALUES ($1, $2, $3, $4)",
		name,
		price,
		category,
		stock)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Insert failed: %v\n", err)
		return false
	}
	return true
}

func (db *Database) DeleteProduct(id int) bool {
	_, err := db.conn.Exec(context.Background(), "DELETE FROM products WHERE id=$1", id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Delete failed: %v\n", err)
		return false
	}
	return true

}

func (db *Database) GetProductOrder() ([]models.Order, bool) {
	var orders []models.Order
	query := `SELECT products.name, orders.quantity 
	FROM orders
	INNER JOIN products
	ON orders.product_id = products.id;
`
	rows, err := db.conn.Query(context.Background(), query)
	if err != nil {
		fmt.Println("ERORO GET SEARCH:", err)
		return orders, false
	}
	defer rows.Close()
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ProductName,
			&order.Quantity,
		)
		if err != nil {
			fmt.Println("SCAN ERROR:", err)
			return orders, false
		}
		orders = append(orders, order)
	}
	return orders, true
}

func (db *Database) StatsOrders() (models.ProductStats, bool) {
	var product models.ProductStats
	query := `SELECT COUNT(*), SUM(orders.quantity), AVG(price), MAX(price), MIN(price) 
	FROM products JOIN orders
	ON orders.product_id = products.id;`
	err := db.conn.QueryRow(
		context.Background(), query,
	).Scan(
		&product.Sum,
		&product.Count,
		&product.Avg,
		&product.Max,
		&product.Min,
	)

	if err != nil {
		fmt.Println("ERROR GET STATS:", err)
		return models.ProductStats{}, false
	}

	return product, true
}
func (db *Database) StatsUsersOrders() ([]models.UserOrder, bool) {
	var products []models.UserOrder

	rows, err := db.conn.Query(context.Background(), "SELECT users.name, orders.id, orders.quantity FROM orders INNER JOIN users ON orders.user_id = users.id;")
	if err != nil {
		fmt.Println("ERORO GET SEARCH:", err)
		return products, false
	}
	defer rows.Close()
	for rows.Next() {
		var product models.UserOrder
		err := rows.Scan(
			&product.UserName,
			&product.OrderID,
			&product.Quantity,
		)
		if err != nil {
			fmt.Println("SCAN ERROR:", err)
			return products, false
		}
		products = append(products, product)
	}
	return products, true
}
func (db *Database) CategoryCount() ([]models.CategoryCount, bool) {
	var Counts []models.CategoryCount
	rows, err := db.conn.Query(context.Background(), "SELECT category, COUNT(*) FROM products GROUP BY category;")
	if err != nil {
		fmt.Println("ERORO GET SEARCH:", err)
		return Counts, false
	}
	defer rows.Close()
	for rows.Next() {
		var count models.CategoryCount
		err = rows.Scan(
			&count.CategoryName,
			&count.CategoryCount,
		)
		if err != nil {
			fmt.Println("SCAN ERROR:", err)
			return Counts, false
		}
		Counts = append(Counts, count)
	}
	if err != nil {
		fmt.Println("ERROR GET STATS:", err)
		return Counts, false
	}

	return Counts, true

}

func (db *Database) UpdateProduct(
	id int,
	name string,
	price float64,
	category string,
	stock int,
) bool {
	_, err := db.conn.Exec(
		context.Background(),
		`UPDATE products
		SET name = $1,
			price = $2,
			category = $3,
			stock = $4
		WHERE id = $5`,
		name,
		price,
		category,
		stock,
		id,
	)
	if err != nil {
		fmt.Println("UPDATE ERROR:", err)
		return false
	}

	return true

}
