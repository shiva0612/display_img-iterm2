package imgcat

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"strings"
)

func NewWriter(w io.Writer, done chan struct{}) io.WriteCloser {
	pr, pw := io.Pipe()
	go func() {
		err := Copy(w, pr) //copying raw imgContent from pr -> os.Stdout (header+base64(body)+footer)
		if err != nil {
			log.Println(err.Error())
		}
		done <- struct{}{}
	}()
	return pw
}

// reads content from reader -> writes(header + base64[content] + footer)
func Copy(w io.Writer, r io.Reader) error {

	header := strings.NewReader("\033]1337;File=inline=1:")
	footer := strings.NewReader("\a")

	bodyr, bodyw := io.Pipe()
	go func() {
		defer bodyw.Close()
		wc := base64.NewEncoder(base64.StdEncoding, bodyw)
		_, err := io.Copy(wc, r)
		if err != nil {
			bodyr.CloseWithError(fmt.Errorf("error copying file to encoder: %s", err.Error()))
		}
		err = wc.Close()
		if err != nil {
			bodyr.CloseWithError(fmt.Errorf("error flushing encoded bytes to writer: %s", err.Error()))
		}
	}()

	_, err := io.Copy(w, io.MultiReader(header, bodyr, footer))
	if err != nil {
		return err
	}

	return nil
}
