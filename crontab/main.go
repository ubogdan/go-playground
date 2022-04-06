package main

import (
	"log"

	"github.com/mileusna/crontab"
)

func main() {
	ctab := crontab.New() // create cron table

	err := ctab.AddJob("10 3 * * *", dailyCron) // daily at 3:10
	if err != nil {
		log.Printf("crontab.AddJob failed: %s", err)
		return
	}

	select {}
}

func dailyCron() {
	log.Printf("daily job is running ....")
}
