package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
)

func StructToJson(v any) []byte {
	jsonData, err := json.Marshal(v)
	if err != nil {
		fmt.Println("Error encoding to JSON:", err)
	}
	return jsonData
}
func ToJson(v any) []byte {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Ошибка:", err)

	}
	return data
}
func WriteAnswer(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(ToJson(data))
}
func GenerateRequestid(n int) (string, error) {
	// Алфавит из заглавных, строчных букв и цифр
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)
	for i := range bytes {
		// Генерируем случайный индекс в пределах длины алфавита
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		bytes[i] = charset[num.Int64()]
	}
	return string(bytes), nil
}
