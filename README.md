# criterio

A composable validation library for Go. Validators are functions that can be combined, reused, and composed to build complex validation logic.

```bash
go get github.com/hay-kot/criterio
```

Requires Go 1.25+

## Core Concepts

All validators implement a single type:

```go
type Validator[T any] func(val T) error
```

Validators return `nil` on success or an error describing the failure.

## API Overview

### Running Validators

```go
// Stop on first failure
err := criterio.Run("email", user.Email,
    criterio.Required[string],
    criterio.StrEmail,
)

// Collect all failures
err := criterio.RunAll("email", user.Email,
    criterio.Required[string],
    criterio.StrMin(5),
    criterio.StrEmail,
)
```

### Reusable Validators

```go
validateEmail := criterio.New("email",
    criterio.Required[string],
    criterio.StrEmail,
)

if err := validateEmail(input); err != nil {
    // handle error
}
```

### Struct Validation

```go
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
        criterio.Nest("address", u.Address.Validate()),
    )
}
```

### Conditional Validation

```go
// Only validate if condition is true
criterio.When(u.IsPremium, criterio.Required[string])

// Skip validation if condition is true
criterio.SkipIf(u.IsGuest, criterio.StrEmail)
```

### Slice Validation

```go
func (u User) Validate() error {
    return criterio.ValidateStruct(
        criterio.Run("name", u.Name, criterio.Required[string]),
        criterio.ValidateSlice("addresses", u.Addresses, Address.Validate),
    )
}
// Error: "addresses[0].street: is required"
```

For simple slices:

```go
criterio.ValidateSlice("emails", emails, criterio.StrEmail)
// Error: "emails[1]: must be a valid email address"
```

### Combinators

```go
// Pass if any validator passes
criterio.Or(criterio.StrEmail, criterio.StrNumeric)

// Invert a validator
criterio.Not(criterio.StrNumeric, "cannot be all numbers")
```

## Available Validators

### General

| Validator                      | Description             |
| ------------------------------ | ----------------------- |
| `Required[T comparable]`       | Value is non-zero       |
| `OneOf[T comparable](vals...)` | Value is in allowed set |

### Numeric (ordered types)

| Validator            | Description                |
| -------------------- | -------------------------- |
| `Min(n)`             | Value >= n                 |
| `Max(n)`             | Value <= n                 |
| `Between(low, high)` | Value in range [low, high] |
| `Positive`           | Value > 0                  |
| `Negative`           | Value < 0                  |
| `NonZero`            | Value != 0                 |
| `MultipleOf(n)`      | Value divisible by n       |

### Strings

| Validator               | Description                         |
| ----------------------- | ----------------------------------- |
| `StrNotEmpty`           | Non-empty after trimming whitespace |
| `StrMin(n)`             | At least n characters               |
| `StrMax(n)`             | At most n characters                |
| `StrBetween(low, high)` | Length in range [low, high]         |
| `StrMatches(re)`        | Matches regex pattern               |
| `StrEmail`              | Valid email format                  |
| `StrURL`                | Valid URL with scheme and host      |
| `StrUUID`               | Valid UUID format                   |
| `StrAlpha`              | Letters only                        |
| `StrAlphanumeric`       | Letters and numbers only            |
| `StrNumeric`            | Digits only                         |
| `StrContains(s)`        | Contains substring                  |
| `StrHasPrefix(s)`       | Starts with prefix                  |
| `StrHasSuffix(s)`       | Ends with suffix                    |
| `StrNoWhitespace`       | No whitespace characters            |
| `StrLowercase`          | All lowercase                       |
| `StrUppercase`          | All uppercase                       |
| `StrOneOf(vals...)`     | Value is in allowed set             |

### Slices

| Validator                    | Description                     |
| ---------------------------- | ------------------------------- |
| `SliceNotEmpty`              | At least one element            |
| `SliceLenMin(n)`             | At least n elements             |
| `SliceLenMax(n)`             | At most n elements              |
| `SliceLenBetween(low, high)` | Length in range [low, high]     |
| `SliceUnique`                | All elements are unique         |
| `SliceEach(v)`               | Apply validator to each element |

### Maps

| Validator                  | Description                      |
| -------------------------- | -------------------------------- |
| `MapNotEmpty`              | At least one entry               |
| `MapLenMin(n)`             | At least n entries               |
| `MapLenMax(n)`             | At most n entries                |
| `MapLenBetween(low, high)` | Entry count in range [low, high] |
| `MapKeys(v)`               | Apply validator to all keys      |
| `MapValues(v)`             | Apply validator to all values    |

### Network

| Validator    | Description                 |
| ------------ | --------------------------- |
| `NetIP`      | Valid IPv4 or IPv6 address  |
| `NetIPv4`    | Valid IPv4 address          |
| `NetIPv6`    | Valid IPv6 address          |
| `NetCIDR`    | Valid CIDR notation         |
| `NetHost`    | Valid hostname              |
| `NetPort`    | Valid port number (1-65535) |
| `NetPortStr` | Valid port number as string |

### Time

| Validator                 | Description                |
| ------------------------- | -------------------------- |
| `TimeFuture`              | Time is in the future      |
| `TimePast`                | Time is in the past        |
| `TimeAfter(t)`            | Time is after t            |
| `TimeBefore(t)`           | Time is before t           |
| `TimeBetween(start, end)` | Time in range [start, end] |

### Duration

| Validator               | Description                   |
| ----------------------- | ----------------------------- |
| `DurPositive`           | Duration > 0                  |
| `DurMin(d)`             | Duration >= d                 |
| `DurMax(d)`             | Duration <= d                 |
| `DurBetween(low, high)` | Duration in range [low, high] |

## Error Handling

Errors are returned as `FieldErrors`, a slice of `FieldError`:

```go
type FieldError struct {
    Field   string
    Message string
}
```

### Error Output

Single field error:

```
email: must be a valid email address
```

Multiple errors (from `RunAll` or `ValidateStruct`):

```
validation failed: name: must be between 2 and 100 characters; email: must be a valid email address
```

Nested struct errors use dot notation:

```
address.street: is required
address.zip: must contain only digits
```

Slice errors use bracket notation:

```
addresses[0].street: is required
emails[1]: must be a valid email address
```

### Accessing Individual Errors

```go
if err := user.Validate(); err != nil {
    if fieldErrs, ok := err.(criterio.FieldErrors); ok {
        for _, fe := range fieldErrs {
            fmt.Printf("%s: %s\n", fe.Field, fe.Message)
        }
    }
}
```

## Customization

String validators use exported regex patterns that can be modified globally:

```go
var (
    EmailRegex        *regexp.Regexp
    UUIDRegex         *regexp.Regexp
    AlphaRegex        *regexp.Regexp
    AlphanumericRegex *regexp.Regexp
    NumericRegex      *regexp.Regexp
    WhitespaceRegex   *regexp.Regexp
)
```

Override in an init function to change validation behavior:

```go
func init() {
    criterio.EmailRegex = regexp.MustCompile(`your-pattern`)
}
```

## Benchmarks

Comparison against manual validation and [go-playground/validator](https://github.com/go-playground/validator).

Run `task bench:readme` to regenerate.

**Key takeaways:**

- **Zero allocations** on valid inputs - the happy path is allocation-free
- **~2x faster** than go-playground/validator for struct validation
- Comparable to manual validation for simple cases, with the benefit of structured error handling
- Error paths allocate more than manual validation due to `FieldErrors` construction, but provide richer error context

<!-- BENCHMARK_START -->

| Benchmark                | ns/op | B/op | allocs/op |
| ------------------------ | ----: | ---: | --------: |
| Email_Manual_Valid       | 368.8 |    0 |         0 |
| Email_Playground_Valid   | 654.4 |    0 |         0 |
| Email_Criterio_Valid     | 374.1 |    0 |         0 |
| Email_Manual_Invalid     | 365.5 |   16 |         1 |
| Email_Playground_Invalid | 699.7 |  185 |         3 |
| Email_Criterio_Invalid   | 418.1 |   56 |         2 |
| User_Manual_Valid        | 233.2 |   88 |         5 |
| User_Playground_Valid    | 877.2 |   48 |         1 |
| User_Criterio_Valid      | 409.6 |    0 |         0 |
| User_Manual_Invalid      | 19.93 |   16 |         1 |
| User_Playground_Invalid  |  1145 |  664 |        11 |
| User_Criterio_Invalid    |  1205 |  584 |        17 |
| Order_Manual_Valid       | 540.4 |   88 |         5 |
| Order_Playground_Valid   |  1956 |  235 |         5 |
| Order_Criterio_Valid     | 816.9 |    0 |         0 |
| Order_Manual_Invalid     | 23.52 |   16 |         1 |
| Order_Playground_Invalid |  2129 | 2188 |        25 |
| Order_Criterio_Invalid   |  3451 | 2785 |        59 |
| Slice_Manual_Valid       | 4.161 |    0 |         0 |
| Slice_Playground_Valid   | 467.9 |  104 |         7 |
| Slice_Criterio_Valid     | 79.73 |   40 |         2 |
| User_Manual_Parallel     | 72.41 |   88 |         5 |
| User_Playground_Parallel | 180.9 |   49 |         1 |
| User_Criterio_Parallel   | 69.96 |    0 |         0 |

<!-- BENCHMARK_END -->

## License

MIT
