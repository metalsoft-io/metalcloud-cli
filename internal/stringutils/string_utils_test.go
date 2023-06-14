package stringutils

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestWrapToLengthCmd(t *testing.T) {
	RegisterTestingT(t)
	s := "lorem ipsum dolor si amet and"

	ws := WrapToLength(s, 10)

	Expect(ws).To(Equal("lorem ipsu\nm dolor si\n amet and"))

}

func TestMakeLabel(t *testing.T) {
	RegisterTestingT(t)

	label, err := MakeLabel("Test11")
	Expect(err).To(BeNil())
	Expect(label).To(Equal("test11"))

	label, err = MakeLabel("Tes-t11")
	Expect(err).To(BeNil())
	Expect(label).To(Equal("tes-t11"))

	label, err = MakeLabel("!Tes-t11")
	Expect(err).NotTo(BeNil())

	label, err = MakeLabel("Tes-t11 ")
	Expect(err).To(BeNil())
	Expect(label).To(Equal("tes-t11"))

	label, err = MakeLabel("1Tes-t11")
	Expect(err).NotTo(BeNil())

	label, err = MakeLabel("Tes-$t11")
	Expect(err).To(BeNil())
	Expect(label).To(Equal("tes-t11"))

	label, err = MakeLabel("$#!")
	Expect(err).NotTo(BeNil())

	label, err = MakeLabel("")
	Expect(err).NotTo(BeNil())
}
