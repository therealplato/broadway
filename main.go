package main

import (
	"fmt"
	"github.com/namely/broadway/instance"
	"github.com/namely/broadway/playbook"
	"log"
	"path/filepath"
)

func main() {
	/*
		args := os.Args
		yamlFileDescriptor := args[1:][0]
	*/
	playbooks, err := LoadPlaybookFolder("playbooks/")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v+\n", playbooks)
	fmt.Println(instance.InstanceStatusNew)
}

func LoadPlaybookFolder(dir string) ([]playbook.Playbook, error) {
	var AllPlaybooks []playbook.Playbook
	paths, err := filepath.Glob(dir + "/*")
	if err != nil {
		return AllPlaybooks, err
	}
	for _, path := range paths {
		playbookBytes, err := playbook.ReadPlaybookFromDisk(path)
		if err != nil {
			fmt.Printf("Warning: Failed to read %s\n", path)
			continue
		}
		parsed, err := playbook.ParsePlaybook(playbookBytes)
		if err != nil {
			fmt.Printf("Warning: Failed to parse %s\n", path)
			continue
		}
		err = parsed.Validate()
		if err != nil {
			fmt.Printf("Warning: Playbook %s invalid: %s\n", path, err)
			continue
		}
		AllPlaybooks = append(AllPlaybooks, parsed)
	}
	return AllPlaybooks, nil
}
