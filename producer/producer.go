package main

import (
	"github.com/streadway/amqp"
	"encoding/json"
	"log"
	"fmt"
)

type jsStruct struct {
	AccessToken string `json:"access_token"`
	EventCode string `json:"event_code"`
	StreamType string `json:"stream_type"`
	Data map[string]interface{} `json:"data"`

}


type jsOutStruct struct {
	jsStruct
	To string `json:"to"`
}


func main() {
	//"access_token": "0d10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
	//"event_code": "ispp",
	outString := `{
  "access_token": "0d10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
  "event_code": "ispp",
  "stream_type": "email",
  "data": {
    "person_name": "Иван",
    "date": "2016-03-03",
    "person_email": "ivanivanov@gmail.com"
  }
}
`



	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	defer conn.Close()

	if err != nil {
		log.Print(err)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Print(err)
		return
	}

	defer ch.Close()

	q, err := ch.QueueDeclare("queueTest",false,false,false,false,nil)

	if err != nil {
		log.Print(err)
		return
	}


	jsonObject := jsStruct{}
	err = json.Unmarshal([]byte(outString), &jsonObject)
	if err != nil {
		fmt.Print(err)
		return
	}



	outData, err := json.Marshal(jsonObject)
	if err != nil {
		log.Print(err)
		return
	}

	ch.Publish("",q.Name,false,false, amqp.Publishing{
		ContentType: "application/json",
		Body: outData,
	})

	log.Print(string(outData))

}
