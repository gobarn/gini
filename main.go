package main

import (
	"strings"
	"bufio"
	"io"
	"fmt"
	"log"
	"os"
)

type Section struct {
	Name string
	Values map[string]string
}

type Gini struct {
	Path string
	Sections map[string]*Section
}

func (gini *Gini) Parse(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	var section *Section
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// skip empty lines
		if line == "" {
			continue
		}

		// scan comments
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		// scan section header
		if strings.HasPrefix(line, "[") {
			sectionName := line[1:len(line)-1]

			if sec, ok := gini.Sections[sectionName]; !ok {
				section = &Section{
					Name: sectionName,
					Values: make(map[string]string),
				}
				gini.Sections[section.Name] = section
			} else {
				section = sec
			}

			continue
		}

		idx := strings.Index(line, "=")

		key := strings.ToLower(strings.TrimSpace(line[0:idx]))
		value := strings.ToLower(strings.TrimSpace(line[idx+1:]))
		value = strings.Trim(value, "\"'")
		section.Values[key] = value
	}
}

func (gini *Gini) Get(section string, key string) string {
	return gini.Sections[section].Values[key]
}

func ParseFile(path string) *Gini {
	gini := &Gini{
		Path: path,
		Sections: make(map[string]*Section),
	}

	file, err := os.Open(gini.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	gini.Parse(file)

	return gini
}


func main() {
	ini := ParseFile("test.ini")
	fmt.Println(ini.Get("mysql", "user"))
	fmt.Println(ini.Get("mysql", "pass"))
}
