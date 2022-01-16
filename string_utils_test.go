package main

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestWrapToLengthCmd(t *testing.T) {
	RegisterTestingT(t)
	s := "lorem ipsum dolor si amet and"

	ws := wrapToLength(s, 10)

	Expect(ws).To(Equal("lorem ipsu\nm dolor si\n amet and"))

}
