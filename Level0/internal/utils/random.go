package utils

import (
	"Level0/internal/entity"
	"math/rand"
	"strconv"
	"time"
)

func RandomOrder() entity.Order {
	rand.Seed(time.Now().UnixNano())

	return entity.Order{
		OrderUID:          randString("order_", 8),
		TrackNumber:       randString("track_", 6),
		Entry:             randomEntry(),
		Delivery:          randomDelivery(),
		Payment:           randomPayment(),
		Items:             randomItems(rand.Intn(3) + 1),
		Locale:            "en",
		InternalSignature: randString("", 10),
		CustomerID:        randString("cust_", 5),
		DeliveryService:   randomDeliveryService(),
		ShardKey:          strconv.Itoa(rand.Intn(100)),
		SmID:              rand.Intn(1000),
		DataCreated:       time.Now().Format(time.RFC3339),
		OofShard:          strconv.Itoa(rand.Intn(10)),
	}
}

func randString(prefix string, n int) string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	res := make([]byte, n)
	for i := range res {
		res[i] = chars[rand.Intn(len(chars))]
	}
	return prefix + string(res)
}

func randomEntry() string {
	entries := []string{"WBIL", "WBILM", "WB"}
	return entries[rand.Intn(len(entries))]
}

func randomDeliveryService() string {
	services := []string{"dpd", "cdek", "boxberry", "dhl"}
	return services[rand.Intn(len(services))]
}

func randomDelivery() entity.Delivery {
	return entity.Delivery{
		Name:    "John Doe",
		Phone:   "+79001234567",
		Zip:     strconv.Itoa(100000 + rand.Intn(899999)),
		City:    "City" + randString("", 3),
		Address: "Street " + strconv.Itoa(rand.Intn(100)),
		Region:  "Region " + strconv.Itoa(rand.Intn(10)),
		Email:   randString("user", 4) + "@example.com",
	}
}

func randomPayment() entity.Payment {
	return entity.Payment{
		Transaction:  randString("tx_", 12),
		RequestID:    randString("req_", 6),
		Currency:     "USD",
		Provider:     "visa",
		Amount:       int32(1000 + rand.Intn(9000)),
		PaymentDT:    time.Now().Unix(),
		Bank:         "SomeBank",
		DeliveryCost: int32(rand.Intn(500)),
		GoodsTotal:   int32(rand.Intn(5000)),
		CustomFee:    int32(rand.Intn(100)),
	}
}

func randomItems(n int) []entity.Items {
	items := make([]entity.Items, n)
	for i := 0; i < n; i++ {
		items[i] = entity.Items{
			ChrtID:      rand.Int63(),
			TrackNumber: randString("track_", 6),
			Price:       int32(100 + rand.Intn(1000)),
			Rid:         randString("rid_", 6),
			Name:        "Item " + strconv.Itoa(i+1),
			Sale:        int32(rand.Intn(100)),
			Size:        []string{"S", "M", "L", "XL"}[rand.Intn(4)],
			TotalPrice:  int32(200 + rand.Intn(1500)),
			NmID:        rand.Int63(),
			Brand:       "Brand" + strconv.Itoa(rand.Intn(100)),
			Status:      rand.Intn(10),
		}
	}
	return items
}
