package database

import (
	"Level0/internal/entity"
	"Level0/pkg/postgres"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	DeliveryLabel = "D"
	PaymentLabel  = "P"
	ItemLabel     = "I"
)

type CRUDRepository interface {
	GetAllOrders(ctx context.Context) ([]entity.Order, error)
	AddOrder(ctx context.Context, order entity.Order) error
}

type DatabaseRepository struct {
	DB *postgres.DatabaseSource
}

func unmarshalling(src []byte, order *entity.Order, str, label string) error {
	if len(src) > 0 {
		var err error
		switch label {
		case "D":
			err = json.Unmarshal(src, &order.Delivery)
		case "P":
			err = json.Unmarshal(src, &order.Payment)
		case "I":
			err = json.Unmarshal(src, &order.Items)
		default:
			return fmt.Errorf("unknown unmarshal target")
		}

		if err != nil {
			return fmt.Errorf("%s : %v", str, err)
		}
		return nil
	}
	return errors.New("no data found")
}

func CreateNewDBRepository(db *postgres.DatabaseSource) DatabaseRepository {
	return DatabaseRepository{DB: db}
}

func (r DatabaseRepository) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
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

	rows, err := r.DB.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying orders: %w", err)
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var o entity.Order
		var deliveryJSON, paymentJSON, itemsJSON []byte

		err = rows.Scan(
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
			return nil, fmt.Errorf("scan error: %w", err)
		}

		sliceBytes := [][]byte{deliveryJSON, paymentJSON, itemsJSON}
		phrases := []string{"delivery parse error", "payment parse error", "items parse error"}
		labels := []string{DeliveryLabel, PaymentLabel, ItemLabel}
		for i := 0; i < 3; i++ {
			if err = unmarshalling(sliceBytes[i], &o, phrases[i], labels[i]); err != nil {
				return nil, err
			}
		}

		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
	/*rows, err := r.DB.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("incorrect getting all order")
	}
	defer rows.Close()

	var orders []entity.Order

	for rows.Next() {
		var o entity.Order
		var deliveryJSON, paymentJSON, itemsJSON []byte

		err = rows.Scan(
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
			return nil, fmt.Errorf("scan error, something wrong")
		}

		sliceBytes := [][]byte{deliveryJSON, paymentJSON, itemsJSON}
		phrases := []string{"delivery parse error", "payment parse error", "items parse error"}
		labels := []string{DeliveryLabel, PaymentLabel, ItemLabel}
		for i := 0; i < 3; i++ {
			if err = unmarshalling(sliceBytes[i], &o, phrases[i], labels[i]); err != nil {
				return nil, err
			}
		}

		orders = append(orders, o)
	}

	return orders, nil*/
}

func (r DatabaseRepository) AddOrder(ctx context.Context, order entity.Order) error {
	if err := postgres.AddOrdersToDB(r.DB, order); err != nil {
		return fmt.Errorf("error adding order: %v", err)
	}
	return nil
}
