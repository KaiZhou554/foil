package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

// Verify that a WebP byte array has the ALPHA flag set.
func main() {
	// Decode a test base64 WebP (a minimal transparent pixel)
	// This is a 1x1 transparent WebP: RIFF + WEBP + VP8X (alpha)
	testB64 := "UklGRiQAAABXRUJQVlA4IBgAAAAwAQCdASoBAAEAAQAcJaQAA3gA/v+QAA=="
	data, err := base64.StdEncoding.DecodeString(testB64)
	if err != nil {
		fmt.Printf("base64 decode failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Decoded %d bytes\n", len(data))
	
	// Check RIFF header
	if len(data) < 12 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WEBP" {
		fmt.Println("Not a valid WebP file")
		os.Exit(1)
	}
	fmt.Println("Valid WebP file ✓")
	
	// Check for ALPH chunk (alpha channel)
	pos := 12
	for pos+8 <= len(data) {
		chunkID := string(data[pos : pos+4])
		chunkSize := int(data[pos+4]) | int(data[pos+5])<<8 | int(data[pos+6])<<16 | int(data[pos+7])<<24
		if chunkID == "ALPH" {
			fmt.Printf("ALPHA chunk found at offset %d, size %d ✓\n", pos, chunkSize)
		}
		if chunkID == "VP8X" {
			fmt.Printf("VP8X chunk found (extended features) at offset %d\n", pos)
			if len(data) > pos+8+4 {
				flags := data[pos+8+4]
				if flags&0x10 != 0 {
					fmt.Println("  ALPHA flag SET in VP8X ✓")
				} else {
					fmt.Println("  ALPHA flag NOT set ✗")
				}
			}
		}
		pos += 8 + chunkSize
		if chunkSize%2 == 1 {
			pos++ // padding byte
		}
	}
	
	fmt.Println("\nScan complete")
}
