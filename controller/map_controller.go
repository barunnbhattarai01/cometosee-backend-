package controller

import (
	"cometosee/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type MapController struct {
	mapService *service.MapService
}

func NewMapController(mapService *service.MapService) *MapController {
	return &MapController{mapService: mapService}
}

func (c *MapController) GetMapEventPins(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	latStr := q.Get("lat")
	if latStr == "" {
		http.Error(w, "lat is required", http.StatusBadRequest)
		return
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "invalid lat", http.StatusBadRequest)
		return
	}

	lonStr := q.Get("lon")
	if lonStr == "" {
		http.Error(w, "lon is required", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "invalid lon", http.StatusBadRequest)
		return
	}

	radius := 0
	if rStr := q.Get("radius"); rStr != "" {
		radius, err = strconv.Atoi(rStr)
		if err != nil || radius < 0 {
			http.Error(w, "invalid radius", http.StatusBadRequest)
			return
		}
	}

	skill := q.Get("skill")

	pins, err := c.mapService.GetMapEventPins(r.Context(), lat, lon, radius, skill)
	if err != nil {
		http.Error(w, "failed to fetch map pins", http.StatusInternalServerError)
		return
	}

	if pins == nil {
		pins = []map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pins": pins,
	})
}
