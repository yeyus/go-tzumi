package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yeyus/go-tzumi/tzumi/responses"

	"github.com/yeyus/go-tzumi/tzumi"
)

var atscFreqs = [...]int{ /*57000, 63000, 69000, 79000, 85000, 177000, 183000, 189000, 195000, 201000, 207000, 213000, 473000, 479000, 485000, 491000, 497000, 503000, 509000, 515000, 521000, 527000, 533000, 539000, 545000, 551000, 557000, 563000, 569000, 575000, 581000, 587000, 593000, 599000, 605000, 611000,*/ 617000, 623000, 629000, 635000, 641000, 647000, 653000, 659000, 665000, 671000, 677000, 683000, 689000, 695000, 701000, 707000, 713000, 719000, 725000, 731000, 737000, 743000, 749000, 755000, 761000, 767000, 773000, 779000, 785000, 791000, 797000, 803000}

// 177000
// 617000

func main() {
	t, err := tzumi.NewTzumiMagicTV("192.168.1.1")
	defer t.Close()
	if err != nil {
		panic(err)
	}

	t.Debug = true

	// Instantiate close channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		t.Close()
		os.Exit(1)
	}()

	t.Login(func(r responses.Response) {
		log.Printf("[Main] First connection established")

		freqIndex := 0

		for freqIndex < len(atscFreqs) {
			log.Printf("[Main] Start of tuning loop for index %d => %d", freqIndex, atscFreqs[freqIndex])
			t.Tune(atscFreqs[freqIndex], 0, func(r responses.Response) {
				if t.State == tzumi.TUNED {
					log.Printf("[Main] Tuning succeed @ %d !!", atscFreqs[freqIndex])
					//t.CheckLock()
					//t.GetSignalStatus()
				} else {
					log.Printf("[Main] Tuning failed!")
				}
			})
			freqIndex++
		}
	})

}
