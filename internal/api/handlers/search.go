package handlers

import (
	"net/http"

	"github.com/zh0glikk/elastic-app/internal/api/requests"
	"github.com/zh0glikk/elastic-app/internal/service/data"
)

func SearchItems(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewSearchRequest(r)
	if err != nil {
		renderErr(w, BadRequest(err))
		return
	}

	items, err := data.NewItemsQ(ElasticCli(r)).Search(request.Key, request.Limit, request.Offset)
	if err != nil {
		Log(r).WithError(err).Error("failed to exec search query")
		renderErr(w, InternalError())
		return
	}

	render(w, items)
}
