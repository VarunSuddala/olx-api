package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

func Get_listings(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`select id,title,description,price,city,created_at from listings
		order by created_at Desc
		limit 100`)
		if err != nil {
			log.Printf("Query : %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		listings := []listing{}
		for rows.Next() {
			l := listing{}
			if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
				log.Printf("rows.err %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			listings = append(listings, l)
		}

		if err := rows.Err(); err != nil {
			log.Printf("rows.err : %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(listings)
	}
}
