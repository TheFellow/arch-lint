package beta

import (
	_ "example/patterns/orders" // want `\[pattern boundaries\] forbidden import`
	_ "example/patterns/orders-v2"
)
