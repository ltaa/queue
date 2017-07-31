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



func pushRabbitRequestGoroutine(jsonString string, ch *amqp.Channel ) {
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

	ch.Publish("", cfg.AmqpChannel,false,false, amqp.Publishing{
		ContentType: "application/json",
		Body: outData,
	})

}



func BenchmarkValidTokenGoroutines(b *testing.B) {
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


	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	defer conn.Close()

	if err != nil {
		return
	}

	ch, err := conn.Channel()

	if err != nil {
		return
	}

	defer ch.Close()

	for n := 0; n < 10000; n++ {
		for _, data := range testData {
			pushRabbitRequestGoroutine(data.src, ch)
		}
	}

	cfg.amqpChan.Close()
	//count := runtime.NumCPU()


	tmpChan := make (chan *amqp.Delivery)

	f := func(msgs <- chan *amqp.Delivery) {
		for d := range msgs {
			err := jobHandling(d)
			if err != nil {
				cfg.logger.Print(err)
				continue
			}
		}
	}

	for i := 0; i < 5; i++  {
		go f(tmpChan)
	}


	for n := 0; n < b.N; n++ {
		for d := range cfg.msgs {
			//jobHandling(&d)
			//tVal:= <-d

			tmpChan <- &d
		}

	}

}



func BenchmarkValidToken(b *testing.B) {
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


	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	defer conn.Close()

	if err != nil {
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		return
	}

	defer ch.Close()

	for n := 0; n < 10000; n++ {
		for _, data := range testData {
			pushRabbitRequestGoroutine(data.src, ch)
		}
	}

	cfg.amqpChan.Close()

	for n := 0; n < b.N; n++ {
		for d := range cfg.msgs {
			jobHandling(&d)

		}

	}

}