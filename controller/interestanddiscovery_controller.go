package controller

import (
	"cometosee/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type DiscoveryController struct {
	service service.DiscoveryService
}

func NewDiscoveryController(s service.DiscoveryService) *DiscoveryController {
	return &DiscoveryController{service: s}
}

func (c *DiscoveryController) DiscoverUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed ", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	lat, _ := strconv.ParseFloat(query.Get("lat"), 64)
	lon, _ := strconv.ParseFloat(query.Get("lon"), 64)
	radius, _ := strconv.Atoi(query.Get("radius"))
	sport := query.Get("sport")
	skill := query.Get("skill")
	userId, _ := strconv.Atoi(query.Get("user_id"))

	users, err := c.service.DiscoverUsers(lat, lon, radius, sport, skill, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
