package pass

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/jgsheppa/localPass/generator"
)

type Pass struct {
	URL      string
	password string
	reader   *bufio.Reader
	output   io.Writer
}

func (p *Pass) EnterURL() error {
	fmt.Fprint(p.output, "Enter a URL: ")
	url, err := p.reader.ReadString('\n')
	if err != nil {
		return err
	}

	trimURL := strings.TrimSpace(url)
	p.URL = trimURL

	return nil
}

func (p *Pass) EnterPassword() error {
	fmt.Fprint(p.output, "Generate password? (Y/n): ")
	answer, err := p.reader.ReadString('\n')
	if err != nil {
		return err
	}

	trimAnswer := strings.TrimSpace(answer)

	if trimAnswer == "Y" || trimAnswer == "y" {
		conditions := generator.NewPasswordConditions()
		pw, err := conditions.GeneratePassword()
		if err != nil {
			return err
		}
		p.password = pw
		return nil
	}

	fmt.Fprint(p.output, "Enter your password: ")
	userPassword, err := p.reader.ReadString('\n')
	if err != nil {
		return err
	}
	p.password = userPassword

	return nil
}

func CreatePass(output io.Writer, rd io.Reader) (Pass, error) {
	var pass Pass

	// TODO: create pass generator with default values
	pass.reader = bufio.NewReader(rd)
	pass.output = output
	pass.EnterURL()
	pass.EnterPassword()

	return pass, nil
}
