package main

import (
	"encoding/json"

	"github.com/nats-io/nats.go"

	"log"
	"os"
)

type Order struct {
	OrderUID        string   `json:"order_uid"`
	TrackNumber     string   `json:"track_number"`
	Entry           string   `json:"entry"`
	Delivery        Delivery `json:"delivery"`
	Payment         Payment  `json:"payment"`
	Items           []Item   `json:"items"`
	Locale          string   `json:"locale"`
	InternalSign    string   `json:"internal_signature"`
	CustomerID      string   `json:"customer_id"`
	DeliveryService string   `json:"delivery_service"`
	ShardKey        string   `json:"shardkey"`
	SMID            int      `json:"sm_id"`
	DateCreated     string   `json:"date_created"`
	OOFShard        string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	// Подключение к серверу NATS Streaming
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Ошибка подключения к NATS Streaming: %v", err)
	}
	defer nc.Close()

	data, err := os.ReadFile("model.json")
	if err != nil {
		log.Fatalf("Ошибка чтения JSON: %v", err)
	}

	// Парсинг данных JSON в структуру данных
	var order Order
	err = json.Unmarshal(data, &order)
	if err != nil {
		log.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	// Отправка сообщений
	subject := "test.subject"
	jsonData, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Ошибка сериализации JSON: %v", err)
	}

	err = nc.Publish(subject, jsonData)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	} else {
		log.Println("Отправлены данные:")
		log.Printf("Order UID: %s", order.OrderUID)
		log.Printf("Track Number: %s", order.TrackNumber)
	}
}
