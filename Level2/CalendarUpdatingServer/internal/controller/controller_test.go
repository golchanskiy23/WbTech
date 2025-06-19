package controller

import (
	"bytes"
	"calendar/entity"
	"calendar/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockCache struct {
	createFunc   func(string, string, time.Time, time.Time) (entity.Event, error)
	updateFunc   func(int, entity.Event) error
	deleteFunc   func(int) error
	getDayFunc   func(time.Time) []entity.Event
	getWeekFunc  func(time.Time) []entity.Event
	getMonthFunc func(int, time.Month) []entity.Event
}

func (m *mockCache) CreateEvent(t, d string, s, e time.Time) (entity.Event, error) {
	if m.createFunc != nil {
		return m.createFunc(t, d, s, e)
	}
	return entity.Event{}, nil
}
func (m *mockCache) UpdateEvent(id int, e entity.Event) error {
	if m.updateFunc != nil {
		return m.updateFunc(id, e)
	}
	return nil
}
func (m *mockCache) DeleteEvent(id int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}
func (m *mockCache) GetEventsForDay(day time.Time) []entity.Event {
	if m.getDayFunc != nil {
		return m.getDayFunc(day)
	}
	return nil
}
func (m *mockCache) GetEventsForWeek(start time.Time) []entity.Event {
	if m.getWeekFunc != nil {
		return m.getWeekFunc(start)
	}
	return nil
}
func (m *mockCache) GetEventsForMonth(y int, mth time.Month) []entity.Event {
	if m.getMonthFunc != nil {
		return m.getMonthFunc(y, mth)
	}
	return nil
}

func setupMockCache(f func(string, string, time.Time, time.Time) (entity.Event, error)) {
	service.Cache = &mockCache{createFunc: f}
}

func TestCreateEventHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/create_event", nil)
	rr := httptest.NewRecorder()

	CreateEventHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestCreateEventHandler_BadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBuffer([]byte("not json")))
	rr := httptest.NewRecorder()

	CreateEventHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCreateEventHandler_BadStartTime(t *testing.T) {
	body := `{"title":"T","description":"D","start":"bad","end":"2025-06-20T10:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	CreateEventHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCreateEventHandler_BadEndTime(t *testing.T) {
	body := `{"title":"T","description":"D","start":"2025-06-20T09:00:00Z","end":"bad"}`
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	CreateEventHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCreateEventHandler_CreateEventFails(t *testing.T) {
	setupMockCache(func(title, desc string, start, end time.Time) (entity.Event, error) {
		return entity.Event{}, errors.New("creation failed")
	})

	input := entity.URLEvent{
		Title:       "T",
		Description: "D",
		Start:       "2025-06-20T09:00:00Z",
		End:         "2025-06-20T10:00:00Z",
	}
	data, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBuffer(data))
	rr := httptest.NewRecorder()

	CreateEventHandler(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rr.Code)
	}
}

func TestCreateEventHandler_Success(t *testing.T) {
	expectedEvent := entity.Event{
		ID:          1,
		Title:       "T",
		Description: "D",
		Start:       time.Date(2025, 6, 20, 9, 0, 0, 0, time.UTC),
		End:         time.Date(2025, 6, 20, 10, 0, 0, 0, time.UTC),
	}

	setupMockCache(func(title, desc string, start, end time.Time) (entity.Event, error) {
		return expectedEvent, nil
	})

	input := entity.URLEvent{
		Title:       "T",
		Description: "D",
		Start:       expectedEvent.Start.Format(time.RFC3339),
		End:         expectedEvent.End.Format(time.RFC3339),
	}
	data, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBuffer(data))
	rr := httptest.NewRecorder()

	CreateEventHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}

	var got entity.Event
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if got.ID != expectedEvent.ID || got.Title != expectedEvent.Title {
		t.Errorf("unexpected response: %+v", got)
	}
}

func TestDeleteEventHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/delete_event?id=1", nil)
	rr := httptest.NewRecorder()

	DeleteEventHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestDeleteEventHandler_MissingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/delete_event", nil)
	rr := httptest.NewRecorder()

	DeleteEventHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestDeleteEventHandler_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/delete_event?id=abc", nil)
	rr := httptest.NewRecorder()

	DeleteEventHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestDeleteEventHandler_NotFound(t *testing.T) {
	setupMockCache(nil)
	service.Cache.(*mockCache).deleteFunc = func(id int) error {
		return errors.New("not found")
	}

	req := httptest.NewRequest(http.MethodPost, "/delete_event?id=1", nil)
	rr := httptest.NewRecorder()

	DeleteEventHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestDeleteEventHandler_Success(t *testing.T) {
	setupMockCache(nil)
	service.Cache.(*mockCache).deleteFunc = func(id int) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodPost, "/delete_event?id=1", nil)
	rr := httptest.NewRecorder()

	DeleteEventHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetEventsForDayHandler(t *testing.T) {
	setupMockCache(nil)
	date := "2025-06-20"
	service.Cache.(*mockCache).getDayFunc = func(day time.Time) []entity.Event {
		return []entity.Event{{ID: 1, Title: "Event"}}
	}

	req := httptest.NewRequest(http.MethodGet, "/events_for_day?day="+date, nil)
	rr := httptest.NewRecorder()

	GetEventsForDayHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetEventsForDayHandler_BadQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/events_for_day?day=invalid", nil)
	rr := httptest.NewRecorder()

	GetEventsForDayHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestGetEventsForWeekHandler(t *testing.T) {
	setupMockCache(nil)
	date := "2025-06-17"
	service.Cache.(*mockCache).getWeekFunc = func(day time.Time) []entity.Event {
		return []entity.Event{{ID: 1, Title: "WeekEvent"}}
	}

	req := httptest.NewRequest(http.MethodGet, "/events_for_week?week="+date, nil)
	rr := httptest.NewRecorder()

	GetEventsForWeekHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetEventsForMonthHandler(t *testing.T) {
	setupMockCache(nil)
	month := "2025-06"
	service.Cache.(*mockCache).getMonthFunc = func(year int, m time.Month) []entity.Event {
		return []entity.Event{{ID: 1, Title: "MonthEvent"}}
	}

	req := httptest.NewRequest(http.MethodGet, "/events_for_month?month="+month, nil)
	rr := httptest.NewRecorder()

	GetEventsForMonthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetEventsForMonthHandler_BadMonth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/events_for_month?month=bad", nil)
	rr := httptest.NewRecorder()

	GetEventsForMonthHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestUpdateEventHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/update_event", nil)
	rr := httptest.NewRecorder()

	UpdateEventHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestUpdateEventHandler_BadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/update_event", bytes.NewBufferString("notjson"))
	rr := httptest.NewRecorder()

	UpdateEventHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestUpdateEventHandler_NotFound(t *testing.T) {
	setupMockCache(nil)
	service.Cache.(*mockCache).updateFunc = func(id int, e entity.Event) error {
		return errors.New("not found")
	}

	e := entity.Event{ID: 1, Title: "Updated"}
	data, _ := json.Marshal(e)
	req := httptest.NewRequest(http.MethodPost, "/update_event", bytes.NewBuffer(data))
	rr := httptest.NewRecorder()

	UpdateEventHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestUpdateEventHandler_Success(t *testing.T) {
	setupMockCache(nil)
	service.Cache.(*mockCache).updateFunc = func(id int, e entity.Event) error {
		return nil
	}

	e := entity.Event{ID: 1, Title: "Updated"}
	data, _ := json.Marshal(e)
	req := httptest.NewRequest(http.MethodPost, "/update_event", bytes.NewBuffer(data))
	rr := httptest.NewRecorder()

	UpdateEventHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}
