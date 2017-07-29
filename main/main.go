package main

import (
	"queue/handlers"
	//_ "queue/db"
)

func main()  {

	cfg := handlers.NewConfig()
	cfg.AmqpUrl = "amqp://guest:guest@rabbit:5672/"
	cfg.DbUrl= "postgres://postgres:postgres@queue-db:5432/queue?sslmode=disable"
	cfg.Init()


	handlers.Loop()
}
