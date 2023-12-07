package main

import (
	"encoding/json"
	"flag"
	"github.com/achetronic/tapogo/api/types"
	"github.com/achetronic/tapogo/pkg/tapogo"
	"log"
	"os"
)

const (
	EmailFlagDescription    = "(Required) Email to access the device"
	PasswordFlagDescription = "(Required) Password to access the device"
	IpFlagDescription       = "(Required) Device IP address"
	CommandFlagDescription  = "Command to execute: device-info (default), on, off, energy-usage"
)

func main() {
	var tapoClient *tapogo.Tapo
	var response *types.ResponseSpec
	var err error

	// 1. Parse the flags from the command line
	emailFlag := flag.String("email", "", EmailFlagDescription)
	passwordFlag := flag.String("password", "", PasswordFlagDescription)
	ipFlag := flag.String("ip", "", IpFlagDescription)
	commandFlag := flag.String("command", "device-info", CommandFlagDescription)

	flag.Parse()

	if *emailFlag == "" || *passwordFlag == "" || *ipFlag == "" {
		log.Print("Some required flags are missing. execute 'tapogo --help' to help yourself")
		os.Exit(0)
	}

	// 2. Init the client: handshake
	tapoClient, err = tapogo.NewTapo(*ipFlag, *emailFlag, *passwordFlag)
	if err != nil {
		log.Fatalln(err)
	}

	// 3. Execute the action
	switch *commandFlag {
	case "off":
		response, err = tapoClient.TurnOff()
	case "on":
		response, err = tapoClient.TurnOn()
	case "energy-usage":
		response, err = tapoClient.GetEnergyUsage()
	case "device-info":
		response, err = tapoClient.DeviceInfo()
	default:
		log.Fatalf("invalid argument: '%s'. expected one of [off, on, energy-usage, device-info]", *commandFlag)
	}

	if err != nil {
		log.Fatalln(err)
	}

	jsonBytes, err := json.Marshal(*response)
	if err != nil {
		log.Print(err)
	}

	log.Printf("Response (processed): %s", string(jsonBytes))
}
