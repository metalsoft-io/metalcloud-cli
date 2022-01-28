package main

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestConvertToSearchFieldFormat(t *testing.T) {
	RegisterTestingT(t)

	Expect(convertToSearchFieldFormat("id:a,b,c status:available,unavailable")).To(
		Equal("+id:a +id:b +id:c +status:available +status:unavailable"))

	Expect(convertToSearchFieldFormat("id:a,b,c ")).To(
		Equal("+id:a +id:b +id:c"))

	Expect(convertToSearchFieldFormat("server_status:available,used")).To(
		Equal("+server_status:available +server_status:used"))

	Expect(convertToSearchFieldFormat("asasd")).To(
		Equal("asasd"))

}
