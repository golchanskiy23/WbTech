package main

import (
	"errors"
	"fmt"
	"time"
)

var Cache = CreateNewCache()

func CreateNewCache() *EventsCache {
	return &EventsCache{
		events: make(map[int]Event),
		nextID: 1,
	}
}

func (c *EventsCache) CreateEvent(title, description string, start, end time.Time) (Event, error) {
	c.Lock()
	defer c.Unlock()

	if start.After(end) {
		return Event{}, errors.New("start time must be before end time")
	}

	event := Event{
		ID:          c.nextID,
		Title:       title,
		Description: description,
		Start:       start,
		End:         end,
	}

	c.events[event.ID] = event
	c.nextID++

	return event, nil
}

func (c *EventsCache) UpdateEvent(eventID int, updatedEvent Event) error {
	c.Lock()
	defer c.Unlock()

	if _, exists := c.events[eventID]; !exists {
		return fmt.Errorf("event with ID %d does not exist", eventID)
	}

	c.events[eventID] = updatedEvent
	return nil
}

func (c *EventsCache) DeleteEvent(eventID int) error {
	c.Lock()
	defer c.Unlock()

	if _, exists := c.events[eventID]; !exists {
		return fmt.Errorf("event with ID %d does not exist", eventID)
	}

	delete(c.events, eventID)
	return nil
}

func (c *EventsCache) GetEventsForDay(day time.Time) []Event {
	c.RLock()
	defer c.RUnlock()

	var eventsForDay []Event
	for _, event := range c.events {
		if sameDay(event.Start, day) {
			eventsForDay = append(eventsForDay, event)
		}
	}
	return eventsForDay
}

func (c *EventsCache) GetEventsForWeek(startOfWeek time.Time) []Event {
	c.RLock()
	defer c.RUnlock()

	var eventsForWeek []Event
	endOfWeek := startOfWeek.AddDate(0, 0, 7)
	for _, event := range c.events {
		if event.Start.After(startOfWeek) && event.Start.Before(endOfWeek) {
			eventsForWeek = append(eventsForWeek, event)
		}
	}
	return eventsForWeek
}

func (c *EventsCache) GetEventsForMonth(year int, month time.Month) []Event {
	c.RLock()
	defer c.RUnlock()

	var eventsForMonth []Event
	for _, event := range c.events {
		if event.Start.Year() == year && event.Start.Month() == month {
			eventsForMonth = append(eventsForMonth, event)
		}
	}
	return eventsForMonth
}

func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
