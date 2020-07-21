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

## License
[MIT](https://choosealicense.com/licenses/mit/)
