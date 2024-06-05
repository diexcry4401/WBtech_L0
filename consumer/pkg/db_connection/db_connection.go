package db_connection

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"pass"`
	Dbname   string `json:"db_name"`
}

func DB_Connect() (*sql.DB, error) {
	// Записываем путь до текущей директории в dir
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка получения пути до директории:", err)
		return nil, err
	}
	// Чтение конфига для подключения к Posgres
	configFile, err := os.Open(filepath.Join(dir, "pkg/db_connection/db_config.json"))
	if err != nil {
		fmt.Println("Ошибка чтения конфига:", err)
		return nil, err
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Ошибка чтения данных конфига:", err)
		return nil, err
	}
	var connectionStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Dbname)
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return nil, err
	} else {
		log.Println("Успешное подключение к PostgresSQL по порту:", config.Port)
		return db, nil
	}
}
