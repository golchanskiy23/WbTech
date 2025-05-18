package addingutils

import (
	"Level0/internal/entity"
	"Level0/pkg/postgres"
	"context"
	"log"
)

func AddOrdersToDB(db *postgres.DatabaseSource, order entity.Order) error {
	var exists bool
	err := db.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid = $1)", order.OrderUID).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	if exists {
		log.Println("Order already exists, skipping insert")
		return err
	}
	_, err = db.Pool.Exec(context.Background(), `
	INSERT INTO orders (
		order_uid, track_number, entry, locale, internal_signature,
		customer_id, delivery_service, shard_key, sm_id, data_created, oof_shard
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DataCreated, order.OofShard)
	if err != nil {
		log.Fatal("insert orders:", err)
		return err
	}

	_, err = db.Pool.Exec(context.Background(), `
	INSERT INTO deliveries (
		order_uid, name, phone, zip, city, address, region, email
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
`, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		log.Fatal("insert deliveries:", err)
		return err
	}

	_, err = db.Pool.Exec(context.Background(), `
	INSERT INTO payments (
		order_uid, transaction, request_id, currency, provider,
		amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
`, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT,
		order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		log.Fatal("insert payments:", err)
		return err
	}

	for _, item := range order.Items {
		_, err := db.Pool.Exec(context.Background(), `
		INSERT INTO items (
			order_uid, chrt_id, track_number, price, rid,
			name, sale, size, total_price, nm_id, brand, status
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			log.Println("insert item:", err)
			return err
		}
	}
	return nil
}
