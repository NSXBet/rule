package rule

// Hash constants for FNV-1a hash algorithm.
const (
	// hashOffsetBasis is the FNV-1a 32-bit offset basis.
	hashOffsetBasis uint32 = 5381
	// hashPrime is the FNV-1a 32-bit prime.
	hashPrime = 5
)

// Character and parsing constants.
const (
	// tokenSliceInitialCapacity is the initial capacity for token slices.
	tokenSliceInitialCapacity = 32
)

// Large integer precision constants.
const (
	// maxSafeInteger is the maximum integer that can be represented safely in float64 (2^53).
	maxSafeInteger int64 = 9007199254740992
	// minSafeInteger is the minimum integer that can be represented safely in float64 (-2^53).
	minSafeInteger int64 = -9007199254740992
)

// String constants.
const (
	// trueString represents the string "true".
	trueString = "true"
)
