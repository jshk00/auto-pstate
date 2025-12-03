package main

import (
	"fmt"
	"log"
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

	case "set":
	case "auto":
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
