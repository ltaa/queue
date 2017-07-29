package handlers

import (
	"encoding/json"
	"log"
	"github.com/streadway/amqp"
	"errors"
)



func jsonValidate(js *MsgJson) error {


	if ok := regMatcher.MatchString(js.AccessToken); !ok {
		return errors.New("invalid AccessToken")
	}

	eventCodeLen := len(js.EventCode)
	if eventCodeLen < minEventCodeLen || eventCodeLen > maxEventCodeLen {
		return errors.New("invalid EventCode len")
	}

	if _, ok := streamTypeMap[js.StreamType]; !ok {
		errors.New("invalid stream type")
	}

	dataStreamKey := "person_" + js.StreamType
	if _, ok := js.Data[dataStreamKey]; !ok {
		errors.New("stream typed data not exist")
	}
	return nil
}

func Handle() {


	con, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Print(err)
		return
	}
	defer con.Close()

	ch, err := con.Channel()

	if err != nil {
		log.Print(err)
		return
	}

	defer ch.Close()

	q, err := ch.QueueDeclare("queueTest",false,false, false, false, nil)

	if err != nil {
		log.Print(err)
		return
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	if err != nil {
		log.Print(err)
		return
	}

	waitchan := make(chan bool)

	go func() {
		for d := range msgs {

			log.Print("body: ", string(d.Body))
			var js MsgJson
			err := json.Unmarshal(d.Body, &js)
			if err != nil {
				log.Print(err)
				continue
				//return
			}

			err = jsonValidate(&js)
			if err != nil {
				log.Print(err)
				continue
			}

			email := js.Data["person_email"].(string)
			delete(js.Data, "person_email")

			jsOut := OutMsgJson{MsgJson: js, To: email }




			outData, err := json.Marshal(jsOut)
			if err != nil {
				log.Print(err)
				continue
				//return
			}


			err = InsertMessage(&jsOut, outData)

			if err != nil {
				log.Print(err)
				continue
			}

			//fmt.Print(string(outData))

		}
	} ()


	<- waitchan
}