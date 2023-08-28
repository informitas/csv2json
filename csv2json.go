package csv2json

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type csv2json struct {
	mu *sync.Mutex
}

// New initializes a new csv2json instance.
func New() *csv2json {
	return &csv2json{&sync.Mutex{}}
}

// Convert converts CSV data from various sources to JSON format.
// src can be a file path, a URL, or a [][]string.
// dest is the file path where the JSON data will be saved.
// If dest is empty, the JSON data will not be saved to a file. Instead, it will be returned as a []map[string]interface{}.
func (c *csv2json) Convert(src interface{}, dest string) ([]map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if reflect.TypeOf(src).Kind() == reflect.Slice {
		if reflect.TypeOf(src).Elem().Kind() == reflect.Slice {
			if reflect.TypeOf(src).Elem().Elem().Kind() == reflect.String {
				//src is a [][]string so we can convert it to a map[string]interface{}
				result, err := c.transformToMap(src.([][]string))
				if err != nil {
					return nil, err
				}
				if dest == "" {
					return result, nil
				}
				//save to file
				err = c.saveToFile(result, dest)
				if err != nil {
					return nil, err
				}
				return nil, nil
			}
		}
	}

	//check if src is a url
	if reflect.TypeOf(src).Kind() == reflect.String {
		if strings.HasPrefix(src.(string), "http") || strings.HasPrefix(src.(string), "https") {
			resp, err := c.downloadFromURL(src.(string))
			if err != nil {
				return nil, err
			}
			result, err := c.transformToMap(resp)
			if err != nil {
				return nil, err
			}
			if dest == "" {
				return result, nil
			}

			//save to file
			err = c.saveToFile(result, dest)
			if err != nil {
				return nil, err
			}
			return nil, nil
		}
	}

	if reflect.TypeOf(src).Kind() == reflect.String {
		// src is a file path
		file, err := os.Open(src.(string))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		csvReader := csv.NewReader(file)
		data, err := csvReader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		result, err := c.transformToMap(data)
		if err != nil {
			return nil, err
		}
		if dest == "" {
			return result, nil
		}
		//save to file
		err = c.saveToFile(result, dest)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	return nil, fmt.Errorf("invalid src type")
}

// downloadFromURL downloads CSV data from a URL.
// It returns a [][]string.
func (c *csv2json) downloadFromURL(url string) ([][]string, error) {
	httpClient := http.Client{}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return data, nil
}

// transformToMap converts [][]string to []map[string]interface{}.
// It returns a []map[string]interface{}.
func (c *csv2json) transformToMap(src [][]string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	headers := src[0]
	for _, line := range src[1:] {
		result := map[string]interface{}{}
		for i, field := range line {
			header := headers[i]

			nestedObject := strings.Split(header, ".")
			internal := result

			for index, val := range nestedObject {
				key, arrayIndex := c.arrayContentMatch(val)
				if arrayIndex != -1 {
					if internal[key] == nil {
						internal[key] = []interface{}{}
					}
					internalArray := internal[key].([]interface{})
					if index == len(nestedObject)-1 {
						if _, err := strconv.Atoi(field); err == nil {
							intField, _ := strconv.Atoi(field)
							internalArray = append(internalArray, intField)
							internal[key] = internalArray
							break
						}
						if field == "true" || field == "false" {
							boolField, _ := strconv.ParseBool(field)
							internalArray = append(internalArray, boolField)
							internal[key] = internalArray
							break
						}

						internalArray = append(internalArray, field)
						internal[key] = internalArray
						break
					}
					if arrayIndex >= len(internalArray) {
						internalArray = append(internalArray, map[string]interface{}{})
					}
					internal[key] = internalArray
					internal = internalArray[arrayIndex].(map[string]interface{})
				} else {
					if index == len(nestedObject)-1 {
						if _, err := strconv.Atoi(field); err == nil {
							internal[key], _ = strconv.Atoi(field)
							break
						}
						if field == "true" || field == "false" {
							internal[key], _ = strconv.ParseBool(field)
							break
						}

						internal[key] = field
						break
					}
					if internal[key] == nil {
						internal[key] = map[string]interface{}{}
					}
					internal = internal[key].(map[string]interface{})
				}
			}
		}
		results = append(results, result)
	}

	return results, nil
}

// saveToFile saves JSON data to a file.
func (c *csv2json) saveToFile(data []map[string]interface{}, dest string) error {
	bytes, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//save to file
	err = os.WriteFile(dest, bytes, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// This function checks if a string contains an array index.
func (c *csv2json) arrayContentMatch(str string) (string, int) {
	i := strings.Index(str, "[")
	if i >= 0 {
		j := strings.Index(str, "]")
		if j >= 0 {
			index, _ := strconv.Atoi(str[i+1 : j])
			return str[0:i], index
		}
	}
	return str, -1
}
