package main

import (
	db_connection "consumer/pkg/db_connection"
	memcache "consumer/pkg/memcache"
	"log"
	"time"
)

func main() {
	// 1) Подключение к БД
	db, err := db_connection.DB_Connect()
	if err != nil {
		log.Fatalf("Ошибка соединения с базой данных: %v", err)
	}

	// 2) Инициализация кеша cache с последующим кешированием данных из БД
	var cache = memcache.New_cache(2*time.Minute, 10*time.Minute)
	cache.Input(db)
}
