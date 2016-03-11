package main

import (
	"fmt"
	"log"
	"os"

	"github.com/namely/broadway/instance"
	"github.com/namely/broadway/playbook"
)

func main() {
	args := os.Args
	yamlFileDescriptor := args[1:][0]

	playbookBytes, err := playbook.ReadPlaybookFromDisk(yamlFileDescriptor)
	if err != nil {
		log.Fatal(err)
	}

	pb, err := playbook.ParsePlaybook(playbookBytes)
	if err != nil {
		log.Fatal(err)
	}

	if err := pb.ValidateTasks(); err != nil {
		log.Fatalf("Task validation failed: %s", err)
	}

	fmt.Println(instance.InstanceStatusNew)
	fmt.Printf("%+v", pb)
}
