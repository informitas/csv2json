# csv2json Go Package

The `csv2json` Go package provides a convenient way to convert CSV data from various sources into JSON format. It supports conversion from file paths, URLs, and `[][]string` slices to JSON. The resulting JSON data can be either saved to a file or returned as a `[]map[string]interface{}` slice.

## Installation

To use the `csv2json` package in your Go project, you can install it using the following command:

```shell
go get github.com/informitas/csv2json
```

## Usage

Import the `csv2json` package into your Go code and follow the usage instructions below.

```go
import (
	"fmt"
	"github.com/informitas/csv2json"
)

func main() {
	csvConverter := csv2json.New()

	// Convert from a file path to JSON
	jsonData, err := csvConverter.Convert("input.csv", "output.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("JSON data:", jsonData)

	// Convert from a URL to JSON
	jsonData, err = csvConverter.Convert("https://example.com/data.csv", "output.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("JSON data:", jsonData)

	// Convert from [][]string to JSON
	csvData := [][]string{
		{"Name", "Age"},
		{"Alice", "25"},
		{"Bob", "30"},
	}
	jsonData, err = csvConverter.Convert(csvData, "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("JSON data:", jsonData)
}
```

## Methods

### New

The `New` function initializes a new `csv2json` instance.

```go
func New() *csv2json
```

### Convert

The `Convert` method converts CSV data from various sources to JSON format.

```go
func (c *csv2json) Convert(src interface{}, dest string) ([]map[string]interface{}, error)
```

- `src`: The source data for conversion. It can be a file path, a URL, or a `[][]string` slice.
- `dest`: The file path where the JSON data will be saved. If empty, the JSON data will be returned as a `[]map[string]interface{}`.

### Other Internal Methods

- `downloadFromURL(url string) ([][]string, error)`: Downloads CSV data from a URL and returns it as a `[][]string`.
- `transformToMap(src [][]string) ([]map[string]interface{}, error)`: Converts `[][]string` to `[]map[string]interface{}`.
- `saveToFile(data []map[string]interface{}, dest string) error`: Saves JSON data to a file.
- `arrayContentMatch(str string) (string, int)`: Checks if a string contains an array index and returns the modified string and the array index.

## License

This project is licensed under the [MIT License](LICENSE).
