package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/ssh/terminal"
	"net/http"
	"os"
	"syscall"
)

const (
	Title              = "Home Assistant Scene Switcher"
	NotifHintTransient = "int:transient:1"
	KeyringAppName     = "ha-switchscene"
)

var (
	haUrl      *string
	scene      *string
	name       *string
	storeToken *bool
)

func init() {
	storeToken = flag.Bool("storeToken", false, "store home assistant token")
	haUrl = flag.String("url", "", "home assistant url")
	scene = flag.String("scene", "", "scene to switch to")
	name = flag.String("name", "scene", "name of scene to switch to")
	flag.Parse()
}

func main() {
	if *haUrl == "" {
		_ = sendNotification("No Home Assistant URL given.", true)
		os.Exit(1)
	}

	if *storeToken {
		if err := readAndStoreToken(); err != nil {
			fmt.Printf("Could not store token: %s", err)
		}
		return
	}

	token, err := keyring.Get(KeyringAppName, *haUrl)

	if err != nil {
		_ = sendNotification("Could not find credentials in keyring.", true)
		fmt.Println(err)
		os.Exit(1)
	}

	if *scene == "" {
		_ = sendNotification("No scene given.", true)
		os.Exit(1)
	}

	if err = switchScene(token, *scene); err != nil {
		_ = sendNotification("Could not switch scene.", true)
		fmt.Println(err)
		os.Exit(1)
	}

	_ = sendNotification(fmt.Sprintf("Switched scene to %s.", *name), false)
}

func readAndStoreToken() error {
	fmt.Print("Enter Home Assistant token: ")
	bytepw, err := terminal.ReadPassword(syscall.Stdin)

	if err != nil {
		return err
	}

	if err = keyring.Set(KeyringAppName, *haUrl, string(bytepw)); err != nil {
		return err
	}

	fmt.Println("\nToken saved to keyring.")
	return nil
}

func switchScene(token, scene string) error {
	return callService(token, "scene", "turn_on",
		map[string]interface{}{
			"entity_id": scene,
		},
	)
}

func callService(token, domain, service string, data map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/services/%s/%s", *haUrl, domain, service)

	body := new(bytes.Buffer)

	if err := json.NewEncoder(body).Encode(data); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("home assistant error response: %d", res.StatusCode)
	}

	return nil
}
