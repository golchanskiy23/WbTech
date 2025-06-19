package service

import (
	"calendar/entity"
	"errors"
	"fmt"
	"sync"
	"time"
)

type CacheInterface interface {
	CreateEvent(title, description string, start, end time.Time) (entity.Event, error)
	UpdateEvent(eventID int, updatedEvent entity.Event) error
	DeleteEvent(eventID int) error
	GetEventsForDay(day time.Time) []entity.Event
	GetEventsForWeek(startOfWeek time.Time) []entity.Event
	GetEventsForMonth(year int, month time.Month) []entity.Event
}

type EventsCache struct {
	sync.RWMutex
	Events map[int]entity.Event
	NextID int
}

var Cache CacheInterface = CreateNewCache()

func CreateNewCache() *EventsCache {
	return &EventsCache{
		Events: make(map[int]entity.Event),
		NextID: 1,
	}
}

func (c *EventsCache) CreateEvent(title, description string, start, end time.Time) (entity.Event, error) {
	c.Lock()
	defer c.Unlock()

	if start.After(end) {
		return entity.Event{}, errors.New("start time must be before end time")
	}

	event := entity.Event{
		ID:          c.NextID,
		Title:       title,
		Description: description,
		Start:       start,
		End:         end,
	}

	c.Events[event.ID] = event
	c.NextID++

	return event, nil
}

func (c *EventsCache) UpdateEvent(eventID int, updatedEvent entity.Event) error {
	c.Lock()
	defer c.Unlock()

	if _, exists := c.Events[eventID]; !exists {
		return fmt.Errorf("event with ID %d does not exist", eventID)
	}

	c.Events[eventID] = updatedEvent
	return nil
}

func (c *EventsCache) DeleteEvent(eventID int) error {
	c.Lock()
	defer c.Unlock()

	if _, exists := c.Events[eventID]; !exists {
		return fmt.Errorf("event with ID %d does not exist", eventID)
	}

	delete(c.Events, eventID)
	return nil
}

func (c *EventsCache) GetEventsForDay(day time.Time) []entity.Event {
	c.RLock()
	defer c.RUnlock()

	var eventsForDay []entity.Event
	for _, event := range c.Events {
		if sameDay(event.Start, day) {
			eventsForDay = append(eventsForDay, event)
		}
	}
	return eventsForDay
}

func (c *EventsCache) GetEventsForWeek(startOfWeek time.Time) []entity.Event {
	c.RLock()
	defer c.RUnlock()

	var eventsForWeek []entity.Event
	endOfWeek := startOfWeek.AddDate(0, 0, 7)
	for _, event := range c.Events {
		if event.Start.After(startOfWeek) && event.Start.Before(endOfWeek) {
			eventsForWeek = append(eventsForWeek, event)
		}
	}
	return eventsForWeek
}

func (c *EventsCache) GetEventsForMonth(year int, month time.Month) []entity.Event {
	c.RLock()
	defer c.RUnlock()

	var eventsForMonth []entity.Event
	for _, event := range c.Events {
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
