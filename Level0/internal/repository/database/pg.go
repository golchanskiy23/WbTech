package database

import (
	"Level0/internal/entity"
	"Level0/pkg/postgres"
	"context"
	"encoding/json"
	"log"
)

type DatabaseRepository struct {
	DB *postgres.DatabaseSource
}

// скрывает базу данных
func CreateNewDBRepository(db *postgres.DatabaseSource) *DatabaseRepository {
	return &DatabaseRepository{DB: db}
}

// набор методов, реализующих работу с базами данных
func (r *DatabaseRepository) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	// скорее всего временно
	query := `SELECT
  o.*,
  (
    SELECT json_build_object(
      'name', d.name,
      'phone', d.phone,
      'zip', d.zip,
      'city', d.city,
      'address', d.address,
      'region', d.region,
      'email', d.email
    )
    FROM deliveries d
    WHERE d.order_uid = o.order_uid
  ) AS delivery,
  (
    SELECT json_build_object(
      'transaction', p.transaction,
      'request_id', p.request_id,
      'currency', p.currency,
      'provider', p.provider,
      'amount', p.amount,
      'payment_dt', p.payment_dt,
      'bank', p.bank,
      'delivery_cost', p.delivery_cost,
      'goods_total', p.goods_total,
      'custom_fee', p.custom_fee
    )
    FROM payments p
    WHERE p.order_uid = o.order_uid
  ) AS payment,
  (
    SELECT json_agg(
      json_build_object(
        'chrt_id', i.chrt_id,
        'track_number', i.track_number,
        'price', i.price,
        'rid', i.rid,
        'name', i.name,
        'sale', i.sale,
        'size', i.size,
        'total_price', i.total_price,
        'nm_id', i.nm_id,
        'brand', i.brand,
        'status', i.status
      )
    )
    FROM items i
    WHERE i.order_uid = o.order_uid
  ) AS items
FROM orders o;
`
	rows, err := r.DB.Pool.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var orders []entity.Order

	for rows.Next() {
		var o entity.Order
		var deliveryJSON, paymentJSON, itemsJSON []byte

		err := rows.Scan(
			&o.OrderUID,
			&o.TrackNumber,
			&o.Entry,
			&o.Locale,
			&o.InternalSignature,
			&o.CustomerID,
			&o.DeliveryService,
			&o.ShardKey,
			&o.SmID,
			&o.DataCreated,
			&o.OofShard,
			&deliveryJSON,
			&paymentJSON,
			&itemsJSON,
		)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}

		if len(deliveryJSON) > 0 {
			if err := json.Unmarshal(deliveryJSON, &o.Delivery); err != nil {
				log.Println("Delivery parse error:", err)
			}
		}
		if len(paymentJSON) > 0 {
			if err := json.Unmarshal(paymentJSON, &o.Payment); err != nil {
				log.Println("Payment parse error:", err)
			}
		}
		if len(itemsJSON) > 0 {
			if err := json.Unmarshal(itemsJSON, &o.Items); err != nil {
				log.Println("Items parse error:", err)
			}
		}

		orders = append(orders, o)
	}

	return orders, nil
}
