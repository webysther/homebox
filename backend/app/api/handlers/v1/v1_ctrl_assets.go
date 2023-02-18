package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/hay-kot/homebox/backend/internal/core/services"
	"github.com/hay-kot/homebox/backend/internal/data/repo"
	"github.com/hay-kot/homebox/backend/internal/sys/validate"
	"github.com/hay-kot/homebox/backend/pkgs/server"

	"github.com/rs/zerolog/log"
)

// HandleItemGet godocs
// @Summary  Gets an item by Asset ID
// @Tags     Assets
// @Produce  json
// @Param    id  path     string true "Asset ID"
// @Success  200       {object} repo.PaginationResult[repo.ItemSummary]{}
// @Router   /v1/assets/{id} [GET]
// @Security Bearer
func (ctrl *V1Controller) HandleAssetGet() server.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := services.NewContext(r.Context())
		assetIdParam := chi.URLParam(r, "id")
		assetIdParam = strings.ReplaceAll(assetIdParam, "-", "") // Remove dashes
		// Convert the asset ID to an int64
		assetId, err := strconv.ParseInt(assetIdParam, 10, 64)
		if err != nil {
			return err
		}
		pageParam := r.URL.Query().Get("page")
		var page int64 = -1
		if pageParam != "" {
			page, err = strconv.ParseInt(pageParam, 10, 64)
			if err != nil {
				return server.Respond(w, http.StatusBadRequest, "Invalid page number")
			}
		}

		pageSizeParam := r.URL.Query().Get("pageSize")
		var pageSize int64 = -1
		if pageSizeParam != "" {
			pageSize, err = strconv.ParseInt(pageSizeParam, 10, 64)
			if err != nil {
				return server.Respond(w, http.StatusBadRequest, "Invalid page size")
			}
		}

		items, err := ctrl.repo.Items.QueryByAssetID(r.Context(), ctx.GID, repo.AssetID(assetId), int(page), int(pageSize))
		if err != nil {
			log.Err(err).Msg("failed to get item")
			return validate.NewRequestError(err, http.StatusInternalServerError)
		}
		return server.Respond(w, http.StatusOK, items)
	}
}