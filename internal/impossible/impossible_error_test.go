package impossible

import (
	"fmt"
	"testing"
)

func TestHandleError(t *testing.T) {
	defer (func() {
		r := recover()
		if r == nil {
			t.Fatalf("did not panic")
		}

		switch x := r.(type) {
		case string:
			if x != "impossible error: abc" {
				t.Fatalf("error not as expected: %v", x)
			}
		default:
			t.Fatalf("unknown panic: %T:%v", x, x)
		}
	})()

	HandleError(fmt.Errorf("abc"))
}
