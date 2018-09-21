package list

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"
)

// GenList generates a list of string IPs of length len
func GenList(len int) []string {
	var retList []string
	for i := 0; i < len; i++ {
		retList = append(retList, GenIP())
	}

	return retList
}

// GenIP creates an IP address
func GenIP() string {
	var ip bytes.Buffer
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := strconv.Itoa(randGen.Intn(256))
	ip.WriteString(s)
	for i := 0; i < 3; i++ {
		ip.WriteString(".")
		s = strconv.Itoa(randGen.Intn(256))
		ip.WriteString(s)

	}

	return ip.String()
}
