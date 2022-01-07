package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const VERSION  = `0.2`

var (
	DebugLog = log.New(os.Stderr, `debug#`, log.Lshortfile)
	ErrorLog = log.New(os.Stderr, `error#`, log.Lshortfile)
)

func helpText() {
	fmt.Print("This tool encrypts data with OTP protection\n\n")
	fmt.Print("USAGE:\n")
	fmt.Printf(" %s -e\t# To encrypt STDIN to STDOUT\n", os.Args[0])
	fmt.Printf(" %s -d\t# To decrypt STDIN to STDOUT\n", os.Args[0])
	fmt.Print(" You must specify -e, -d or -c\n\n")
	flag.PrintDefaults()
}


func main() {
	help := flag.Bool("h", false, "print this help")
	ver := flag.Bool("v", false, "Show version")
	encryptFlag := flag.Bool("e", false, "Encrypt mode")
	decryptFlag := flag.Bool("d", false, "Decrypt mode")
	askKeyCmdFlag := flag.String("ask-key", "", "Get key command")
	askOtpCmdFlag := flag.String("ask-otp", "", "Get OTP command")
	qrFilePathFlag := flag.String("qr", "", "Write QR-code to PNG file (used with -c)")
	inputFlag := flag.String("i", "", "Input file")
	outputFlag := flag.String("o", "", "Output file")
	noPromptFlag := flag.Bool("no-prompt", false, "Disable console prompt")
	createOtpAuthUrlFlag := flag.Bool(`c`, false, `Create OTP auth URL mode`)
	flag.Parse()

	if *help {
		helpText()
		os.Exit(0)
	}

	if *ver {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if *createOtpAuthUrlFlag {
		if err := CreateOtpAuthUrl(*qrFilePathFlag); err != nil { os.Exit(1) }
		os.Exit(0)
	}

	if *encryptFlag == *decryptFlag {
		helpText()
		os.Exit(1)
	}

	denyConsoleInteractive := false
	if *noPromptFlag || *inputFlag == `` || *outputFlag == `` { denyConsoleInteractive = true }

	if *encryptFlag {
		if err := Encrypt(*inputFlag, *outputFlag, *askKeyCmdFlag, denyConsoleInteractive); err != nil { os.Exit(1) }
		os.Exit(0)
	}

	if *decryptFlag {
		if err := Decrypt(*inputFlag, *outputFlag, *askKeyCmdFlag, *askOtpCmdFlag, denyConsoleInteractive); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
}
