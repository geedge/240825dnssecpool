package internal

import (
	"bufio"
	"encoding/csv"
	"os"
	"strings"
)

func ReadResolversTXT(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var resolvers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		resolver := strings.TrimSpace(scanner.Text())
		if resolver != "" {
			resolvers = append(resolvers, resolver)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return resolvers, nil
}

func ReadResolversCSV(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var resolvers []string
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		if len(record) > 0 {
			resolver := strings.TrimSpace(record[0])
			if resolver != "" {
				resolvers = append(resolvers, resolver)
			}
		}
	}

	return resolvers, nil
}
