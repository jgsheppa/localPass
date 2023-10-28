package list

import (
	"fmt"
	"io"
	"os"
)

type PassList struct {
	Sites  []string
	writer io.Writer
}

func NewList() *PassList {
	return &PassList{writer: os.Stdout}
}

func (pl *PassList) WithOutput(output io.Writer) *PassList {
	pl.writer = output
	return pl
}

func (pl *PassList) WithSites(sites []string) *PassList {
	pl.Sites = sites
	return pl
}

func (pl *PassList) PrintList() error {
	for _, site := range pl.Sites {
		_, err := fmt.Fprintf(pl.writer, "%s \n", site)
		if err != nil {
			return err
		}
	}
	return nil
}

func Run() (int, error) {
	// TODO: remove these sites once database is integrated
	err := NewList().WithSites([]string{"www.acme.xyz"}).PrintList()
	if err != nil {
		return 1, err
	}

	return 0, nil
}
