package rapid7_test

import (
	"fmt"
	"testing"

	"github.com/stellaraf/go-rapid7"
	"github.com/stretchr/testify/assert"
)

func Test_InvestigationIDFromRRN(t *testing.T) {
	cases := [][]string{
		{"", ""},
		{"rrn:investigation:", ""},
		{"rrn:stuff", ""},
	}
	t.Run("correct", func(t *testing.T) {
		t.Parallel()
		id := "0287ca61-d643-488d-9613-b6a12c736bfa"
		inv := &rapid7.Investigation{
			RRN: fmt.Sprintf("rrn:investigation:us2:%s:investigation:00WQ1YRXS8AF", id),
		}
		result := inv.ID()
		assert.Equal(t, id, result)
	})
	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("incorrect-%d", i), func(t *testing.T) {
			t.Parallel()
			i := &rapid7.Investigation{
				RRN: c[0],
			}
			result := i.ID()
			assert.Equal(t, c[1], result)
		})
	}
}
