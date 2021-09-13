package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

func run() error {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "%s: Parse SMART Health Cards' QR code payloads.\n", os.Args[0])
		fmt.Fprintf(out, "SMART Health Cards: https://smarthealth.cards/.\n\n")
		fmt.Fprintf(out, "USAGE:\n")
		fmt.Fprintf(out, "  %s <QR_PAYLOAD>\n\n", os.Args[0])
		fmt.Fprintf(out, "Decodes the JSON Web Signature embedded in the QR code payload\n\n")
		fmt.Fprintf(out, "\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return fmt.Errorf("Must pass exactly one argument: <QR_PAYLOAD>")
	}

	qrPayload := args[0]

	header, card, signature, err := DecodeSmartVaccineCard(qrPayload)
	if err != nil {
		return err
	}

	fmt.Println("Decoded JSON Web Signature vaccine card data")
	fmt.Println()
	fmt.Println("Header:")
	fmt.Println("---------------------------------------------")
	fmt.Println(string(header))
	fmt.Println()
	fmt.Println("Payload:")
	fmt.Println("---------------------------------------------")
	fmt.Println(string(card))
	fmt.Println()
	fmt.Println("Signature:")
	fmt.Println("---------------------------------------------")
	fmt.Println(hex.EncodeToString(signature))

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
