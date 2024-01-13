package save

import (
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const aliasLength = 6

type Request struct {
	URL   string `json:"url" validate: "required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req) //unmarshal

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			storage.Log("failed to decode request body", slog.LevelError)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))
			storage.Log("invalid request", slog.LevelError)

			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
		err = storage.SaveURL(req.URL, alias)
		if err != nil {
			existedUrl := slog.String("url", req.URL)
			errMsg := "url already exists" + existedUrl.String()

			log.Info("url already exists", existedUrl)
			storage.Log(errMsg, slog.LevelInfo)

			render.JSON(w, r, resp.Error("alias already exists"))
			return
		}
		log.Info("url added")
		storage.Log("url added", slog.LevelInfo)

		responseOK(w, r, alias)
	}

}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
