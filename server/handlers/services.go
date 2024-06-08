package handlers

import (
	"encoding/json"
	"net/http"
)

type ServicesResp struct {
	Services []string `json:"services"`
	Message  string   `json:"message"`
	Status   string   `json:"status"`
}

var s = [][]string{{"sales", "purchases", "inventory"}, {"sales", "purchases"}, {"purchases"}}

const NUM_SERVICE_OPTIONS = 3

var c = 0
var failCounter = 0

func Services(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only HTTP GET request is allowed"))
		return
	}
	idx := c % NUM_SERVICE_OPTIONS
	c += 1
	failCounter += 1
	svcs := &ServicesResp{
		Services: s[idx],
		Message:  "",
		Status:   "success",
	}
	data, err := json.Marshal(svcs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Marshalling services unsuccessful response failed: " + err.Error()))
		return
	}

	if (failCounter % 2) != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error encountered while generating services"))
		return
	}
	w.Write([]byte(data))
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
}
