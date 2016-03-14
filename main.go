package main

import (
	"fmt"
	"log"

	"github.com/namely/broadway/instance"
	"github.com/namely/broadway/playbook"
)

func main() {
	/*
		args := os.Args
		yamlFileDescriptor := args[1:][0]
	*/
	playbooks, err := playbook.LoadPlaybookFolder("playbooks/")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v+\n", playbooks)
	fmt.Println(instance.InstanceStatusNew)
}
