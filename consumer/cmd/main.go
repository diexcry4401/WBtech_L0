package main

import (
	connect_nats_streaming "consumer/pkg/connect_nats_streaming"
	db_connection "consumer/pkg/db_connection"
	memcache "consumer/pkg/memcache"
	"consumer/pkg/server"
	"os"
	"os/signal"
	"syscall"

	"log"
	"time"
)

func main() {
	// 1) Подключение к БД
	db, err := db_connection.DB_Connect()
	if err != nil {
		log.Fatalf("Ошибка соединения с базой данных: %v", err)
	}

	// 2) Инициализация кеша cache с последующим импортом данных из БД
	var cache = memcache.New_cache(2*time.Minute, 10*time.Minute)
	cache.Input(db)

	// 3) Подключение к NATS Streaming серверу

	go func() {
		err := connect_nats_streaming.СonnectingNats(db, cache)
		if err != nil {
			log.Fatalf("Ошибка при подключении к NATS: %v", err)
		}
	}()
	log.Println("Consumer запущен. Ожидание сообщений...")

	// 4) Запуск сервера
	go func() {
		err = server.ServerStart(cache)
		if err != nil {
			log.Fatalf("Ошибка при подключении к NATS: %v", err)
		}
	}()
	// Ожидание сигнала для завершения работы приложения
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	log.Println("Consumer завершает работу.")

}
