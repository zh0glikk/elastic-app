package requests

import (
	"net/http"
	"strconv"
)

type SearchRequest struct {
	Key    string
	Limit  int
	Offset int
}

func NewSearchRequest(r *http.Request) (*SearchRequest, error) {
	var req SearchRequest

	values := r.URL.Query()

	keys, ok := values["key"]
	if ok {
		req.Key = keys[0]
	}
	limits, ok := values["limit"]
	if ok {
		limit, err := strconv.ParseInt(limits[0], 10, 64)
		if err != nil {
			return nil, err
		}
		req.Limit = int(limit)
	}

	offsets, ok := values["offset"]
	if ok {
		offset, err := strconv.ParseInt(offsets[0], 10, 64)
		if err != nil {
			return nil, err
		}
		req.Offset = int(offset)
	}

	return &req, nil
}
