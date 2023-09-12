package main

import (
	"display-img-iterm/imgcat"
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

	wc := imgcat.NewWriter(os.Stdout)
	_, err = io.Copy(wc, f)
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}
	return nil

}
