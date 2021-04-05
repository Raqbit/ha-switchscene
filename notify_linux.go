package main

import "os/exec"

func sendNotification(body string, isError bool) error {
	icon := "emblem-default"

	if isError {
		icon = "emblem-important"
	}

	return exec.Command("notify-send", "-i", icon, Title, body, "-h", NotifHintTransient).Run()
}
