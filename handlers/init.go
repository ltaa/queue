package handlers

import (
	"regexp"
	"log"
)

var regMatcher *regexp.Regexp
var streamTypeMap map[string]bool
func init() {

	var err error
	regMatcher, err = regexp.Compile(regexpString)

	if err != nil {
		log.Print(err)
		panic(err)
		return
	}
}

func init() {
	streamTypeMap = make(map[string]bool)
	streamTypeMap["sms"] = true
	streamTypeMap["push"] = true
	streamTypeMap["email"] = true

}