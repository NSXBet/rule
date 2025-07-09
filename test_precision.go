package rule

import (
	"fmt"
	"strconv"
)

func main() {
	// Test the exact values that are failing
	str1 := "9223372036854775806"
	str2 := "9223372036854775807"

	f1, _ := strconv.ParseFloat(str1, 64)
	f2, _ := strconv.ParseFloat(str2, 64)

	fmt.Printf("String: %s -> Float64: %.0f\n", str1, f1)
	fmt.Printf("String: %s -> Float64: %.0f\n", str2, f2)
	fmt.Printf("Are they equal? %v\n", f1 == f2)

	// Test int64 values
	i1 := int64(9223372036854775806)
	i2 := int64(9223372036854775807)

	fmt.Printf("Int64: %d -> Float64: %.0f\n", i1, float64(i1))
	fmt.Printf("Int64: %d -> Float64: %.0f\n", i2, float64(i2))
	fmt.Printf("Are they equal as float64? %v\n", float64(i1) == float64(i2))
}
