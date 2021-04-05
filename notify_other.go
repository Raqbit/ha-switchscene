// +build !linux

package main

import "fmt"

func sendNotification(body string, isError bool) error {
	fmt.Println(body)
	return nil
}
