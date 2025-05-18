package controller

import (
	"Level0/internal/service"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"reflect"
)

type Controller struct {
	Service *service.OrderService
}

func (c Controller) GetOrderByIdHandler(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	fmt.Println(orderID, reflect.TypeOf(orderID))
	order, err := c.Service.GetOrderById(orderID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := fmt.Sprintf("Don't have order with such ID in database: %v", err)
		w.Write([]byte(resp))
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		ans, err := json.Marshal(&order)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := fmt.Sprintf("Error during write response in database: %v", err)
			w.Write([]byte(resp))
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(ans)
		}
	}
}

func CreateNewOrderController(service *service.OrderService) *chi.Mux {
	controller := Controller{
		Service: service,
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/order/{id}", controller.GetOrderByIdHandler)
	return r
}
