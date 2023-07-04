package filtering

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestConvertToSearchFieldFormat(t *testing.T) {
	RegisterTestingT(t)

	Expect(ConvertToSearchFieldFormat("id:a,b,c status:available,unavailable")).To(
		Equal("+id:a +id:b +id:c +status:available +status:unavailable"))

	Expect(ConvertToSearchFieldFormat("id:a,b,c ")).To(
		Equal("+id:a +id:b +id:c"))

	Expect(ConvertToSearchFieldFormat("server_status:available,used")).To(
		Equal("+server_status:available +server_status:used"))

	Expect(ConvertToSearchFieldFormat("asasd")).To(
		Equal("asasd"))

}
