package imgcat

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"strings"
)

type writer struct {
	pw   *io.PipeWriter
	done chan struct{}
}

func (w *writer) Write(data []byte) (int, error) {
	return w.pw.Write(data)

}

func (w *writer) Close() error {
	if err := w.pw.Close(); err != nil {
		return err
	}
	<-w.done
	return nil
}
func NewWriter(w io.Writer) *writer {
	pr, pw := io.Pipe()
	custW := &writer{pw: pw, done: make(chan struct{})}
	go func() {
		err := Copy(w, pr)
		if err != nil {
			log.Println(err.Error())
		}
		custW.done <- struct{}{}
	}()
	return custW
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
