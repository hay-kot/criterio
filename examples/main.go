package main

import (
	"fmt"

	"github.com/hay-kot/criterio"
)

func main() {
	fmt.Println("=== Criterio Validation Examples ===")
	fmt.Println()

	exampleStructValidation()
	exampleNestedErrors()
	exampleConditionalValidation()
	exampleCombinators()
	exampleRunAllVsRun()
}

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

type Address struct {
	Street string
	City   string
	Zip    string
}

func (a Address) Validate() error {
	return criterio.ValidateStruct(
		criterio.Run("street", a.Street, criterio.Required[string]),
		criterio.Run("city", a.City, criterio.Required[string]),
		criterio.Run("zip", a.Zip, criterio.Required[string], criterio.StrNumeric),
	)
}

type User struct {
	Name      string
	Email     string
	Age       int
	Phone     string
	IsPremium bool
	Address   Address
}

func (u User) Validate() error {
	return criterio.ValidateStruct(
		criterio.Run("name", u.Name,
			criterio.Required[string],
			criterio.StrBetween(2, 100),
		),
		criterio.Run("email", u.Email,
			criterio.Required[string],
			criterio.StrEmail,
		),
		criterio.Run("age", u.Age,
			criterio.Min(0),
			criterio.Max(150),
		),
		criterio.Run("phone", u.Phone,
			criterio.When(u.IsPremium, criterio.Required[string]),
		),
		criterio.Nest("address", u.Address.Validate()),
	)
}

// -----------------------------------------------------------------------------
// Examples
// -----------------------------------------------------------------------------

func exampleStructValidation() {
	fmt.Println("1. Struct Validation")
	fmt.Println("   -----------------")

	validUser := User{
		Name:  "Hayden",
		Email: "hayden@example.com",
		Age:   30,
		Address: Address{
			Street: "123 Main St",
			City:   "Austin",
			Zip:    "78701",
		},
	}

	if err := validUser.Validate(); err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Println("   Valid user passed validation")
	}
	fmt.Println()
}

func exampleNestedErrors() {
	fmt.Println("2. Nested Struct Errors")
	fmt.Println("   --------------------")

	invalidUser := User{
		Name:  "H", // too short
		Email: "not-an-email",
		Age:   -5, // negative
		Address: Address{
			Street: "", // missing
			City:   "Austin",
			Zip:    "abc", // not numeric
		},
	}

	if err := invalidUser.Validate(); err != nil {
		// Show individual field errors
		if fieldErrs, ok := err.(criterio.FieldErrors); ok {
			for _, fe := range fieldErrs {
				fmt.Printf("   - %s: %s\n", fe.Field, fe.Message)
			}
		}
	}
	fmt.Println()
}

func exampleConditionalValidation() {
	fmt.Println("3. Conditional Validation (When/SkipIf)")
	fmt.Println("   ------------------------------------")

	// Phone not required for regular users
	regularUser := User{
		Name:      "Regular",
		Email:     "regular@example.com",
		Age:       25,
		Phone:     "", // empty phone OK
		IsPremium: false,
		Address:   Address{Street: "1 St", City: "NYC", Zip: "10001"},
	}

	if err := regularUser.Validate(); err != nil {
		fmt.Printf("   Regular user error: %v\n", err)
	} else {
		fmt.Println("   Regular user: phone not required (passed)")
	}

	// Phone required for premium users
	premiumUser := User{
		Name:      "Premium",
		Email:     "premium@example.com",
		Age:       30,
		Phone:     "", // empty phone NOT OK
		IsPremium: true,
		Address:   Address{Street: "1 St", City: "NYC", Zip: "10001"},
	}

	if err := premiumUser.Validate(); err != nil {
		if fieldErrs, ok := err.(criterio.FieldErrors); ok {
			for _, fe := range fieldErrs {
				fmt.Printf("   Premium user error - %s: %s\n", fe.Field, fe.Message)
			}
		}
	}
	fmt.Println()
}

func exampleCombinators() {
	fmt.Println("4. Validator Combinators (Or, Not)")
	fmt.Println("   --------------------------------")

	// Or: accept email OR numeric string
	contactValidator := criterio.Or(criterio.StrEmail, criterio.StrNumeric)

	fmt.Print("   'test@example.com' (Or email|numeric): ")
	if err := contactValidator("test@example.com"); err != nil {
		fmt.Printf("failed - %v\n", err)
	} else {
		fmt.Println("passed")
	}

	fmt.Print("   '5551234567' (Or email|numeric): ")
	if err := contactValidator("5551234567"); err != nil {
		fmt.Printf("failed - %v\n", err)
	} else {
		fmt.Println("passed")
	}

	fmt.Print("   'invalid' (Or email|numeric): ")
	if err := contactValidator("invalid"); err != nil {
		fmt.Printf("failed - %v\n", err)
	} else {
		fmt.Println("passed")
	}

	// Not: username cannot be all numbers
	usernameValidator := criterio.Not(criterio.StrNumeric, "cannot be all numbers")

	fmt.Print("   'user123' (Not numeric): ")
	if err := usernameValidator("user123"); err != nil {
		fmt.Printf("failed - %v\n", err)
	} else {
		fmt.Println("passed")
	}

	fmt.Print("   '123456' (Not numeric): ")
	if err := usernameValidator("123456"); err != nil {
		fmt.Printf("failed - %v\n", err)
	} else {
		fmt.Println("passed")
	}
	fmt.Println()
}

func exampleRunAllVsRun() {
	fmt.Println("5. Run (fail-fast) vs RunAll (collect all)")
	fmt.Println("   ----------------------------------------")

	value := ""

	// Run stops on first error
	fmt.Println("   Run (fail-fast):")
	if err := criterio.Run("field", value,
		criterio.Required[string],
		criterio.StrMin(5),
		criterio.StrEmail,
	); err != nil {
		if fieldErrs, ok := err.(criterio.FieldErrors); ok {
			fmt.Printf("   Errors collected: %d\n", len(fieldErrs))
			for _, fe := range fieldErrs {
				fmt.Printf("     - %s\n", fe.Message)
			}
		}
	}

	// RunAll collects all errors
	fmt.Println("   RunAll (collect all):")
	if err := criterio.RunAll("field", value,
		criterio.Required[string],
		criterio.StrMin(5),
		criterio.StrEmail,
	); err != nil {
		if fieldErrs, ok := err.(criterio.FieldErrors); ok {
			fmt.Printf("   Errors collected: %d\n", len(fieldErrs))
			for _, fe := range fieldErrs {
				fmt.Printf("     - %s\n", fe.Message)
			}
		}
	}
}
