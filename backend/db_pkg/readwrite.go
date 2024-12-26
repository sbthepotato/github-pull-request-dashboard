package db_pkg

import (
	"errors"
	"io"
	"os"
)

/*
Write a json to file
*/
func Write(fileName string, jsonData []byte) error {

	file, err := os.Create("db_files/" + fileName + ".json")
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

/*
Read a json from file
*/
func Read(fileName string) ([]byte, error) {
	file, err := os.Open("db_files" + fileName + ".json")

	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
