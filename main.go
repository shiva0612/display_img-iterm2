package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("path is not provided...")
	}

	for _, path := range os.Args[1:] {
		if err := cat(path); err != nil {
			fmt.Fprintf(os.Stderr, "could load img: error: %v\n", err)
		}
	}
}

func cat(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.New("error opening the file: " + path + " error =  " + err.Error())
	}
	defer f.Close()

	fmt.Printf("\033]1337;File=inline=1:")

	wc := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	_, err = io.Copy(wc, f)
	if err != nil {
		return errors.New("could not bas64 encode: " + err.Error())
	}
	err = wc.Close()
	if err != nil {
		return errors.New("could not close the writeCloser" + err.Error())
	}

	fmt.Printf("\a")

	return nil
}
