![build](https://github.com/mottaquikarim/go-airtable/workflows/Build%20Status/badge.svg)
# [Go Airtable](https://godoc.org/github.com/mottaquikarim/go-airtable)
Simple Airtable client written in go

## Installation

```bash
go get github.com/mottaquikarim/go-airtable
```

## Usage

```go
import (
  "github.com/mottaquikarim/go-airtable"
)

acc := airtable.Account{
  ApiKey: "XXXX",
  BaseId: "XXXX",
}
pokédex := airtable.NewTable("pokémon", acc)
original_generation, err := pokédex.List(airtable.Options{})
if err != nil {
  // handle error
}
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

### Run Tests

```bash
make test-dev
```

### Format

```bash
make fmt
```

### Lint

```bash
make lint
```

### CLI

```bash
make build
```

Then,

```bash
make run arguments="-help"
```

Output:
```
Usage: ./airtbl [flags]
  -api-key="XXXXX": Airtable API Key
  -base-id="XXXXX": Airtable Base Id
```

## License
[MIT](https://choosealicense.com/licenses/mit/)
