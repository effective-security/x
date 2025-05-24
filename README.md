# x

`x` is a collection of small, focused Go utility packages grouped into a single Go module.

## Supported Packages

| Package        | Description                                             |
| -------------- | ------------------------------------------------------- |
| configloader   | Load JSON/YAML config with environment overrides        |
| ctl            | Kong-based CLI helpers (version flag, bool-pointer)     |
| enum           | Generic enum and bit-flag helpers (Go 1.18+ generics)   |
| fileutil       | File and folder utilities, live file reloader, resolver |
| flake          | Sonyflake-inspired unique ID generator                  |
| format         | String, number, boolean, and time formatting utilities  |
| guid           | Simple GUID/UUID generator                              |
| netutil        | Free port discovery, local IP, node info, URL parsing   |
| print          | JSON/YAML/table/text printers                           |
| slices         | Common slice operations (equal, contains, unique, etc.) |
| ticker         | Context-aware ticker with callback and status tracking  |
| urlutil        | URL and querystring helpers, public endpoint builder    |
| values         | Generic any-to-typed conversion and MapAny utilities    |

## Installation

```bash
go get github.com/effective-security/x
```

## Usage

Import the package(s) you need in your Go code. For example, to load configuration:

```go
import "github.com/effective-security/x/configloader"

func main() {
    factory, err := configloader.NewFactory(nil, nil, "MYAPP_")
    if err != nil {
        log.Fatal(err)
    }
    var cfg Config
    if err := factory.Load("config.yaml", &cfg); err != nil {
        log.Fatal(err)
    }
    // Use cfg...
}
```

And to find a free TCP port:

```go
import "github.com/effective-security/x/netutil"

port, err := netutil.FindFreePort("localhost", 10)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Listening on port %d\n", port)
```

See individual package directories for more examples and detailed documentation.

## Development

Requires Go 1.24 or newer and GNU Make.

```bash
# List available make targets
make help

# Format, lint, generate code, and run tests with coverage
make all
```

The current release version is managed in the `.VERSION` file and published via GitHub Actions.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
