package controller

import (
	"calendar/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func GetEventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	day, err := parseDateQuery(r, "day")
	if err != nil {
		http.Error(w, "Invalid day query parameter", http.StatusBadRequest)
		return
	}

	events := service.Cache.GetEventsForDay(day)
	json.NewEncoder(w).Encode(events)
}

func GetEventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	startOfWeek, err := parseDateQuery(r, "week")
	if err != nil {
		http.Error(w, "Invalid week query parameter", http.StatusBadRequest)
		return
	}

	events := service.Cache.GetEventsForWeek(startOfWeek)
	json.NewEncoder(w).Encode(events)
}

func GetEventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	year, month, err := parseMonthQuery(r, "month")
	if err != nil {
		http.Error(w, "Invalid month query parameter", http.StatusBadRequest)
		return
	}

	events := service.Cache.GetEventsForMonth(year, month)
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
