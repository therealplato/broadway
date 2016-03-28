package main

import (
	"fmt"
	"log"
	"os"

	"github.com/namely/broadway/instance"
	"github.com/namely/broadway/playbook"
	"github.com/namely/broadway/server"
	"github.com/namely/broadway/services"
	"github.com/namely/broadway/store"
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

	MS := services.NewManifestService()
	manifests, err := MS.LoadManifestFolder()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", playbooks)
	fmt.Printf("%+v\n", manifests)
	fmt.Println(instance.StatusNew)

	server := server.New(store.New())
	server.SetPlaybooks(playbooks)
	err = server.Run(os.Getenv("HOST"))
	if err != nil {
		panic(err)
	}

}
