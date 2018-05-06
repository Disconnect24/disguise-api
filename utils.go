package disguise

import (
	"context"
	"fmt"
	"google.golang.org/appengine/log"
	"math/rand"
	"strconv"
	"time"
	"regexp"
)

// https://stackoverflow.com/a/31832326/3874884
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var MailRegex = regexp.MustCompile(`w\d{16}`)

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
func GenMailErrorCode(ctx context.Context, mailNumber string, error int, reason string) string {
	if error != 100 {
		log.Warningf(ctx, "Encountered error", error, "with reason", reason)
	}

	return fmt.Sprint(
		"cd", mailNumber[1:], "=", strconv.Itoa(error), "\n",
		"msg", mailNumber[1:], "=", reason, "\n")
}

// GenNormalErrorCode formulates a proper response for overall errors.
func GenNormalErrorCode(ctx context.Context, error int, reason string) string {
	if error != 100 {
		log.Warningf(ctx, "Encountered error", error, "with reason", reason)
	}
	return fmt.Sprint(
		"cd=", strconv.Itoa(error), "\n",
		"msg=", reason, "\n")
}

// FriendCodeIsValid determines if a friend code is valid by
// checking not empty, is 17 in length, starts with w.
func FriendCodeIsValid(wiiID string) bool {
	return MailRegex.MatchString(wiiID)
}

// GenerateBoundary returns a string in the Wii specific boundary format.
func GenerateBoundary() string {
	return fmt.Sprint("BoundaryForDL", fmt.Sprint(time.Now().Format("200601021504")), "/", random(1000000, 9999999))
}
