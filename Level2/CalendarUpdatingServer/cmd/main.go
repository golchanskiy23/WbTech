package main

import (
	"calendar/internal"
	"calendar/internal/controller"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/create_event", internal.LoggingMiddleware(controller.CreateEventHandler))
	http.HandleFunc("/update_event", internal.LoggingMiddleware(controller.UpdateEventHandler))
	http.HandleFunc("/events_for_day", internal.LoggingMiddleware(controller.GetEventsForDayHandler))
	http.HandleFunc("/events_for_week", internal.LoggingMiddleware(controller.GetEventsForWeekHandler))
	http.HandleFunc("/events_for_month", internal.LoggingMiddleware(controller.GetEventsForMonthHandler))
	http.HandleFunc("/delete_event", internal.LoggingMiddleware(controller.DeleteEventHandler))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
