package csv2json

//write test to Convert method

import (
	"encoding/csv"
	"net/http"
	"os"
	"testing"
)

func TestConvert(t *testing.T) {
	url := "https://media.githubusercontent.com/media/datablist/sample-csv-files/main/files/organizations/organizations-100.csv"
	csv2json := New()

	_, err := csv2json.Convert(url, "test.json")
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Stat("test.json"); os.IsNotExist(err) {
		t.Error("File was not created")
	}
	err = os.Remove("test.json")
	if err != nil {
		t.Error(err)
	}

	httpClient := http.Client{}
	resp, err := httpClient.Get(url)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)
	//save csv file
	file, err := os.Create("test.csv")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	records, err := csvReader.ReadAll()
	if err != nil {
		t.Error(err)
	}

	err = writer.WriteAll(records)
	if err != nil {
		t.Error(err)
	}

	_, err = csv2json.Convert("test.csv", "test.json")
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Stat("test.json"); os.IsNotExist(err) {
		t.Error("File was not created")
	}
	err = os.Remove("test.json")
	if err != nil {
		t.Error(err)
	}

	err = os.Remove("test.csv")
	if err != nil {
		t.Error(err)
	}
}
