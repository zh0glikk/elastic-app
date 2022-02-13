package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"

	"github.com/zh0glikk/elastic-app/internal/api/handlers"
	"github.com/zh0glikk/elastic-app/internal/api/middlewares"
	"github.com/zh0glikk/elastic-app/internal/config"
)

func Router(log *logrus.Entry, cfg config.Config) (chi.Router, error) {
	r := chi.NewRouter()

	client, err := elastic.NewClient(elastic.SetURL(cfg.ElasticCfg.URL), elastic.SetSniff(false))
	if err != nil {
		return nil, errors.Wrap(err, "failed to init elastic client")
	}

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middlewares.CtxMiddlewares(
			handlers.CtxLog(log),
			handlers.CtxElasticCli(client),
		),
	)

	r.Route("/items", func(r chi.Router) {
		r.Get("/", handlers.SearchItems)
	})


	//FIXME FIXME hardcoded
	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("home.html")
		if err != nil {
			logrus.WithError(err).Error("failed to parse file")
			return
		}

		t.Execute(w, "")
	})

	return r, nil
}
