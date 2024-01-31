package showall

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	resp "myURLShortener/internal/lib/api/response"
	"myURLShortener/internal/lib/logger/slogger"
	"myURLShortener/internal/storage"
	"net/http"
)

type Response struct {
	resp.Response
	Data []*AliasUrl `json:"data,omitempty"`
}

type AliasUrl struct {
	Alias string `json:"alias,omitempty"`
	URL   string `json:"url,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLShower
type URLShower interface {
	ShowURL() ([]*AliasUrl, error)
}

func New(log *slog.Logger, urlShower URLShower) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.show.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		resURLS, err := urlShower.ShowURL()
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("Url not found")

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("Failed to get url", slogger.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("Showed url", slog.Any("urls", resURLS))

		responseOK(w, r, resURLS)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, aliasUrl []*AliasUrl) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Data:     aliasUrl,
	})
}
