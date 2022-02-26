package bootstrap

import (
	"fmt"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(bootstrap),
)

func bootstrap() {
	fmt.Println("Bootstrap !")
}
