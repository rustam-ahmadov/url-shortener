package redirect

import (
	"fmt"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func New(log *slog.Logger, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "*")
		if alias == "" {
			log.Error("alias is empty")
			storage.Log("alias is empty", slog.LevelError)
			render.JSON(w, r, response.Error("alias is empty"))
		}

		url, err := storage.GetURL(alias)
		if err != nil {
			errStr := fmt.Sprintf("url has not been found by alias: %s", alias)
			log.Error(errStr)
			storage.Log(errStr, slog.LevelError)
			render.JSON(w, r, response.Error("not found"))
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	}
}
