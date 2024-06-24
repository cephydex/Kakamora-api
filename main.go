package main

import (
	"fmt"
	"net/http"
	"strconv"
	"svclookup/api"
	"svclookup/xutil"
	"time"

	"github.com/robfig/cron"
	zlog "github.com/rs/zerolog/log"
)

func main() {

	// create initial file
	api.CheckNCreate()
	fmt.Printf("\nSVC :: MAIN @ %s \n", time.Now().Format(time.RFC3339))

	loc := api.Accra()
	cj := cron.NewWithLocation(loc)

	// cronJob.AddFunc("*/20 * * * * *", func () {
	// cronJob.AddFunc("10 * * * *", func () {
	// cronJob.AddFunc("5 */1 * * *", func () {
	// c.AddFunc("@every 00h00m10s", GuestGreeting)
	cj.AddFunc("@every 00h05m00s", func () {
		now_time := time.Now()
		ttf := now_time.Format(time.RFC3339)
		fmt.Printf("\nINIT :: Task @ %s \n", ttf)
		
		api.AutoDiscover(now_time)
	})
	cj.Start()

	config, err := xutil.LoadConfig(".")
    if err != nil {
        zlog.Fatal().Msgf("cannot load config (main):", err)
    }
	strPort := config.AppPort
	serverPort, err := strconv.Atoi(strPort); if err != nil {
        // panic(err)
		zlog.Fatal().Msgf("string conversion:", err)
    }
	mux := http.NewServeMux()

	handleRequests(mux)

	server(mux, serverPort)
}
