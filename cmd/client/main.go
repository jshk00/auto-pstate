package main

import (
	"context"
	"fmt"
	"io"
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
	switch arg(1) {
	case "list-prefs":
		prefs, err := internal.GetPreferences()
		if err != nil {
			log.Fatal(err)
		}
		for _, p := range prefs {
			fmt.Println(p)
		}
	case "status":
		res, err := client.Post("http://localhost:9010/status", "", nil)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		fmt.Println(string(b))
	case "set":
		profile := arg(2)
		if profile == "" {
			log.Fatal("set require addtional argument profile") //nolint
		}
		res, err := client.Post("http://localhost:9010/set/"+profile, "", nil)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode == http.StatusBadRequest {
			log.Printf("%s profile is invalid\n", profile)
			return
		}
		log.Printf("successfully set profile %s and manual mode has been enabled\n", profile)
	case "auto":
		on := arg(2)
		if on == "on" || on == "off" {
			res, err := client.Post("http://localhost:9010/auto/"+on, "", nil)
			if err != nil {
				log.Println(err)
			}
			defer res.Body.Close()
			if res.StatusCode == http.StatusOK {
				fmt.Println("pstated auto mode has been set to", on)
			}
			return
		}
		fmt.Println("invalid option should be on/off")
	default:
		fmt.Print(banner)
	}
}

func arg(i int) string {
	if len(os.Args) > i {
		return os.Args[i]
	}
	return ""
}
