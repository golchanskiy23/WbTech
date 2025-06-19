package service

import (
	"calendar/entity"
	"testing"
	"time"
)

func TestCreateEvent_Success(t *testing.T) {
	cache := CreateNewCache()
	start := time.Now()
	end := start.Add(time.Hour)

	event, err := cache.CreateEvent("Title", "Desc", start, end)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if event.ID != 1 || event.Title != "Title" {
		t.Errorf("unexpected event: %+v", event)
	}
	if cache.NextID != 2 {
		t.Errorf("expected NextID = 2, got %d", cache.NextID)
	}
}

func TestCreateEvent_InvalidTime(t *testing.T) {
	cache := CreateNewCache()
	start := time.Now()
	end := start.Add(-time.Hour)

	_, err := cache.CreateEvent("Bad", "Time", start, end)
	if err == nil {
		t.Fatal("expected error for invalid time")
	}
}

func TestUpdateEvent_Success(t *testing.T) {
	cache := CreateNewCache()
	start := time.Now()
	end := start.Add(time.Hour)

	original, _ := cache.CreateEvent("Old", "Desc", start, end)
	updated := original
	updated.Title = "New"

	err := cache.UpdateEvent(original.ID, updated)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cache.Events[original.ID].Title != "New" {
		t.Errorf("event not updated")
	}
}

func TestUpdateEvent_NotFound(t *testing.T) {
	cache := CreateNewCache()
	event := entity.Event{ID: 42, Title: "Missing"}

	err := cache.UpdateEvent(event.ID, event)
	if err == nil {
		t.Fatal("expected error for non-existent event")
	}
}

func TestDeleteEvent_Success(t *testing.T) {
	cache := CreateNewCache()
	ev, _ := cache.CreateEvent("ToDelete", "desc", time.Now(), time.Now().Add(time.Hour))

	err := cache.DeleteEvent(ev.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := cache.Events[ev.ID]; ok {
		t.Errorf("event not deleted")
	}
}

func TestDeleteEvent_NotFound(t *testing.T) {
	cache := CreateNewCache()

	err := cache.DeleteEvent(999)
	if err == nil {
		t.Fatal("expected error for missing ID")
	}
}

func TestGetEventsForDay(t *testing.T) {
	cache := CreateNewCache()
	today := time.Date(2025, 6, 19, 10, 0, 0, 0, time.UTC)
	cache.CreateEvent("Today", "Event", today, today.Add(time.Hour))
	cache.CreateEvent("Other", "Event", today.AddDate(0, 0, 1), today.AddDate(0, 0, 1).Add(time.Hour))

	events := cache.GetEventsForDay(today)
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

func TestGetEventsForWeek(t *testing.T) {
	cache := CreateNewCache()
	weekStart := time.Date(2025, 6, 17, 0, 0, 0, 0, time.UTC)
	cache.CreateEvent("MidWeek", "1", weekStart.AddDate(0, 0, 2), weekStart.AddDate(0, 0, 2).Add(time.Hour))
	cache.CreateEvent("OutsideWeek", "2", weekStart.AddDate(0, 0, 8), weekStart.AddDate(0, 0, 8).Add(time.Hour))

	events := cache.GetEventsForWeek(weekStart)
	if len(events) != 1 {
		t.Errorf("expected 1 event in week, got %d", len(events))
	}
}

func TestGetEventsForMonth(t *testing.T) {
	cache := CreateNewCache()
	cache.CreateEvent("JuneEvent", "J", time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC), time.Date(2025, 6, 10, 1, 0, 0, 0, time.UTC))
	cache.CreateEvent("JulyEvent", "J", time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 7, 1, 1, 0, 0, 0, time.UTC))

	events := cache.GetEventsForMonth(2025, 6)
	if len(events) != 1 {
		t.Errorf("expected 1 event in June, got %d", len(events))
	}
}
