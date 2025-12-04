package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jshk00/auto-pstate/internal"
)

var banner = `----------------------------------------------------
pdctl is client library to interact with pstated
----------------------------------------------------
status        - fetches the latest status of pstated.
set [profile] - sets the give profile this will automatically turns on manual mode.
auto [on/off] - set if auto mode should be on or off.
list-perfs    - list the available energy preferences.
`

var client = &http.Client{
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", internal.SockPath)
		},
	},
}

func main() {
	log.SetFlags(0)
	switch arg(1) {
	case "list-prefs":
		listPrefs()
	case "status":
		status()
	case "set":
		set()
	case "auto":
		auto()
	default:
		fmt.Print(banner)
	}
}

func status() {
	res, err := client.Post("http://localhost:9010/status", "", nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	var info map[string]string
	if err = json.NewDecoder(res.Body).Decode(&info); err != nil {
		log.Println(err)
		return
	}
	for k, v := range info {
		log.Printf("%s --> %s\n", k, v)
	}
}

func listPrefs() {
	prefs, err := internal.GetPreferences()
	if err != nil {
		log.Println(err)
		return
	}
	for _, p := range prefs {
		fmt.Println(p)
	}
}

func set() {
	profile := arg(2)
	if profile == "" {
		log.Fatal("set requires addtional argument profile") //nolint
	}
	res, err := client.Post("http://localhost:9010/set/"+profile, "", nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusBadRequest {
		log.Printf("%s profile is invalid\n", profile)
		return
	}
	log.Printf("p-state set to %s(manual mode)\n", profile)
}

func auto() {
	on := arg(2)
	if on == "on" || on == "off" {
		res, err := client.Post("http://localhost:9010/auto/"+on, "", nil)
		if err != nil {
			log.Println(err)
		}
		defer res.Body.Close()
		if res.StatusCode == http.StatusOK {
			if on == "on" {
				log.Println("auto mode enabled")
				return
			}
			log.Println("auto mode disabled")
		}
		return
	}
	fmt.Println("invalid option should be on/off")
}

func arg(i int) string {
	if len(os.Args) > i {
		return os.Args[i]
	}
	return ""
}
