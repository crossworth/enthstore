package enthstore

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_quoteKey(t *testing.T) {
	t.Parallel()
	tests := []struct {
		args string
		want string
	}{
		{args: "a", want: "'a'"},
		{args: "a'b", want: "'a''b'"},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			if got := quoteKey(tt.args); got != tt.want {
				t.Errorf("quoteKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_quoteValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		args string
		want string
	}{
		{args: `test`, want: `"test"`},
		{args: `te st`, want: `"te st"`},
		{args: `"a"`, want: `"\"a\""`},
		{args: `\"a\"`, want: `"\\\"a\\\""`},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, quoteValue(tt.args))
		})
	}
}
