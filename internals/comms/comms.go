package comms

import "errors"

var errNoAction = errors.New("error : the action in the request json doesn't exist or is not specified")

type Request struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}
