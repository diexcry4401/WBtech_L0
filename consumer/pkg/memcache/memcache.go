package memcache

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"consumer/pkg/order_structure"

	_ "github.com/lib/pq"
)

func New_cache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Items)
	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}
	if cleanupInterval > 0 {
		cache.StartGC()
	}
	return &cache
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {

	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()
	defer c.Unlock()

	c.items[key] = Items{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}

}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	if !found {
		return nil, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}

	return item.Value, true

}

func (c *Cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("key not found")
	}

	delete(c.items, key)
	return nil

}

func (c *Cache) Input(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Итерация по результатам запроса и заполнение кеша
	for rows.Next() {
		var order order_structure.Order
		var id_order int
		err := rows.Scan(
			&id_order,
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSign,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SMID,
			&order.DateCreated,
			&order.OOFShard)
		if err != nil {
			panic(err)
		}
		var delivery order_structure.Delivery

		err = db.QueryRow("SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_id = $1", id_order).Scan(
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
		)
		if err != nil {
			panic(err)
		}
		order.Delivery = delivery
		var payment order_structure.Payment

		err = db.QueryRow("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_id = $1", id_order).Scan(
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
		)
		if err != nil {
			panic(err)
		}

		order.Payment = payment

		rows_item, err := db.Query("SELECT chrt_id, track_number, price,rid,  name,sale,size,total_price, nm_id, brand, status FROM items WHERE order_id = $1", id_order)
		if err != nil {
			panic(err)
		}
		defer rows_item.Close()

		var itemsOrder []order_structure.Item

		// Итерация по результатам запроса и заполнение кеша
		for rows_item.Next() {
			var item order_structure.Item
			err := rows_item.Scan(
				&item.ChrtID,
				&item.TrackNumber,
				&item.Price,
				&item.RID,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.TotalPrice,
				&item.NmID,
				&item.Brand,
				&item.Status)
			if err != nil {
				panic(err)
			}
			itemsOrder = append(itemsOrder, item)
		}
		order.Items = itemsOrder

		c.Set(fmt.Sprintf("%d", id_order), order, 0) // Предположим, что id является уникальным
	}

	return nil
}

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Items
}

type Items struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

func (c *Cache) clearItems(keys []string) {
	c.Lock()
	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}

func (c *Cache) expiredKeys() (keys []string) {
	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}
	return
}

func (c *Cache) GC() {
	for {
		<-time.After(c.cleanupInterval)
		if c.items == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

func (c *Cache) StartGC() {
	go c.GC()
}
