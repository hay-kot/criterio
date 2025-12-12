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
	// Create a reusable validator
	validateName := criterio.New("name",
		criterio.Required[string],
		criterio.StrBetween(8, 128),
		criterio.OneOf("Hayden L", "Hayden J", "Hayden F"),
	)

	validateCount := criterio.New("count",
		criterio.Min(1),
		criterio.Max(100),
	)

	// Run validators inline
	if err := criterio.Run("name", "Hayden K",
		criterio.Required[string],
		criterio.StrBetween(8, 128),
		criterio.OneOf("Hayden L", "Hayden J", "Hayden F"),
	); err != nil {
		return err
	}

	// Use reusable validators
	if err := validateName("Hayden L"); err != nil {
		return err
	}

	if err := validateCount(50); err != nil {
		return err
	}

	return nil
}
