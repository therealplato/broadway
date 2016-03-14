package main

import (
	"fmt"
	"github.com/namely/broadway/instance"
	"github.com/namely/broadway/playbook"
	"log"
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
