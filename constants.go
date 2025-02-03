package main

import "os"

var (
	DbUser = os.Getenv("DB_USER")
	DbPass = os.Getenv("DB_PASSWORD")
	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbName = os.Getenv("DB_NAME")
)
