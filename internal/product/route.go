package product

import (
	"net/http"
	"strconv"

	"github.com/flukis/pagination-comparation/utils/response"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Router struct {
	svcProduct Service
}

func NewRouter(
	svcProduct Service,
) *Router {
	return &Router{svcProduct}
}

func (r *Router) Routes() *chi.Mux {
	route := chi.NewMux()

	route.Get("/page", r.GetPaginatePage)
	route.Get("/cursor", r.GetPaginateCursor)

	return route
}

func (r *Router) GetPaginatePage(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	limitStr := req.URL.Query().Get("per_page")
	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		if err = response.WriteError(w, http.StatusBadRequest, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	pageStr := req.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		if err = response.WriteError(w, http.StatusBadRequest, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	sortBy := req.URL.Query().Get("sort_by")
	sortOrdered := req.URL.Query().Get("sort_order")
	res := r.svcProduct.GetPaginationWithPage(ctx, pageInt, limitInt, sortBy, sortOrdered)
	if res.Error != nil {
		if err = response.WriteError(w, http.StatusInternalServerError, res.Error); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}

	var metaResp struct {
		Limit    int `json:"limit"`
		ThisPage int `json:"total_data"`
		Next     int `json:"next"`
	}

	metaResp.Limit = res.Meta.Limit
	metaResp.Next = res.Meta.Next
	metaResp.ThisPage = res.Meta.Total
	if err = response.WriteResponse(w, "get all products success", http.StatusOK, res.Data, metaResp); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) GetPaginateCursor(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	limitStr := req.URL.Query().Get("limit")
	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		if err = response.WriteError(w, http.StatusBadRequest, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	cursor := req.URL.Query().Get("cursor")
	res := r.svcProduct.GetPaginationWithCursor(ctx, cursor, limitInt)
	if res.Error != nil {
		if err = response.WriteError(w, http.StatusInternalServerError, res.Error); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}

	var metaResp struct {
		Limit    int    `json:"limit"`
		ThisPage int    `json:"total_data"`
		Next     string `json:"next_cursor"`
	}

	metaResp.Limit = res.Meta.Limit
	metaResp.Next = res.Meta.Next
	metaResp.ThisPage = res.Meta.Total
	if err = response.WriteResponse(w, "get all products success", http.StatusOK, res.Data, metaResp); err != nil {
		log.Error().Err(err)
		return
	}
}
