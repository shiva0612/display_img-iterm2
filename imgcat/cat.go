package imgcat

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

func Copy(w io.Writer, r io.Reader) error {
	_, err := fmt.Fprintf(w, "\033]1337;File=inline=1:")
	if err != nil {
		return err
	}

	wc := base64.NewEncoder(base64.StdEncoding, w)
	_, err = io.Copy(wc, r)
	if err != nil {
		return errors.New("could not bas64 encode: " + err.Error())
	}
	err = wc.Close()
	if err != nil {
		return errors.New("could not close the writeCloser" + err.Error())
	}

	_, err = fmt.Fprintf(w, "\a")
	if err != nil {
		return err
	}

	return nil
}
