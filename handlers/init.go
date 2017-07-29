package handlers

import (
	"regexp"
)

var regMatcher *regexp.Regexp
var streamTypeMap map[string]bool
func init() {

	var err error
	regMatcher, err = regexp.Compile(regexpString)

	if err != nil {
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