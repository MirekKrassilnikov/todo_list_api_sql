package config

import "database/sql"

// Порт, на котором будет работать сервер
const Port = "7540"

// Директория для сервирования файлов
const WebDir = "./web"
const Layout = "20060102"

var db *sql.DB
