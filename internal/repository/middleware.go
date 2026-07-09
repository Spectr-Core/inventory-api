package repository

import (
	"ProductAPI/internal/utils"
	"log"
	"net/http"
	"time"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqid := ""

		cookie, err := r.Cookie("request_id")
		if err == http.ErrNoCookie {
			id, errr := utils.GenerateRequestid(10)
			if errr != nil {
				log.Println("EROR WITH GenerateRequestid LOGGER")
			}
			cookieset := http.Cookie{
				Name:    "request_id",
				Value:   id,
				Path:    "/",
				Expires: time.Now().Add(24 * time.Hour),
			}
			reqid = id
			http.SetCookie(w, &cookieset)
		} else {
			reqid = cookie.Value
		}
		next(w, r)
		elapsed := time.Since(start)
		log.Printf("%s | %s | %s | %s\n", r.Method, r.URL.Path, reqid, elapsed)
	}
}
