package gini

import (
	"fmt"
	"os"
	"io"
	"bufio"
	"strings"
)

const (
	DEFAULT_SECTION = "default"
	DEFAULT_COMMENT = "#"
	DEFAULT_COMMENT_SEM = ";"
)

type Pair map[string]string
type Section map[string]Pair

type Config struct {
	filepath string
	sections Section
}

func Parse(filepath string) (*Config, error) {
	config := &Config{
		filepath: filepath,
		sections: make(Section),
	}

	err := config.ParseFile()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (config *Config) ParseFile() error {
	file, err := os.Open(config.filepath)

	if err != nil {
		return err
	}

	defer file.Close()

	err = config.Parse(file)
	if err != nil {
		return err
	}

	return nil
}

func (config *Config) Parse(src io.Reader) error {
	var section string
	scanner := bufio.NewScanner(src)
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, DEFAULT_COMMENT):
			continue
		case strings.HasPrefix(line, DEFAULT_COMMENT_SEM):
			continue
		case strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]"):
			section = line[1:len(line)-1]
		default:
			item := strings.SplitN(line, "=", 2)
			if len(item) != 2 {
				return fmt.Errorf("parse error: %s:%d %s", config.filepath, lineNum, item[0])
			}

			key := strings.TrimSpace(item[0])
			value := strings.TrimSpace(item[1])
			config.Add(section, key, value)
		}

		lineNum++
	}

	return nil
}

func (config *Config) Add(section string, key string, value string) bool {
	if section == "" {
		section = DEFAULT_SECTION
	}

	if _, ok := config.sections[section]; !ok {
		config.sections[section] = make(Pair)
	}

	_, ok := config.sections[section][key]
	config.sections[section][key] = value;

	return !ok
}

func (config *Config) Get(section string, key string) string {
	if section == "" {
		section = DEFAULT_SECTION
	}

	return config.sections[section][key]
}



