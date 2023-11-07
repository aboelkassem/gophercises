package deck

import (
	"fmt"
)

// to document examples

func ExampleNew() {
	deck := New(
		OptionShuffle(),
	)
	fmt.Println(deck)
}
