package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/danielpieper/alexa/internal/config"
	"github.com/danielpieper/alexa/internal/service"
	"github.com/danielpieper/lgtv-go"
	"github.com/mdlayher/wol"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	tvClient, err := lgtv.New("http://" + cfg.TVEndpoint)
	if err != nil {
		panic(err)
	}
	if err = tvClient.Authenticate(cfg.TVAuth); err != nil {
		panic(err)
	}

	wolClient, err := wol.NewClient()
	if err != nil {
		panic(err)
	}

	svc, err := service.New(
		wolClient,
		&tvClient,
		cfg.TVEndpoint,
		cfg.TVAuth,
		cfg.WOLBroadcast,
		cfg.WOLMac,
	)
	if err != nil {
		panic(err)
	}

	// if svc.IsPoweredOn() {
	// 	log.Println("powerered on")
	// }
	// os.Exit(0)

	// if err = tvClient.KeyExternalInput(); err != nil {
	// 	panic(err)
	// }
	// os.Exit(0)

	http.HandleFunc("/ps5", func(w http.ResponseWriter, _ *http.Request) {
		log.Println("switch to steam deck")
		svc.SwitchToPS5()

		w.WriteHeader(http.StatusNoContent)
	})
	http.HandleFunc("/firetv", func(w http.ResponseWriter, _ *http.Request) {
		log.Println("switch to Fire TV")
		svc.SwitchToFireTV()

		w.WriteHeader(http.StatusNoContent)
	})
	http.HandleFunc("/tv", func(w http.ResponseWriter, _ *http.Request) {
		log.Println("switch to TV")
		svc.SwitchToTV()

		w.WriteHeader(http.StatusNoContent)
	})
	// http.HandleFunc("/ps5", func(w http.ResponseWriter, _ *http.Request) {
	// 	log.Println("switch to PS5")
	// 	switchToPS5(client)
	//
	// 	w.WriteHeader(http.StatusNoContent)
	// })

	err = http.ListenAndServe(":"+cfg.HttpPort, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
