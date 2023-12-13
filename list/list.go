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
	Passes []models.Pass
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

func (pl *PassList) WithPasses(passes []models.Pass) *PassList {
	pl.Passes = passes
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

func (pl *PassList) PrintPasses() error {
	for _, pass := range pl.Passes {
		_, err := fmt.Fprintf(pl.writer, "%d | %s | %s \n", pass.ID, pass.URL, pass.Password)
		if err != nil {
			return err
		}
	}
	return nil
}

func Run(service models.PassService) (int, error) {
	passes, err := service.Get()
	if err != nil {
		return 1, err
	}

	err = NewList().WithPasses(passes).PrintPasses()
	if err != nil {
		return 1, err
	}

	return 0, nil
}
