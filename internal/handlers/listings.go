package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/VarunSuddala/olx-api/internal/httpx"
	"github.com/VarunSuddala/olx-api/internal/middleware"
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
type response struct {
	Msg       string `json:"msg"`
	RequestID string `json:"request_id"`
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db:     db,
		logger: logger,
	}
}

func (lh ListingHandler) Get(w http.ResponseWriter, r *http.Request) {

	// request scoped context
	ctx := r.Context()
	rows, err := lh.db.QueryContext(ctx,
		`select id,title,description,price,city,created_at from listings
		order by created_at Desc
		limit 100`)
	if err != nil {
		lh.logger.Error("lisitng error", "err", err)
		httpx.Error(w, 500, "internal error", string(httpx.CodeInternalError))
		return
	}
	defer rows.Close()
	listings := []listing{}
	for rows.Next() {
		l := listing{}
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			log.Printf("rows.err %v", err)
			httpx.Error(w, 500, "internal error", string(httpx.CodeInternalError))
			return
		}
		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows.err : %v", err)
		httpx.Error(w, 500, "internal error", string(httpx.CodeInternalError))

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(listings)
}

func (lh ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestId := middleware.RequestIDContext(ctx)
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"msg":"id is required"}`)
		return
	}
	res, err := lh.db.ExecContext(ctx,
		`delete from listings where id = $1`, id)

	if err != nil {

		lh.logger.Error("delete failed", "listing_id", id, "request_id", requestId, "err", err)
		// http.Error(w, "internal error", http.StatusInternalServerError)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", string(httpx.CodeInternalError))
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		resp := response{
			Msg:       fmt.Sprintf("no listing found with id %s", id),
			RequestID: requestId,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := response{
		Msg:       "deleted successfully",
		RequestID: requestId,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (lh ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request_id := middleware.RequestIDContext(ctx)
	var req listing
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lh.logger.Error("failed to decode", "request_id", request_id, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", string(httpx.CodeMalformedJSON))
		return
	}
	var id string
	row := lh.db.QueryRowContext(ctx,
		`INSERT INTO listings (title, description, price, city)
         VALUES ($1, $2, $3, $4)
         RETURNING id`, req.Title, req.Description, req.Price, req.City)
	if err := row.Scan(&id); err != nil {
		lh.logger.Error("failed to insert", "request_id", request_id, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something wentwrong", string(httpx.CodeInternalError))
		return
	}
	lh.logger.Info("listing created", "request_id", request_id, "listing_id", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id})
}
