package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var input URLEvent

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	start, err := time.Parse(time.RFC3339, input.Start)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}
	end, err := time.Parse(time.RFC3339, input.End)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	event, err := Cache.CreateEvent(input.Title, input.Description, start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	eventIDStr := r.URL.Query().Get("id")
	if eventIDStr == "" {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	err = Cache.DeleteEvent(eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": "Event deleted"})
}

func GetEventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	day, err := parseDateQuery(r, "day")
	if err != nil {
		http.Error(w, "Invalid day query parameter", http.StatusBadRequest)
		return
	}

	events := Cache.GetEventsForDay(day)
	json.NewEncoder(w).Encode(events)
}

func GetEventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	startOfWeek, err := parseDateQuery(r, "week")
	if err != nil {
		http.Error(w, "Invalid week query parameter", http.StatusBadRequest)
		return
	}

	events := Cache.GetEventsForWeek(startOfWeek)
	json.NewEncoder(w).Encode(events)
}

func GetEventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	year, month, err := parseMonthQuery(r, "month")
	if err != nil {
		http.Error(w, "Invalid month query parameter", http.StatusBadRequest)
		return
	}

	events := Cache.GetEventsForMonth(year, month)
	json.NewEncoder(w).Encode(events)
}

func parseDateQuery(r *http.Request, queryParam string) (time.Time, error) {
	dateStr := r.URL.Query().Get(queryParam)
	if dateStr == "" {
		return time.Time{}, errors.New("query parameter missing")
	}
	return time.Parse("2006-01-02", dateStr)
}

func parseMonthQuery(r *http.Request, queryParam string) (int, time.Month, error) {
	monthStr := r.URL.Query().Get(queryParam)
	if monthStr == "" {
		return 0, 0, errors.New("query parameter missing")
	}
	date, err := time.Parse("2006-01", monthStr)
	if err != nil {
		return 0, 0, err
	}
	return date.Year(), date.Month(), nil
}

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

	err = Cache.UpdateEvent(updatedEvent.ID, updatedEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": "Event updated"})
}

func parseAndUpdateEvent(r *http.Request) (Event, error) {
	var updatedEvent Event
	err := json.NewDecoder(r.Body).Decode(&updatedEvent)
	if err != nil {
		return Event{}, err
	}
	return updatedEvent, nil
}
