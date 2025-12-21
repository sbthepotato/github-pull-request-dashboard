package web_pkg

import (
	"context"
	"database/sql"
	"encoding/json"
	"github-pull-request-dashboard/db_pkg"
	"io"
	"net/http"
)

var cachedRegexJson []byte

func GetTitleRegexList(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		setHeaders(&w, "json")

		if cachedRegexJson == nil {
			cachedRegexList, err := db_pkg.GetTitleRegexList(ctx, db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			cachedRegexJson, err = json.Marshal(cachedRegexList)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Write(cachedRegexJson)

	}
}

func SetTitleRegex(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w, "text")

		mu.Lock()
		defer mu.Unlock()

		if r.Method != http.MethodPost {
			http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		titleRegexList := make([]*db_pkg.TitleRegex, 0)

		err = json.Unmarshal(body, &titleRegexList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = db_pkg.UpsertTitleRegex(ctx, db, titleRegexList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cachedRegexList, err := db_pkg.GetTitleRegexList(ctx, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cachedRegexJson, err = json.Marshal(cachedRegexList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Title regex data saved successfully"))

	}
}
