package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		// .env file not found is not an error
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error opening .env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set the environment variable
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting env var %s: %w", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	return nil
}
