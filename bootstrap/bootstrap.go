package bootstrap

import (
	"boilerplate/infrastructure"
	"fmt"

	"go.uber.org/fx"
)

var Module = fx.Options(
	infrastructure.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(env infrastructure.Env) {
	fmt.Println(env.DBName)
	fmt.Println("Bootstrap !")
}
