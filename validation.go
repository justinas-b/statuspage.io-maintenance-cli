package main

import (
	"fmt"
	"strconv"
	"time"
)

func validateEmptyString(s string) error {
	if s == "" {
		return fmt.Errorf("title cannot be empty")
	}
	return nil
}

func validateDateTimeFormat(s string) error {
	startTime, err := time.Parse(time.DateTime, s)
	if err != nil {
		return err
	}

	inFuture := startTime.After(time.Now().UTC())
	if !inFuture {
		return fmt.Errorf("provided time is not in the future")
	}

	return nil
}

func validateInteger(s string) error {
	_, err := strconv.Atoi(s)
	return err
}
