package handlers


const (
	regexpString = `^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-(8|9|a|b)[a-f0-9]{3}-[a-f0-9]{12}$`
	minEventCodeLen = 1
	maxEventCodeLen = 255
)



type MsgJson struct {
	AccessToken string `json:"access_token"`
	EventCode string `json:"event_code"`
	StreamType string `json:"stream_type"`
	Data map[string]interface{} `json:"data"`

}


type OutMsgJson struct {
	MsgJson
	To string `json:"to"`
}