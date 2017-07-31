package handlers

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"errors"
	"fmt"
)



func jsonValidate(js *MsgJson) (err error) {

	if ok := regMatcher.MatchString(js.AccessToken); !ok {
		return errors.New("invalid AccessToken")
	}

	eventCodeLen := len(js.EventCode)
	if eventCodeLen < minEventCodeLen || eventCodeLen > maxEventCodeLen {
		return errors.New("invalid EventCode len")
	}

	if _, ok := streamTypeMap[js.StreamType]; !ok {
		return errors.New("invalid stream type")
	}

	dataStreamKey := "person_" + js.StreamType
	if _, ok := js.Data[dataStreamKey]; !ok {
		return errors.New("stream typed data not exist")
	}
	return nil
}

func Loop() {

	cfg.logger.Print("init done:")


	foo := func() {
		for d := range cfg.msgs {
			err := jobHandling(&d)
			if err != nil {
				cfg.logger.Print(err)
				continue
			}
		}
	}

	for i := 0; i < cfg.GoroutinesNum - 1; i++ {
		go foo()
	}

	foo()

	
}


func jobHandling(d *amqp.Delivery) (err error) {
	defer func () {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint("handling error: ", r))
			return
		}
		return

	}()

	jsOut, err := msgHandle(d)
	if err != nil {
		return err
	}

	outData, err := json.Marshal(jsOut)
	if err != nil {
		return err
	}

	err = InsertMessage(jsOut, outData)

	if err != nil {
		return err
	}

	return nil

}

func msgHandle(d *amqp.Delivery) (jsOut *OutMsgJson, err error) {


	var js MsgJson
	err = json.Unmarshal(d.Body, &js)
	if err != nil {
		return nil, err
	}

	err = jsonValidate(&js)
	if err != nil {
		return nil, err
	}

	personStreamType := "person_" + js.StreamType

	personStreamVal := js.Data[personStreamType].(string)
	delete(js.Data, personStreamType)

	jsOut = &OutMsgJson{MsgJson: js, To: personStreamVal }

	return jsOut, nil

}