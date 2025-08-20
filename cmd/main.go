package main

import (
	"github.com/noxhalley/funken/internal/initializer"
	"go.uber.org/fx"
)

func main() {
	fx.New(initializer.Build()).Run()
}
