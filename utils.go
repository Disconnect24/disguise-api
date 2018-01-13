package disguise

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// https://stackoverflow.com/a/31832326/3874884
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStringBytesMaskImprSrc makes a random string with the specified size.
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// GenMailErrorCode formulates a proper response needed for mail-specific errors.
func GenMailErrorCode(mailNumber string, error int, reason string) string {
	if error != 100 {
		log.Println("[Warning] Encountered error", error, "with reason", reason)
	}

	return fmt.Sprint(
		"cd", mailNumber[1:], "=", strconv.Itoa(error), "\n",
		"msg", mailNumber[1:], "=", reason, "\n")
}

// GenNormalErrorCode formulates a proper response for overall errors.
func GenNormalErrorCode(error int, reason string) string {
	if error != 100 {
		log.Println("[Warning] Encountered error", error, "with reason", reason)
	}
	return fmt.Sprint(
		"cd=", strconv.Itoa(error), "\n",
		"msg=", reason, "\n")
}

// FriendCodeIsValid determines if a friend code is valid by
// checking not empty, is 17 in length, starts with w.
// BUG(spotlightishere): does not actually determine at a numerical level if valid.
func FriendCodeIsValid(wiiID string) bool {
	if wiiID == "" {
		return false
	} else if wiiID[0:1] != "w" {
		return false
	} else if len(wiiID) != 17 {
		return false
	} else {
		return true
	}
}

// GenerateBoundary returns a string in the Wii specific boundary format.
func GenerateBoundary() string {
	return fmt.Sprint("BoundaryForDL", fmt.Sprint(time.Now().Format("200601021504")), "/", random(1000000, 9999999))
}
