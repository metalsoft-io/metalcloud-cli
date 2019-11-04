package main

import (
	"fmt"
	"math/rand"
	"testing"

	. "github.com/onsi/gomega"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestValidateAPIKey(t *testing.T) {
	RegisterTestingT(t)

	Expect(len(RandStringBytes(64))).To(Equal(64))
	goodKey := fmt.Sprintf("%d:%s", rand.Intn(100), RandStringBytes(63))

	badKey1 := fmt.Sprintf("%d:%s", rand.Intn(100), RandStringBytes(64))
	badKey2 := fmt.Sprintf(":%s", RandStringBytes(63))

	Expect(validateAPIKey(goodKey)).To(BeNil())
	Expect(validateAPIKey(badKey1)).NotTo(BeNil())
	Expect(validateAPIKey(badKey2)).NotTo(BeNil())

}
