package protector

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"unicode"
)

// GetSessionKey generate 10 char random string
func GetSessionKey() string {
	rand.Seed((time.Now().Unix()))
	var digits = []rune("0123456789")
	result := make([]rune, 10)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}

// GetHashStr calculate initial hash string
func GetHashStr() string {
	rand.Seed((time.Now().Unix()))
	var digits = []rune("0123456789")
	li := make([]rune, 6)
	for i := range li {
		li[i] = digits[rand.Intn(len(digits))]
	}
	return string(li)
}

// NextSessionKey generate next session key
func NextSessionKey(hashString string, sessionKey string) string {

	// verify hashcode

	if hashString == "" {
		return "Hash code is empty"
	}

	runes := []rune(hashString)

	for idx := range runes {
		i := runes[idx]
		if !unicode.IsDigit(i) {
			return fmt.Sprintf("Hash code contains non-digit letter \"%c\"", i)
		}
	}
	num := 0
	for idx := 0; idx < len(runes); idx++ {
		i := runes[idx]
		hashNum, _ := strconv.Atoi(calcHash(sessionKey, int(i)))
		num += hashNum
	}
	str := "0000000000" + strconv.Itoa(num)[0:10]
	return str[len(str)-10:]
}

func calcHash(sessionKey string, value int) string {
	result := ""
	switch value {
	case 1:
		i, _ := strconv.Atoi(sessionKey[0:5])
		str := "00" + strconv.Itoa(i%97)
		return str[len(str)-2:]
	case 2:
		for i := 1; i < len(sessionKey); i++ {
			result = result + string(sessionKey[len(sessionKey)-i])
		}
		return result + string(sessionKey[0])
	case 3:
		return string(sessionKey[len(sessionKey)-5]) + sessionKey[0:5]
	case 4:
		num := 0
		for i := 1; i < 9; i++ {
			x, _ := strconv.Atoi(string(sessionKey[i]))
			num += x + 41
		}
		return strconv.Itoa(num)
	case 5:
		num := 0
		for i := 0; i < len(sessionKey); i++ {
			ch := rune(int(sessionKey[i]) + 41)
			if !unicode.IsDigit(ch) {
				ch = rune(strconv.Itoa(int(ch))[0])
			}
			num += int(ch)
		}
		return strconv.Itoa(num)
	default:
		i, _ := strconv.Atoi(sessionKey)
		return strconv.Itoa(i + value)
	}
}
