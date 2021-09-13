package main

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// NumericSegmentCharOffset is used to decrement the ascii byte values of each
// character in a JWS so that they can be encoded as two-digit numbers.
const NumericSegmentCharOffset = 45

// findNumericSegment parses the numeric segment from a QR code payload.
func findNumericSegment(qrPayload string) (string, error) {
	re, err := regexp.Compile("shc:/([1-9]/[1-9]/|)([0-9]+)")
	if err != nil {
		return "", err
	}
	matches := re.FindStringSubmatch(qrPayload)
	if len(matches) != 3 {
		return "", fmt.Errorf("failed to find numerically-encoded JWS in given data")
	}
	if matches[1] != "" {
		return "", fmt.Errorf("cannot decode chunked JWS")
	}
	return matches[2], nil
}

// decodeNumericSegment decodes the numerically encoded segment of the QR payload.
func decodeNumericSegment(numericSegment string) ([]byte, error) {
	if len(numericSegment)%2 != 0 {
		return nil, fmt.Errorf("invalid numeric segment to decode; length must be an even number")
	}
	decodedSegment := make([]byte, len(numericSegment)/2)
	for i := 0; i*2 < len(numericSegment); i++ {
		intStr := numericSegment[i*2 : i*2+2]
		parsed, _ := strconv.Atoi(intStr)
		decodedSegment[i] = byte(parsed + NumericSegmentCharOffset)
	}
	return decodedSegment, nil
}

// decodeJWS decodes a base64-URL-encoded JSON Web Signature
// and returns its header, payload, and signature data.
func decodeJWS(jwsToken string) (header, payload, signature []byte, err error) {
	segments := strings.Split(jwsToken, ".")
	if len(segments) != 3 {
		err = fmt.Errorf("Expected 3 JWS segments; found %d", len(segments))
		return
	}

	decodedSegments := make([][]byte, len(segments))

	for i, segment := range segments {
		decodedSegment, err := base64.RawURLEncoding.DecodeString(segment)
		if err != nil {
			return nil, nil, nil, err
		}

		decodedSegments[i] = decodedSegment
	}

	header = decodedSegments[0]
	payload = decodedSegments[1]
	signature = decodedSegments[2]
	return
}

// inflate reverses the deflate compression algorithm on a given slice of binary data.
func inflate(compressed []byte) ([]byte, error) {
	inflateReader := flate.NewReader(bytes.NewReader(compressed))
	uncompressed, err := io.ReadAll(inflateReader)
	if err != nil {
		return nil, err
	}
	return uncompressed, nil
}

// DecodeSmartVaccineCard decodes the given qrPayload (meant to come directly from a scanned QR code)
// and returns the JSON strings it encodes, and the signature it carries.
func DecodeSmartVaccineCard(qrPayload string) (headerJSON, cardJSON string, signature []byte, err error) {
	numericSegment, err := findNumericSegment(qrPayload)
	if err != nil {
		return "", "", nil, err
	}

	decodedNumericSegment, err := decodeNumericSegment(numericSegment)
	if err != nil {
		return "", "", nil, err
	}

	header, compressedPayload, signature, err := decodeJWS(string(decodedNumericSegment))
	if err != nil {
		return "", "", nil, err
	}

	payload, err := inflate(compressedPayload)
	if err != nil {
		return "", "", nil, err
	}

	headerJSON = string(header)
	cardJSON = string(payload)

	return
}
