package utils

import (
	"Level0/internal/entity"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Checker = map[string]struct{}

func RandomOrder(c Checker) entity.Order {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderUID := randString(r, c, "order_", 8, true)
	log.Print(orderUID, " ")
	trackNumber := randString(r, c, "track_", 6, false)

	return entity.Order{
		OrderUID:          orderUID,
		TrackNumber:       trackNumber,
		Entry:             randomEntry(r),
		Delivery:          randomDelivery(r, c),
		Payment:           randomPayment(r, orderUID, c),
		Items:             randomItems(r, r.Intn(3)+1, trackNumber, c),
		Locale:            "en",
		InternalSignature: randString(r, c, "", 10, false),
		CustomerID:        randString(r, c, "cust_", 5, false),
		DeliveryService:   randomDeliveryService(r),
		ShardKey:          strconv.Itoa(r.Intn(100)),
		SmID:              r.Intn(1000),
		DataCreated:       time.Now().Format(time.RFC3339),
		OofShard:          strconv.Itoa(r.Intn(10)),
	}
}

func randString(r *rand.Rand, c Checker, prefix string, n int, flag bool) string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	res := make([]byte, n)
	for i := range res {
		res[i] = chars[r.Intn(len(chars))]
	}
	curr := prefix + string(res)
	if _, ok := c[curr]; ok && flag {
		return randString(r, c, prefix, n, flag)
	}
	c[curr] = struct{}{}
	return curr
}

func randomEntry(r *rand.Rand) string {
	entries := []string{"WBIL", "WBILM", "WB"}
	return entries[r.Intn(len(entries))]
}

func randomDeliveryService(r *rand.Rand) string {
	services := []string{"dpd", "cdek", "boxberry", "dhl"}
	return services[r.Intn(len(services))]
}

func randomDelivery(r *rand.Rand, c Checker) entity.Delivery {
	return entity.Delivery{
		Name:    "John Doe",
		Phone:   "+79001234567",
		Zip:     strconv.Itoa(100000 + r.Intn(899999)),
		City:    "City" + randString(r, c, "", 3, false),
		Address: "Street " + strconv.Itoa(r.Intn(100)),
		Region:  "Region " + strconv.Itoa(r.Intn(10)),
		Email:   randString(r, c, "user", 4, false) + "@example.com",
	}
}

func randomPayment(r *rand.Rand, orderUID string, c Checker) entity.Payment {
	return entity.Payment{
		Transaction:  orderUID,
		RequestID:    randString(r, c, "req_", 6, false),
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

func randomItems(r *rand.Rand, n int, t string, c Checker) []entity.Items {
	items := make([]entity.Items, n)
	log.Println(n)
	for i := 0; i < n; i++ {
		items[i] = entity.Items{
			ChrtID:      r.Int63(),
			TrackNumber: t,
			Price:       int32(100 + r.Intn(1000)),
			Rid:         randString(r, c, "rid_", 6, false),
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
