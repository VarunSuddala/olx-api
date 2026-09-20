package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
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

type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db:     db,
		logger: logger,
	}
}

func (lh ListingHandler) Get_listings(w http.ResponseWriter, r *http.Request) {

	// request scoped context
	ctx := r.Context()
	rows, err := lh.db.QueryContext(ctx,
		`select id,title,description,price,city,created_at from listings
		order by created_at Desc
		limit 100`)
	if err != nil {
		lh.logger.Error("lisitng error", "err", err)
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

func (lh ListingHandler) Delete_listing(w http.ResponseWriter, r *http.Request) {
	// generate request id
	
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"msg":"id is required"}`)
		return
	}
	_, err := lh.db.ExecContext(ctx,
		`delete from listings where id = $1`, id)

	if err != nil {

		lh.logger.Error("delete failed", "listing_id", id, "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (lh ListingHandler) post_listing(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("posting_listing : %v", err)
		http.Error(w, "failed to read body", http.StatusBadRequest)
	}
	fmt.Println(string(body))
}
