package controller

import (
	"calendar/entity"
	"calendar/internal/service"
	"encoding/json"
	"net/http"
)

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	updatedEvent, err := parseAndUpdateEvent(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = service.Cache.UpdateEvent(updatedEvent.ID, updatedEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": "Event updated"})
}

func parseAndUpdateEvent(r *http.Request) (entity.Event, error) {
	var updatedEvent entity.Event
	err := json.NewDecoder(r.Body).Decode(&updatedEvent)
	if err != nil {
		return entity.Event{}, err
	}
	return updatedEvent, nil
}
