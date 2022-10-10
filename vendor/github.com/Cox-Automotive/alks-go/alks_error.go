package alks

import (
	"fmt"
)

type AlksError struct {
	StatusCode int
	RequestId  string `json:"requestId"`
	Err        error
}

func (r *AlksError) Error() string {
	return fmt.Sprintf("status %d: requestID %s: err %v", r.StatusCode, r.RequestId, r.Err)
}

type AlksResponseError struct {
	StatusMessage string   `json:"statusMessage"`
	Errors        []string `json:"errors"`
	RequestId     string   `json:"requestId"`
}

var ErrorStringFull = "[%s] ALKS Error %d Msg: %s\n Contact the ALKS Team for assistance on Slack at #alks-client-support"
var ErrorStringNoReqId = "ALKS Error %d Msg: %s\n Contact the ALKS Team for assistance on Slack at #alks-client-support"
var ErrorStringOnlyCodeAndReqId = "[%s] ALKS Error %d\n Contact the ALKS Team for assistance on Slack at #alks-client-support"
var ErrorStringOnlyCode = "ALKS Error %d\n Contact the ALKS Team for assistance on Slack at #alks-client-support"
var ParseErrorReqId = "[%s] Error parsing ALKS Error response: %s"
var ParseError = "Error parsing ALKS Error response: %s"
