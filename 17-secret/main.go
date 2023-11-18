package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aboelkassem/gophercises/secret/secret"
)

func main() {
	var (
		flagEncodingKey = flag.String("key", "", "encoding key")
		flagPath        = flag.String("path", "vault.enc", "path to vault file")
	)
	flag.Parse()

	v := secret.FileVault(*flagEncodingKey, *flagPath)

	// read package arguments
	switch cmd := flag.Arg(0); cmd {
	case "set":
		keyName, keyValue := flag.Arg(1), flag.Arg(2)
		if keyName == "" || keyValue == "" {
			log.Fatal("Key and value must be provided")
		}
		err := v.Set(keyName, keyValue)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Value set!")
	case "get":
		keyName := flag.Arg(1)
		if keyName == "" {
			log.Fatal("Key must be provided")
		}
		value, err := v.Get(keyName)
		if err != nil {
			// if err == secret.ErrKeyNotFound {
			// 	log.Fatalf("Key %q not found", keyName)
			// }
			log.Fatalln(err)
		}
		fmt.Println(value)
	case "list":
		keyNames, err := v.List()
		if err != nil {
			log.Fatalln(err)
		}

		if keyNames == nil {
			fmt.Println("Empty vault")
			return
		}
		for _, name := range keyNames {
			fmt.Println(name)
		}
	case "delete":
		keyName := flag.Arg(1)
		if keyName == "" {
			log.Fatal("Key and value must be provided")
		}
		err := v.Delete(keyName)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Key deleted!")
	default:
		log.Fatalf("Unknown command %q, please use set or get.", cmd)
	}

}
