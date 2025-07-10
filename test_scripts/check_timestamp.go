package main

import (
	"fmt"
	"time"
)

func main() {
	// Check what Unix timestamp 1720558320 converts to
	t := time.Unix(1720558320, 0).UTC()
	fmt.Printf("Unix 1720558320 = %s\n", t.String())
	fmt.Printf("RFC3339: %s\n", t.Format(time.RFC3339))
	
	// Check what 2024-07-09T22:12:00Z converts to
	target, _ := time.Parse(time.RFC3339, "2024-07-09T22:12:00Z")
	fmt.Printf("Target time: %s\n", target.String())
	fmt.Printf("Target Unix: %d\n", target.Unix())
	
	// Are they equal?
	fmt.Printf("Equal: %v\n", t.Equal(target))
}