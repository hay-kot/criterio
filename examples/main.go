package main

import (
	"log"

	"github.com/hay-kot/criterio"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("failed to run program, %s", err.Error())
	}
}

func run() error {
	validateName := criterio.New("name",
		criterio.Required[string](),
		criterio.StrBetween(8, 128),
		criterio.OneOf("Hayden L", "Hayden J", "Hayden F"),
	)

	err := criterio.Run("name", "Hayden K",
		criterio.Required[string](),
		criterio.StrBetween(8, 128),
		criterio.OneOf("Hayden L", "Hayden J", "Hayden F"),
	)

	validateName("Hayden K")

	return err
}
