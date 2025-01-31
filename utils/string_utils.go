package utils

import (
	"encoding/json"
	"fmt"
	"handbook-scraper/utils/log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// RemoveHTMLTags removes HTML tags from a string, replacing <br> tags with \n and removes others.
func RemoveHTMLTags(s string) string {
	// Replace <br /> or <br/> tags with newlines
	re := regexp.MustCompile(`<br\s*/?>`)
	s = re.ReplaceAllString(s, "\n")

	// Remove all other HTML tags
	re = regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(s, "")
}

// StringToInt extracts the first integer from a string
func StringToInt(s string) int {
	re := regexp.MustCompile("[0-9]+")
	numStr := re.FindString(s)
	num, _ := strconv.Atoi(numStr)
	return num
}

// StringToArray converts a string to an array of strings
func StringToArray(s string) []string {
	// Trim whitespace and split by newlines
	lines := strings.Split(strings.TrimSpace(s), "\n")
	return lines
}

// Jsonify converts anything to a JSON string
func Jsonify(any interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(any, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error converting data to JSON")
	}
	return string(jsonData), nil
}

// SaveJSONToFile saves a JSON string to a file
func SaveJSONToFile(jsonData string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("error closing file: %v\n", err)
		}
	}(file)

	_, err = file.WriteString(jsonData)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %w", err)
	}
	return nil
}

// SaveDataToFile saves data to a file
func SaveDataToFile(data []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing data to file: %w", err)
	}
	return nil
}

// LoadDataFromFile loads data from a file
func LoadDataFromFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	return data, nil
}
