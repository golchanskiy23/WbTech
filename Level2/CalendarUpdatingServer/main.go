package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/create_event", LoggingMiddleware(CreateEventHandler))
	http.HandleFunc("/update_event", LoggingMiddleware(UpdateEventHandler))
	http.HandleFunc("/events_for_day", LoggingMiddleware(GetEventsForDayHandler))
	http.HandleFunc("/events_for_week", LoggingMiddleware(GetEventsForWeekHandler))
	http.HandleFunc("/events_for_month", LoggingMiddleware(GetEventsForMonthHandler))
	http.HandleFunc("/delete_event", LoggingMiddleware(DeleteEventHandler))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
