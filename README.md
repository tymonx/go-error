# Go Error

The Go Error library implements runtime error with message string formatted
using **replacement fields** surrounded by curly braces `{}` format strings from
the [Go Formatter](https://gitlab.com/tymonx/go-formatter) library.

[[_TOC_]]

## Features

* Format string by providing arguments without using placeholders or format verbs `%`
* Format string using automatic placeholder `{p}`
* Format string using positional placeholders `{pN}`
* Format string using named placeholders `{name}`
* Format string using object placeholders `{.Field}`, `{p.Field}` and `{pN.Field}` where `Field` is an exported `struct` field or method
* Set custom format error message string. Default is `{.FileBase}:{.Line}:{.FunctionBase}(): {.String}`
* Error message contains file path, line number, function name from where was called
* Compatible with the standard `errors` package with `As`, `Is` and `Unwrap` functions
* It uses the [Go Formatter](https://gitlab.com/tymonx/go-formatter) library

## Usage

Import `rterror` package:

```go
import "gitlab.com/tymonx/go-error/rterror"
```

### Without arguments

```go
err := rterror.New("Error message")

fmt.Println(err)
```

Output:

```plaintext
<file>:<line>:<function>(): Error message
```

### With arguments

```go
err := rterror.New("Error message {p1} -", 3, "bar")

fmt.Println(err)
```

Output:

```plaintext
<file>:<line>:<function>(): Error message bar - 3
```

### Wrapped

```go
wrapped := rterror.New("Wrapped error")

err := rterror.New("Error message {p1} -", 3, err, "bar")

fmt.Println(errors.Is(err, wrapped))
```

Output:

```plaintext
true
```

### Custom format

```go
err := rterror.New("Error message {p1} -", 3, "bar").SetFormat("#{.Function} := '{.String}' <-")

fmt.Println(err)
```

Output:

```plaintext
#<function> := 'Error message bar - 3' <-
```

### Custom error type

```go
type MyError struct {
    rterror.RuntimeError
}

func New(message string, arguments ...interface{}) *MyError {
    return &MyError{
        RuntimeError: *rterror.NewSkipCaller(rterror.SkipCall, message, arguments...),
    }
}

err := New("My custom error")

fmt.Println(err)
```

Output:

```plaintext
<file>:<line>:<function>(): My custom error
```
