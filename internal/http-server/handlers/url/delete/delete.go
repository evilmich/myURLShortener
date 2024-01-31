package delete

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	resp "myURLShortener/internal/lib/api/response"
	"myURLShortener/internal/lib/logger/slogger"
	"myURLShortener/internal/storage"
	"net/http"
)

type Request struct {
	Alias string `json:"alias"`
	URL   string `json:"url"`
}

type ResponseAlias struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
	URL   string `json:"url,omitempty"`
}

type ResponseURL struct {
	resp.Response
	Aliases []*AliasData `json:"aliases,omitempty"`
	URL     string       `json:"url,omitempty"`
}

type AliasData struct {
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLDeleter
type URLDeleter interface {
	GetURL(alias string) (string, error)
	GetAliasAndURL(alias, url string) (string, string, error)
	DeleteByAliasAndURL(alias, url string) error
	DeleteURLByAlias(alias string) error
	DeleteAliasByURL(url string) ([]*AliasData, error)
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("Request body is empty")

			render.JSON(w, r, resp.Error("Empty request"))

			return
		}
		if err != nil {
			log.Error("Failed to decode request body", slogger.Err(err))

			render.JSON(w, r, resp.Error("Failed to decode request"))

			return
		}

		log.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("Invalid request", slogger.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		url := req.URL

		if alias == "" && url == "" {
			log.Info("Alias and URL is empty")

			render.JSON(w, r, resp.Error("empty params 'alias' and 'url'"))

			return
		} else if alias != "" && url != "" {
			resAlias, resURL, err := urlDeleter.GetAliasAndURL(alias, url)
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("Alias or Url not found", "alias", alias, "url", url)

				render.JSON(w, r, resp.Error("not found"))

				return
			}
			if err != nil {
				log.Error("Failed to get url or alias", slogger.Err(err))

				render.JSON(w, r, resp.Error("internal error"))

				return
			}

			err = urlDeleter.DeleteByAliasAndURL(resAlias, resURL)
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("Alias or Url not found", "alias", alias, "url", url)

				render.JSON(w, r, resp.Error("not found"))

				return
			}
			if err != nil {
				log.Error("Failed to delete url or alias", slogger.Err(err))

				render.JSON(w, r, resp.Error("internal error"))

				return
			}

			log.Info("Alias and Url deleted", slog.String("alias", alias), slog.String("url", url))

			responseOKForAlias(w, r, alias, url)
		} else if url == "" {
			resURL, err := urlDeleter.GetURL(alias)
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("Url not found", "alias", alias)

				render.JSON(w, r, resp.Error("not found"))

				return
			}
			if err != nil {
				log.Error("Failed to get url", slogger.Err(err))

				render.JSON(w, r, resp.Error("internal error"))

				return
			}

			err = urlDeleter.DeleteURLByAlias(alias)
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("Alias not found", "alias", alias)

				render.JSON(w, r, resp.Error("not found"))

				return
			}
			if err != nil {
				log.Error("Failed to delete url", slogger.Err(err))

				render.JSON(w, r, resp.Error("internal error"))

				return
			}

			log.Info("Url deleted", slog.String("url", resURL))

			responseOKForAlias(w, r, alias, resURL)
		} else {
			resAlias, err := urlDeleter.DeleteAliasByURL(url)
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("URL not found", "url", url)

				render.JSON(w, r, resp.Error("not found"))

				return
			}
			if err != nil {
				log.Error("Failed to delete url", slogger.Err(err))

				render.JSON(w, r, resp.Error("internal error"))

				return
			}

			if len(resAlias) == 0 {
				log.Info("URL not found", "url", url)

				render.JSON(w, r, resp.Error("not found"))

				return
			}

			log.Info("Url deleted", slog.String("url", url))

			responseOKForURL(w, r, resAlias, url)

		}

	}
}

func responseOKForAlias(w http.ResponseWriter, r *http.Request, alias string, url string) {
	render.JSON(w, r, ResponseAlias{
		Response: resp.OK(),
		Alias:    alias,
		URL:      url,
	})
}

func responseOKForURL(w http.ResponseWriter, r *http.Request, data []*AliasData, url string) {
	render.JSON(w, r, ResponseURL{
		Response: resp.OK(),
		Aliases:  data,
		URL:      url,
	})
}
