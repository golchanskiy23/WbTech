package utils

import (
	"Level0/internal/entity"
	"math/rand"
	"strconv"
	"time"
)

func RandomOrder() entity.Order {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return entity.Order{
		OrderUID:          randString(r, "order_", 8),
		TrackNumber:       randString(r, "track_", 6),
		Entry:             randomEntry(r),
		Delivery:          randomDelivery(r),
		Payment:           randomPayment(r),
		Items:             randomItems(r, r.Intn(3)+1),
		Locale:            "en",
		InternalSignature: randString(r, "", 10),
		CustomerID:        randString(r, "cust_", 5),
		DeliveryService:   randomDeliveryService(r),
		ShardKey:          strconv.Itoa(r.Intn(100)),
		SmID:              r.Intn(1000),
		DataCreated:       time.Now().Format(time.RFC3339),
		OofShard:          strconv.Itoa(r.Intn(10)),
	}
}

func randString(r *rand.Rand, prefix string, n int) string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	res := make([]byte, n)
	for i := range res {
		res[i] = chars[r.Intn(len(chars))]
	}
	return prefix + string(res)
}

func randomEntry(r *rand.Rand) string {
	entries := []string{"WBIL", "WBILM", "WB"}
	return entries[r.Intn(len(entries))]
}

func randomDeliveryService(r *rand.Rand) string {
	services := []string{"dpd", "cdek", "boxberry", "dhl"}
	return services[r.Intn(len(services))]
}

func randomDelivery(r *rand.Rand) entity.Delivery {
	return entity.Delivery{
		Name:    "John Doe",
		Phone:   "+79001234567",
		Zip:     strconv.Itoa(100000 + r.Intn(899999)),
		City:    "City" + randString(r, "", 3),
		Address: "Street " + strconv.Itoa(r.Intn(100)),
		Region:  "Region " + strconv.Itoa(r.Intn(10)),
		Email:   randString(r, "user", 4) + "@example.com",
	}
}

func randomPayment(r *rand.Rand) entity.Payment {
	return entity.Payment{
		Transaction:  randString(r, "tx_", 12),
		RequestID:    randString(r, "req_", 6),
		Currency:     "USD",
		Provider:     "visa",
		Amount:       int32(1000 + r.Intn(9000)),
		PaymentDT:    time.Now().Unix(),
		Bank:         "SomeBank",
		DeliveryCost: int32(r.Intn(500)),
		GoodsTotal:   int32(r.Intn(5000)),
		CustomFee:    int32(r.Intn(100)),
	}
}

func randomItems(r *rand.Rand, n int) []entity.Items {
	items := make([]entity.Items, n)
	for i := 0; i < n; i++ {
		items[i] = entity.Items{
			ChrtID:      r.Int63(),
			TrackNumber: randString(r, "track_", 6),
			Price:       int32(100 + r.Intn(1000)),
			Rid:         randString(r, "rid_", 6),
			Name:        "Item " + strconv.Itoa(i+1),
			Sale:        int32(r.Intn(100)),
			Size:        []string{"S", "M", "L", "XL"}[r.Intn(4)],
			TotalPrice:  int32(200 + r.Intn(1500)),
			NmID:        r.Int63(),
			Brand:       "Brand" + strconv.Itoa(r.Intn(100)),
			Status:      r.Intn(10),
		}
	}
	return items
}
