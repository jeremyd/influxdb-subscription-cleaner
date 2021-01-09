package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/spf13/viper"
)

// This handy script will DELETE ALL INFLUXDB SUBSCRIPTIONS!
// Workaround for https://github.com/influxdata/kapacitor/issues/870
// Configure using Environment variables:
// INFLUXDB_URL=http://myinflux:8086
// INFLUXDB_DRYRUN=true // for dry run

func checkIfSet(name string) {
	if !viper.IsSet(name) {
		fmt.Printf("You must set the environment variable $INFLUXDB_%s, exiting..\n", name)
		os.Exit(1)
	}
}

func clean() {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: viper.GetString("url"),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	q := client.NewQuery("SHOW SUBSCRIPTIONS", "", "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, sub := range response.Results[0].Series {
			fmt.Printf("Series Name: %s\n", sub.Name)
			for _, v := range sub.Values {
				retention := v[0]
				dbname := sub.Name
				name := v[1]
				dropQuery := fmt.Sprintf("drop subscription \"%s\" on \"%s\".\"%s\"", name, dbname, retention)
				if viper.IsSet("dryrun") {
					fmt.Printf("we would run: %s\n", dropQuery)
				} else {
					dq := client.NewQuery(dropQuery, "", "")
					fmt.Printf("running: %s\n", dropQuery)
					if deleteResponse, err := c.Query(dq); err == nil && response.Error() == nil {
						fmt.Println(deleteResponse)
					} else {
						fmt.Printf("error while deleting: %s", err)
					}
				}
			}
		}
	}
}

func main() {
	viper.SetEnvPrefix("influxdb")
	viper.AutomaticEnv()
	checkIfSet("url")
	// If interval is set, re-run every interval hours.
	if viper.IsSet("interval") {
		fmt.Printf("Starting up and cleaning.")
		clean()
		fmt.Printf("Cleaning subscriptions every %d seconds\n", viper.GetInt("interval"))
		ticker := time.NewTicker(time.Second * viper.GetDuration("interval"))
		go func() {
			for t := range ticker.C {
				fmt.Println("Starting cleaning run ", t)
				clean()
			}
		}()

		// wait for exit
		exitSignal := make(chan os.Signal)
		signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
		<-exitSignal

	} else {
		clean()
	}
}
