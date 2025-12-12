# criterio

A composable validation library for Go.

[Go Reference](https://pkg.go.dev/github.com/hay-kot/criterio)

## Install

```bash
go get -u github.com/hay-kot/criterio
```

## Usage

```go
// Inline validation
err := criterio.Run("email", userInput,
    criterio.Required[string](),
    criterio.StrEmail(),
)

// Reusable validator
validateName := criterio.New("name",
    criterio.Required[string](),
    criterio.StrBetween(1, 100),
)

if err := validateName(name); err != nil {
    return err
}
```

## Validators

### Core

| Validator              | Description             | Status |
| ---------------------- | ----------------------- | ------ |
| `Required[T]()`        | Value is non-zero       | ✓      |
| `OneOf[T](allowed...)` | Value is in allowed set | ✓      |

### Ordered Types (integers, floats, strings)

| Validator            | Description                | Status |
| -------------------- | -------------------------- | ------ |
| `Min(min)`           | Value >= min               | ✓      |
| `Max(max)`           | Value <= max               | ✓      |
| `Between(low, high)` | Value in range [low, high] | ✓      |
| `Positive()`         | Value > 0                  | ✓      |
| `Negative()`         | Value < 0                  | ✓      |
| `NonZero()`          | Value != 0                 | ✓      |
| `MultipleOf(n)`      | Value is divisible by n    | ✓      |

### Strings

| Validator               | Description                         | Status |
| ----------------------- | ----------------------------------- | ------ |
| `StrNotEmpty()`         | Non-empty after trimming whitespace | ✓      |
| `StrMin(min)`           | Length >= min                       | ✓      |
| `StrMax(max)`           | Length <= max                       | ✓      |
| `StrBetween(low, high)` | Length in range [low, high]         | ✓      |
| `StrMatches(pattern)`   | Matches regex pattern               | ✓      |
| `StrEmail()`            | Valid email format                  | ✓      |
| `StrOneOf(allowed...)`  | Value in allowed set                | ✓      |
| `StrURL()`              | Valid URL format                    | ✓      |
| `StrUUID()`             | Valid UUID format                   | ✓      |
| `StrAlpha()`            | Letters only                        | ✓      |
| `StrAlphanumeric()`     | Letters and numbers only            | ✓      |
| `StrNumeric()`          | Digits only                         | ✓      |
| `StrContains(substr)`   | Contains substring                  | ✓      |
| `StrHasPrefix(prefix)`  | Starts with prefix                  | ✓      |
| `StrHasSuffix(suffix)`  | Ends with suffix                    | ✓      |
| `StrNoWhitespace()`     | No whitespace characters            | ✓      |
| `StrLowercase()`        | All lowercase                       | ✓      |
| `StrUppercase()`        | All uppercase                       | ✓      |

### Slices

| Validator                    | Description                   | Status |
| ---------------------------- | ----------------------------- | ------ |
| `SliceLenMin(min)`           | Length >= min                 | ✓      |
| `SliceLenMax(max)`           | Length <= max                 | ✓      |
| `SliceLenBetween(low, high)` | Length in range [low, high]   | ✓      |
| `SliceNotEmpty()`            | Length > 0                    | ✓      |
| `SliceUnique()`              | All elements unique           | ✓      |
| `SliceEach(validator)`       | Each element passes validator | ✓      |

### Maps

| Validator                  | Description                 | Status |
| -------------------------- | --------------------------- | ------ |
| `MapLenMin(min)`           | Length >= min               | ✓      |
| `MapLenMax(max)`           | Length <= max               | ✓      |
| `MapLenBetween(low, high)` | Length in range [low, high] | ✓      |
| `MapNotEmpty()`            | Length > 0                  | ✓      |
| `MapKeys(validator)`       | All keys pass validator     | ✓      |
| `MapValues(validator)`     | All values pass validator   | ✓      |

### Network

| Validator   | Description                 | Status |
| ----------- | --------------------------- | ------ |
| `NetIP()`   | Valid IP address (v4 or v6) | ✓      |
| `NetIPv4()` | Valid IPv4 address          | ✓      |
| `NetIPv6()` | Valid IPv6 address          | ✓      |
| `NetCIDR()` | Valid CIDR notation         | ✓      |
| `NetHost()` | Valid hostname              | ✓      |
| `NetPort()` | Valid port number (1-65535) | ✓      |

### Time

| Validator                 | Description                | Status |
| ------------------------- | -------------------------- | ------ |
| `TimeFuture()`            | Time is in the future      | ✓      |
| `TimePast()`              | Time is in the past        | ✓      |
| `TimeAfter(t)`            | Time is after t            | ✓      |
| `TimeBefore(t)`           | Time is before t           | ✓      |
| `TimeBetween(start, end)` | Time in range [start, end] | ✓      |

### Duration

| Validator               | Description                   | Status |
| ----------------------- | ----------------------------- | ------ |
| `DurMin(min)`           | Duration >= min               | ✓      |
| `DurMax(max)`           | Duration <= max               | ✓      |
| `DurBetween(low, high)` | Duration in range [low, high] | ✓      |
| `DurPositive()`         | Duration > 0                  | ✓      |

## Error Handling

Validators return `FieldErrors` which provides structured access to validation failures:

```go
err := criterio.Run("age", -5, criterio.Min(0))
if err != nil {
    // err.Error() returns "age: must be at least 0"
}
```

Multiple validators collect all errors:

```go
validateUser := criterio.New("username",
    criterio.Required[string](),
    criterio.StrBetween(3, 20),
    criterio.StrAlphanumeric(),
)
// Returns all failing validations, not just the first
```
