package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli"
	"github.com/yeyus/go-tzumi/tzumi/responses"

	"github.com/yeyus/go-tzumi/tzumi"
)

var atscFreqs = [...]int{57000, 63000, 69000, 79000, 85000, 177000, 183000, 189000, 195000, 201000, 207000, 213000, 473000, 479000, 485000, 491000, 497000, 503000, 509000, 515000, 521000, 527000, 533000, 539000, 545000, 551000, 557000, 563000, 569000, 575000, 581000, 587000, 593000, 599000, 605000, 611000, 617000, 623000, 629000, 635000, 641000, 647000, 653000, 659000, 665000, 671000, 677000, 683000, 689000, 695000, 701000, 707000, 713000, 719000, 725000, 731000, 737000, 743000, 749000, 755000, 761000, 767000, 773000, 779000, 785000, 791000, 797000, 803000}

// Working frequencies in SF
// 2018/12/04 23:00:14 [Find] Results:
// 2018/12/04 23:00:14      177000Khz
// 2018/12/04 23:00:14      207000Khz
// 2018/12/04 23:00:14      557000Khz
// 2018/12/04 23:00:14      563000Khz
// 2018/12/04 23:00:14      569000Khz
// 2018/12/04 23:00:14      593000Khz
// 2018/12/04 23:00:14      605000Khz
// 2018/12/04 23:00:14      617000Khz
// 2018/12/04 23:00:14      623000Khz
// 2018/12/04 23:00:14      629000Khz
// 2018/12/04 23:00:14      635000Khz
// 2018/12/04 23:00:14      647000Khz
// 2018/12/04 23:00:14      653000Khz
// 2018/12/04 23:00:14      659000Khz

var host string
var debug bool
var tuner *tzumi.TzumiMagicTV

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "increased log verbosity",
			Destination: &debug,
		},
		cli.StringFlag{
			Name:        "host",
			Usage:       "tuner ip",
			Destination: &host,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "find",
			Usage: "find the list of active frequencies",
			Action: func(c *cli.Context) (err error) {
				tuner, err = tzumi.NewTzumiMagicTV(host)
				if err != nil {
					tuner.Close()
					return err
				}

				tuner.Debug = debug

				tuner.Login(func(r responses.Response) {
					log.Printf("[Main] First connection established")
					foundList := findChannels(tuner, atscFreqs[:])
					log.Printf("[Find] Results:")
					for i := 0; i < len(foundList); i++ {
						log.Printf("\t %dKhz", foundList[i])
					}
				})
				return nil
			},
		},
		{
			Name:  "tune",
			Usage: "tune to frequency and start TS proxy",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "frequency",
					Usage: "tuner frequency in Khz",
				},
			},
			Action: func(c *cli.Context) (err error) {
				tuner, err = tzumi.NewTzumiMagicTV(host)
				if err != nil {
					tuner.Close()
					return err
				}

				tuner.Debug = debug
				frequency := c.Int("frequency")

				tuner.Login(func(r responses.Response) {
					tuner.Tune(frequency, 0, true, func(r2 responses.Response) {
						log.Printf("Tuned to %dKhz and program 0", frequency)
					})
				})
				return nil
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		if len(host) == 0 {
			return errors.New("Host is a required parameter")
		}

		// Instantiate close channel
		closechannel := make(chan os.Signal, 1)
		signal.Notify(closechannel, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-closechannel
			tuner.Close()
			os.Exit(1)
		}()

		return nil
	}

	app.After = func(c *cli.Context) error {
		if tuner != nil && tuner.State != tzumi.DISCONNECTED {
			tuner.Close()
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func findChannels(tuner *tzumi.TzumiMagicTV, freqList []int) []int {
	freqIndex := 0

	foundBitmap := make([]bool, len(freqList))
	found := 0

	for freqIndex < len(freqList) {
		log.Printf("[findChannels] Start of tuning loop for index %d => %d", freqIndex, freqList[freqIndex])
		tuner.Tune(freqList[freqIndex], 0, false, func(r responses.Response) {
			if tuner.State == tzumi.TUNED {
				log.Printf("[findChannels] Tuning succeed @ %d !!", freqList[freqIndex])
				foundBitmap[freqIndex] = true
				found++
				//t.CheckLock()
				//t.GetSignalStatus()
			} else {
				log.Printf("[findChannels] Tuning failed!")
				foundBitmap[freqIndex] = false
			}
		})
		freqIndex++
	}

	foundList := make([]int, found)
	var j int = 0
	for i := 0; i < len(foundBitmap); i++ {
		if foundBitmap[i] {
			foundList[j] = freqList[i]
			j++
		}
	}

	return foundList
}
