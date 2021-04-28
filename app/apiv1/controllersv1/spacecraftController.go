package controllersv1

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/fittipaldi/go-r3-d3/app/apiv1/modelsv1"
	"github.com/fittipaldi/go-r3-d3/lib/helper"

	"github.com/gorilla/context"
	"gorm.io/gorm"
)

func GetSpacecrafts(w http.ResponseWriter, r *http.Request) {
	gormDB := context.Get(r, "gormDB").(*gorm.DB)

	getParams := r.URL.Query()
	page := 1
	limit := 5
	offset := 0
	pageParam, _ := strconv.Atoi(strings.Join(getParams["page"], ""))
	if pageParam > 0 {
		page = pageParam
	}
	if page > 1 {
		offset = limit * (page - 1)
	}

	var spacecraft modelsv1.Spacecraft
	var spacecrafts []modelsv1.Spacecraft
	var count int64

	//main query
	query := gormDB.Model(&spacecraft)
	//total
	query.Count(&count)
	//pagination
	query.Offset(offset).Limit(limit).Scan(&spacecrafts)
	//total pages
	totalPages := math.Ceil(float64(count) / float64(limit))

	payload := map[string]interface{}{
		"status":         true,
		"items":          spacecrafts,
		"total":          count,
		"total_pages":    totalPages,
		"items_per_page": limit,
	}
	helper.JsonResponse(w, http.StatusOK, payload)
	return
}
