package utils

import "fmt"

var Status map[uint8]string = map[uint8]string{
	0: "Need Mechanic Verification",
	1: "Need Supervisor Verification",
	2: "Need Superintendent Verification",
	3: "Need Ticket",
	4: "Ticket Created",
	5: "Completed",
}

type statusCodeError struct {
	code uint8
}

func (e statusCodeError) Error() string {
	return fmt.Sprintf("Invalid status code: %d", e.code)
}

func GetStatus(statusCode uint8) (string, error) {
	status, exists := Status[statusCode]

	if !exists {
		return "", statusCodeError{code: statusCode}
	}

	return status, nil
}
