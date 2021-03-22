package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	// Password salt number of bytes.
	ARGON2_SALT_LEN = 16

	// Argon2 time parameter.
	ARGON2_TIME = 1

	// Argon2 memory parameter (in KiB).
	ARGON2_MEMORY = 64 * 1024

	// Argon2 thread parameter.
	ARGON2_THREADS = 4

	// Length of the hashed key.
	ARGON2_KEY_LEN = 32

	// Stored password format string.
	// Format: "<algorithm>$<algorithm_ver>$<time>$<memory>$<threads>$<salt>$<hash>"
	ARGON2_FORMAT = "argon2id$%d$%d$%d$%d$%s$%s"

	// Separator between parameters.
	ARGON2_FORMAT_SEP = "$"

	// How many parameters are stored in the ARGON2_FORMAT
	// argon2id + version + time + memory + threads + salt + hash
	ARGON2_FORMAT_COUNT = 7
)

// Generate an argon2 hash of the provided plaintext password using the parameters
// described above.
func Password(plaintext string) (string, error) {
	// create a salt
	salt := make([]byte, ARGON2_SALT_LEN)

	// fill the salt with secure random bytes
	_, e := rand.Read(salt)

	if e != nil {
		return "", nil
	}

	// hash the password
	hash := argon2.IDKey(
		[]byte(plaintext),
		salt,
		ARGON2_TIME,
		ARGON2_MEMORY,
		ARGON2_THREADS,
		ARGON2_KEY_LEN,
	)

	// encode the hash and salt as base64 strings
	hash64 := base64.StdEncoding.EncodeToString(hash)
	salt64 := base64.StdEncoding.EncodeToString(salt)

	// format the parameters used to create the hash along with the actual hash
	// i.e. argon2 parameters and the salt
	return fmt.Sprintf(
		ARGON2_FORMAT,
		argon2.Version,
		ARGON2_TIME,
		ARGON2_MEMORY,
		ARGON2_THREADS,
		salt64,
		hash64,
	), nil
}

// Compare an encoded password using the encoding described in ARGON2_FORMAT with
// a plaintext password.
func CmpPassword(encoded, plaintext string) (bool, error) {
	// decode the encoded string
	// 0: algorithm
	// 1: algorithm version
	// 2: time parameter
	// 3: memory parameter
	// 4: threads parameter
	// 5: salt
	// 6: hash
	slice := strings.Split(encoded, ARGON2_FORMAT_SEP)

	if l := len(slice); l != ARGON2_FORMAT_COUNT {
		return false, fmt.Errorf(
			"Encoding mis-match. Expected: %d. Actual: %d.", ARGON2_FORMAT_COUNT, l)
	}

	// slice from time parameter to threads parameter
	time, memory, threads, e := strParams(slice[2:5])

	if e != nil {
		return false, e
	}

	// decode the salt to a byte slice from base64
	salt, e := base64.StdEncoding.DecodeString(slice[5])

	if e != nil {
		return false, e
	}

	// decode the hash to a byte slice from base64
	hash, e := base64.StdEncoding.DecodeString(slice[6])

	if e != nil {
		return false, e
	}

	// hash the plaintext comparison password
	cmp := argon2.IDKey([]byte(plaintext), salt, time, memory, threads, uint32(len(hash)))

	// compare the byte slices
	return subtle.ConstantTimeCompare(hash, cmp) == 1, nil
}

// Convert the three argon2 cost complexity strings to their respective typed integers.
// Assumes: 0: time, 1: memory, 2: threads.
func strParams(slice []string) (uint32, uint32, uint8, error) {
	// convert the time parameter to uint32
	time, e := strUint32(slice[0])

	if e != nil {
		return 0, 0, 0, e
	}

	// convert the memory parameter to uint32
	memory, e := strUint32(slice[1])

	if e != nil {
		return 0, 0, 0, e
	}

	// convert the threads parameter to uint8
	threads, e := strconv.ParseInt(slice[2], 10, 8)

	if e != nil {
		return 0, 0, 0, e
	}

	return time, memory, uint8(threads), nil
}

// Convert a string to a uint32 (base 10).
func strUint32(str string) (uint32, error) {
	int64, e := strconv.ParseInt(str, 10, 32)

	return uint32(int64), e
}
