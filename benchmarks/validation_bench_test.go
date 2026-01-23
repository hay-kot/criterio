package benchmarks

import (
	"errors"
	"net/mail"
	"regexp"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hay-kot/criterio"
)

// ============================================================================
// Test Data Structures
// ============================================================================

type User struct {
	Name  string `validate:"required,min=2,max=100"`
	Email string `validate:"required,email"`
	Age   int    `validate:"required,gte=0,lte=150"`
}

type Address struct {
	Street  string `validate:"required,min=1,max=200"`
	City    string `validate:"required,min=1,max=100"`
	Country string `validate:"required,len=2"`
	Zip     string `validate:"required,min=5,max=10"`
}

type Order struct {
	ID       string   `validate:"required,uuid"`
	Customer User     `validate:"required"`
	Shipping Address  `validate:"required"`
	Items    []string `validate:"required,min=1,dive,required,min=1"`
}

// ============================================================================
// Test Data
// ============================================================================

var validUser = User{
	Name:  "John Doe",
	Email: "john@example.com",
	Age:   30,
}

var invalidUser = User{
	Name:  "J",
	Email: "invalid-email",
	Age:   -1,
}

var validAddress = Address{
	Street:  "123 Main St",
	City:    "New York",
	Country: "US",
	Zip:     "10001",
}

var validOrder = Order{
	ID:       "550e8400-e29b-41d4-a716-446655440000",
	Customer: validUser,
	Shipping: validAddress,
	Items:    []string{"item1", "item2", "item3"},
}

var invalidOrder = Order{
	ID:       "not-a-uuid",
	Customer: invalidUser,
	Shipping: Address{},
	Items:    []string{},
}

// ============================================================================
// Manual Validation
// ============================================================================

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func validateEmailManual(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("email is invalid")
	}
	return nil
}

func validateUserManual(u User) error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if len(u.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	if len(u.Name) > 100 {
		return errors.New("name must be at most 100 characters")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("email is invalid")
	}
	if u.Age < 0 {
		return errors.New("age must be non-negative")
	}
	if u.Age > 150 {
		return errors.New("age must be at most 150")
	}
	return nil
}

func validateAddressManual(a Address) error {
	if a.Street == "" {
		return errors.New("street is required")
	}
	if len(a.Street) > 200 {
		return errors.New("street must be at most 200 characters")
	}
	if a.City == "" {
		return errors.New("city is required")
	}
	if len(a.City) > 100 {
		return errors.New("city must be at most 100 characters")
	}
	if a.Country == "" {
		return errors.New("country is required")
	}
	if len(a.Country) != 2 {
		return errors.New("country must be 2 characters")
	}
	if a.Zip == "" {
		return errors.New("zip is required")
	}
	if len(a.Zip) < 5 || len(a.Zip) > 10 {
		return errors.New("zip must be 5-10 characters")
	}
	return nil
}

func validateOrderManual(o Order) error {
	if o.ID == "" {
		return errors.New("id is required")
	}
	if !uuidRegex.MatchString(o.ID) {
		return errors.New("id must be a valid UUID")
	}
	if err := validateUserManual(o.Customer); err != nil {
		return err
	}
	if err := validateAddressManual(o.Shipping); err != nil {
		return err
	}
	if len(o.Items) == 0 {
		return errors.New("items is required")
	}
	for _, item := range o.Items {
		if item == "" {
			return errors.New("item cannot be empty")
		}
	}
	return nil
}

// ============================================================================
// go-playground/validator Setup
// ============================================================================

var playgroundValidator = validator.New()

// ============================================================================
// Criterio Validators
// ============================================================================

func validateEmailCriterio(email string) error {
	return criterio.Run("email", email,
		criterio.StrNotEmpty,
		criterio.StrEmail,
	)
}

func validateUserCriterio(u User) error {
	return criterio.ValidateStruct(
		criterio.Run("name", u.Name,
			criterio.StrNotEmpty,
			criterio.StrMin(2),
			criterio.StrMax(100),
		),
		criterio.Run("email", u.Email,
			criterio.StrNotEmpty,
			criterio.StrEmail,
		),
		criterio.Run("age", u.Age,
			criterio.Min(0),
			criterio.Max(150),
		),
	)
}

func validateAddressCriterio(a Address) error {
	return criterio.ValidateStruct(
		criterio.Run("street", a.Street,
			criterio.StrNotEmpty,
			criterio.StrMax(200),
		),
		criterio.Run("city", a.City,
			criterio.StrNotEmpty,
			criterio.StrMax(100),
		),
		criterio.Run("country", a.Country,
			criterio.StrNotEmpty,
			criterio.StrBetween(2, 2),
		),
		criterio.Run("zip", a.Zip,
			criterio.StrNotEmpty,
			criterio.StrBetween(5, 10),
		),
	)
}

func validateOrderCriterio(o Order) error {
	return criterio.ValidateStruct(
		criterio.Run("id", o.ID,
			criterio.StrNotEmpty,
			criterio.StrUUID,
		),
		criterio.Nest("customer", validateUserCriterio(o.Customer)),
		criterio.Nest("shipping", validateAddressCriterio(o.Shipping)),
		criterio.Run("items", o.Items,
			criterio.SliceNotEmpty[string](),
			criterio.SliceEach(criterio.StrNotEmpty),
		),
	)
}

// ============================================================================
// Benchmarks: Simple Email Validation
// ============================================================================

func BenchmarkEmail_Manual_Valid(b *testing.B) {
	email := "john@example.com"
	b.ResetTimer()
	for b.Loop() {
		_ = validateEmailManual(email)
	}
}

func BenchmarkEmail_Playground_Valid(b *testing.B) {
	email := "john@example.com"
	b.ResetTimer()
	for b.Loop() {
		_ = playgroundValidator.Var(email, "required,email")
	}
}

func BenchmarkEmail_Criterio_Valid(b *testing.B) {
	email := "john@example.com"
	b.ResetTimer()
	for b.Loop() {
		_ = validateEmailCriterio(email)
	}
}

func BenchmarkEmail_Manual_Invalid(b *testing.B) {
	email := "invalid-email"
	b.ResetTimer()
	for b.Loop() {
		_ = validateEmailManual(email)
	}
}

func BenchmarkEmail_Playground_Invalid(b *testing.B) {
	email := "invalid-email"
	b.ResetTimer()
	for b.Loop() {
		_ = playgroundValidator.Var(email, "required,email")
	}
}

func BenchmarkEmail_Criterio_Invalid(b *testing.B) {
	email := "invalid-email"
	b.ResetTimer()
	for b.Loop() {
		_ = validateEmailCriterio(email)
	}
}

// ============================================================================
// Benchmarks: Simple Struct Validation (User)
// ============================================================================

func BenchmarkUser_Manual_Valid(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = validateUserManual(validUser)
	}
}

func BenchmarkUser_Playground_Valid(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = playgroundValidator.Struct(validUser)
	}
}

func BenchmarkUser_Criterio_Valid(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = validateUserCriterio(validUser)
	}
}

func BenchmarkUser_Manual_Invalid(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = validateUserManual(invalidUser)
	}
}

func BenchmarkUser_Playground_Invalid(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = playgroundValidator.Struct(invalidUser)
	}
}

func BenchmarkUser_Criterio_Invalid(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = validateUserCriterio(invalidUser)
	}
}

// ============================================================================
// Benchmarks: Nested Struct Validation (Order)
// ============================================================================

func BenchmarkOrder_Manual_Valid(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = validateOrderManual(validOrder)
	}
}

func BenchmarkOrder_Playground_Valid(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = playgroundValidator.Struct(validOrder)
	}
}

func BenchmarkOrder_Criterio_Valid(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = validateOrderCriterio(validOrder)
	}
}

func BenchmarkOrder_Manual_Invalid(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = validateOrderManual(invalidOrder)
	}
}

func BenchmarkOrder_Playground_Invalid(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = playgroundValidator.Struct(invalidOrder)
	}
}

func BenchmarkOrder_Criterio_Invalid(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = validateOrderCriterio(invalidOrder)
	}
}

// ============================================================================
// Benchmarks: Slice Validation
// ============================================================================

func BenchmarkSlice_Manual_Valid(b *testing.B) {
	items := []string{"item1", "item2", "item3", "item4", "item5"}
	b.ResetTimer()
	for b.Loop() {
		for _, item := range items {
			if item == "" {
				break
			}
		}
	}
}

func BenchmarkSlice_Playground_Valid(b *testing.B) {
	items := []string{"item1", "item2", "item3", "item4", "item5"}
	b.ResetTimer()
	for b.Loop() {
		_ = playgroundValidator.Var(items, "required,min=1,dive,required,min=1")
	}
}

func BenchmarkSlice_Criterio_Valid(b *testing.B) {
	items := []string{"item1", "item2", "item3", "item4", "item5"}
	b.ResetTimer()
	for b.Loop() {
		_ = criterio.Run("items", items,
			criterio.SliceNotEmpty[string](),
			criterio.SliceEach(criterio.StrNotEmpty),
		)
	}
}

// ============================================================================
// Benchmarks: Parallel Execution (measures contention)
// ============================================================================

func BenchmarkUser_Manual_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = validateUserManual(validUser)
		}
	})
}

func BenchmarkUser_Playground_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundValidator.Struct(validUser)
		}
	})
}

func BenchmarkUser_Criterio_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = validateUserCriterio(validUser)
		}
	})
}
