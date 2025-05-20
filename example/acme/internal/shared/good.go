package shared

import (
	"fmt"

	"github.com/TheFellow/go-arch-lint/example/acme/experimental"
)

func Test() {
	w := experimental.NewWidget()
	fmt.Println(w)
}
