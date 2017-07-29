package handlers

import (
	"fmt"
	"github.com/streadway/amqp"
	"encoding/json"
	"testing"
	"reflect"
)

type tokenTestData struct {
	src string
	outJs *OutMsgJson
	want bool
}

func TestValidToken(t *testing.T) {

	testData := [] tokenTestData {
		{
			src: `{
				"access_token": "0d10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
				"event_code": "ispp",
				"stream_type": "sms",
				"data": {
					"person_name": "Иван",
					"date": "2016-03-03",
					"person_sms": "ivanivanov@gmail.com"
				}
			}`,
			outJs: &OutMsgJson{
				MsgJson: MsgJson {
					AccessToken: "0d10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
					EventCode: "ispp",
					StreamType: "sms",
					Data: map[string]interface{} {
						"person_name": "Иван",
						"date": "2016-03-03",
					},
				},
				To: "ivanivanov@gmail.com",
			},
			want: true,

		},
		{
			src: `{
				"access_token": "0d10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
				"event_code": "ispp",
				"stream_type": "push",
				"data": {
					"person_name": "Иван",
					"date": "2016-03-03",
					"person_push": "ivanivanov@gmail.com"
				}
			}`,
			outJs: &OutMsgJson{
				MsgJson: MsgJson {
					AccessToken: "0d10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
					EventCode: "ispp",
					StreamType: "push",
					Data: map[string]interface{} {
						"person_name": "Иван",
						"date": "2016-03-03",
					},
				},
				To: "ivanivanov@gmail.com",
			},
			want: true,

		},
		{
			src: `{
				"access_token": "0d10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
				"event_code": "ispp",
				"stream_type": "email",
				"data": {
					"person_name": "Иван",
					"date": "2016-03-03",
					"person_email": "ivanivanov@gmail.com"
				}
			}`,
			outJs: &OutMsgJson{
				MsgJson: MsgJson {
					AccessToken: "0d10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
					EventCode: "ispp",
					StreamType: "email",
					Data: map[string]interface{} {
						"person_name": "Иван",
						"date": "2016-03-03",
					},
				},
				To: "ivanivanov@gmail.com",
			},
			want: true,

		},
		{
			src: `{
				"access_token": "0D10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
				"event_code": "ispp",
				"stream_type": "email",
				"data": {
					"person_name": "Иван",
					"date": "2016-03-03",
					"person_email": "ivanivanov@gmail.com"
				}
			}`,
			outJs: nil,
			want: true,

		},
		{
			src: `{
				"access_token": "0G10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
				"event_code": "ispp",
				"stream_type": "email",
				"data": {
					"person_name": "Иван",
					"date": "2016-03-03",
					"person_email": "ivanivanov@gmail.com"
				}
			}`,
			outJs: &OutMsgJson{
				MsgJson: MsgJson {
					AccessToken: "0G10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
					EventCode: "ispp",
					StreamType: "email",
					Data: map[string]interface{} {
						"person_name": "Иван",
						"date": "2016-03-03",
					},
				},
				To: "ivanivanov@gmail.com",
			},
			want: false,

		},
		{
			src: `{
				"access_token": "0d10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
				"event_code": "",
				"stream_type": "email",
				"data": {
					"person_name": "Иван",
					"date": "2016-03-03",
					"person_email": "ivanivanov@gmail.com"
				}
			}`,
			outJs: &OutMsgJson{
				MsgJson: MsgJson {
					AccessToken: "0d10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
					EventCode: "",
					StreamType: "email",
					Data: map[string]interface{} {
						"person_name": "Иван",
						"date": "2016-03-03",
					},
				},
				To: "ivanivanov@gmail.com",
			},
			want: false,

		},
		{
			src: `{
				"access_token": "0d10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
				"event_code": "",
				"stream_type": "email",
				"data": {
					"person_name": "Иван",
					"date": "2016-03-03",
					"person_email": "ivanivanov@gmail.com"
				}
			}`,
			outJs: nil,
			want: true,

		},
		{
			src: `{
				"access_token": "0d10566b-7e7f-4c17-b2ea-f0e42a4df3c0",
				"event_code": "",
				"stream_type": "qwerty",
				"data": {
					"person_name": "Иван",
					"date": "2016-03-03",
					"person_email": "ivanivanov@gmail.com"
				}
			}`,
			outJs: nil,
			want: true,

		},
	}

	conf := NewConfig()
	conf.Init()

	for _, data := range testData {
		pushRabbitRequest(data.src)
		d := <- cfg.msgs
		jsOut, _ := msgHandle(&d)

		if eq := reflect.DeepEqual(jsOut, data.outJs); eq != data.want {
			if data.want {
				t.Errorf("invalid transform:\n get %v\n want: %v", jsOut, data.outJs)
			} else {
				t.Errorf("invalid transform:\n get %v\n must not equal: %v", jsOut, data.outJs)
			}
		}




	}

}


func pushRabbitRequest(jsonString string) {
	jsonObject := MsgJson{}
	err := json.Unmarshal([]byte(jsonString), &jsonObject)
	if err != nil {
		fmt.Print(err)
		return
	}


	outData, err := json.Marshal(jsonObject)
	if err != nil {
		cfg.logger.Print(err)
		return
	}

	cfg.amqpChan.Publish("", cfg.AmqpChannel,false,false, amqp.Publishing{
		ContentType: "application/json",
		Body: outData,
	})

}