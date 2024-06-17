package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	mrand "math/rand"
	"time"
)

// formatDateForAvro converts a datetime object to an integer for AVRO encoding
func FormatDateForAvro(inputDate time.Time, formatType string) int64 {
	utcDate := inputDate.UTC()
	epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	if formatType == "timestamp-millis" {
		return utcDate.Sub(epoch).Milliseconds()
	}
	if formatType == "date" {
		return int64(utcDate.Sub(epoch).Hours() / 24)
	}

	panic("Invalid formatType. Choose 'timestamp-millis' or 'date'.")
}

// FormatDateForJson converts a datetime object to a string for JSON encoding
func FormatDateForJson(inputDate time.Time, formatType string) string {
	utcDate := inputDate.UTC()

	switch formatType {
	case "timestamp-millis":
		// Return Unix time in milliseconds as a string
		return fmt.Sprintf("%d", utcDate.UnixNano()/int64(time.Millisecond))

	case "date":
		// Return date in YYYY-MM-DD format
		return utcDate.Format("2006-01-02")

	case "iso8601":
		// Return ISO 8601 formatted string
		return utcDate.Format(time.RFC3339)

	default:
		panic("Invalid formatType. Choose 'timestamp-millis', 'date', or 'iso8601'.")
	}
}

// generateRandomString generates a random string with a given prefix and specified length
func GenerateRandomString(prefix string, length int) string {
	return prefix + RandomString(length)
}

// randomString generates a random string of length n.
//
// It takes an integer parameter n, which specifies the length of the generated string.
// The function returns a string.
func RandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		slog.Error("Error generating random string", "Error", err)
	}
	return base64.URLEncoding.EncodeToString(b)[:n]
}

// RandomInt generates a random integer within the specified range [0, max)
func RandomInt(max int) int {
	var n uint32
	if err := binary.Read(rand.Reader, binary.LittleEndian, &n); err != nil {
		return 0
	}

	// Use modulo to limit the range of the random number
	return int(n % uint32(max))
}

// RandomFloat generates a random float64 value between min and max.
func RandomFloat(min, max float64) (float64, error) {
	if min >= max {
		return 0, errors.New("min must be less than max")
	}

	// Determine the range of random numbers
	rangeSize := max - min

	// Get the number of bits required for the range
	bitCount := uint(math.Ceil(math.Log2(rangeSize)))

	// Generate a random big.Int with the specified number of bits
	randomInt, err := rand.Int(rand.Reader, big.NewInt(1).Lsh(big.NewInt(1), bitCount))
	if err != nil {
		return 0, err
	}

	// Convert the big.Int to a float64 value and scale it to the desired range
	randomFloat := new(big.Float).SetInt(randomInt)
	rangeBigFloat := new(big.Float).SetFloat64(rangeSize)
	randomFloat.Mul(randomFloat, rangeBigFloat)
	floatResult, _ := randomFloat.Float64()

	// Normalize the float within the min-max range
	floatResult = floatResult/(math.Pow(2, float64(bitCount)))*rangeSize + min

	return floatResult, nil
}

func DateOrNil() *int {
	if mrand.Float64() < 0.2 { // 20% chance to be less than 0.2
		return nil
	}

	myDate := int(FormatDateForAvro(GenerateRandomDate(), "date"))
	myDatePointer := &myDate
	return myDatePointer
}

// RandomFloatOrNil returns a random float64 pointer or nil based on the specified probability.
func RandomFloatOrNil(min, max float64, nilProbability float64) *float64 {
	// Ensure valid probability range
	if nilProbability < 0.0 || nilProbability > 1.0 {
		panic("nilProbability must be between 0.0 and 1.0")
	}

	// Determine if we should return nil or a valid float
	if mrand.Float64() < nilProbability {
		return nil
	}

	// Generate a random float within the specified range
	randomFloat, err := RandomFloat(min, max)
	if err != nil {
		panic(fmt.Sprintf("Error generating random float %v", err))
	}

	return &randomFloat
}

// IfEmptyReturnNilString accepts a string or a pointer to a string, returning a pointer to a non-empty string or nil if empty.
func IfEmptyReturnNilString(data interface{}) *string {
	var str *string
	switch v := data.(type) {
	case *string:
		// If it's a pointer, check if it's nil or its value is empty
		if v == nil || *v == "" {
			return nil
		}
		str = v
	case string:
		// If it's a direct string, check if it's empty
		if v == "" {
			return nil
		}
		str = &v
	default:
		panic("Invalid input type. Expected string or *string")
	}
	return str
}

// RandomChoiceString selects a random non-empty string from a given slice of strings
func RandomChoiceString(choices []string) *string {
	if len(choices) == 0 {
		return nil
	}
	// Get a random index within the range of the choices
	index := RandomInt(len(choices))
	choice := choices[index]
	if choice == "" {
		return nil
	}
	return &choice
}

// RandomChoiceInt selects a random integer from a given slice of integers
func RandomChoiceInt(choices []int) *int {
	if len(choices) == 0 {
		return nil
	}
	// Get a random index
	index := RandomInt(len(choices))
	choice := choices[index]
	if choice == 0 {
		return nil
	}
	return &choice
}

// RandomInt64 generates a random int64 within the specified range [0, max)
func RandomInt64(max int64) (int64, error) {
	var n uint64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &n); err != nil {
		return 0, err
	}

	// Use modulo to limit the range of the random number
	return int64(n % uint64(max)), nil
}

func RandomBoolPointer() *bool {
	// Generate a random boolean (true or false)
	randomValue := RandomInt(2) == 1
	// Return a pointer to the random boolean
	return &randomValue
}

func RandomIntPointer(max int) *int {
	// Generate a random integer within the specified range
	randomValue := RandomInt(max)
	// Return a pointer to the random integer
	return &randomValue
}

// GenerateRandomDate generates a random datetime object within the last 50 years
func GenerateRandomDate() time.Time {
	start := time.Now().AddDate(-50, 0, 0) // 50 years back from the current time
	durationToNow := time.Since(start)     // Total duration from 50 years ago to now

	// Get a secure random duration
	randomDuration, err := RandomInt64(int64(durationToNow))
	if err != nil {
		return time.Time{}
	}

	// Return the random date by adding the random duration to the start
	return start.Add(time.Duration(randomDuration))
}

func NewString(str string) *string {
	return &str
}

// GetWeightedRandomDatabaseOperation returns a weighted random operation ("INSERT", "UPDATE", "DELETE")
// with a higher likelihood for "INSERT".
func GetWeightedRandomDatabaseOperation() string {
	// Generate a random number from 0 to 9
	randomValue := RandomInt(10)

	// Weighted probabilities for each outcome
	switch {
	case randomValue < 6: // 60% chance of INSERT
		return "INSERT"
	case randomValue < 8: // 20% chance of UPDATE
		return "UPDATE"
	case randomValue < 10: // 20% chance of DELETE
		return "DELETE"
	default:
		panic("Unexpected random value. Expected a number between 0 and 9.")
	}
}
