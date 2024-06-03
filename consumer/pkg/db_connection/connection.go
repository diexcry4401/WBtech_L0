package main

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

	config, err := os.Open(filepath.Join(dir, "db_config.json"))
	if err != nil {
		fmt.Println("Ошибка чтения конфига:", err)
		return nil, err
	}
	defer config.Close()

	var configFile Config
	decoder := json.NewDecoder(config)
	err = decoder.Decode(&configFile)
	if err != nil {
		fmt.Println("Ошибка декодирования конфига:", err)
		return nil, err
	}
	var connectionStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", configFile.Host, configFile.Port, configFile.User, configFile.Password, configFile.Dbname)
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return nil, err
	} else {
		log.Println("Успешное подключение к PostgresSQL по порту:", configFile.Port)
		return db, nil
	}
}

func main() {
	DB_Connect()
}
