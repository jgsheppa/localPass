package list

import (
	"fmt"
	"io"
	"os"

	"github.com/jgsheppa/localPass/models"
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

func Run(service models.PassService) (int, error) {
	var urls []string
	passes, err := service.Get()
	if err != nil {
		return 1, err
	}
	for _, pass := range passes {
		urls = append(urls, pass.URL)
	}

	// TODO: remove these sites once database is integrated
	err = NewList().WithSites(urls).PrintList()
	if err != nil {
		return 1, err
	}

	return 0, nil
}
