package prompt

import (
	"dcs/scraper"
	"errors"
	"os"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
)

func Drama(dramas []scraper.DramaInfo) (*scraper.DramaInfo, error) {
	type DramaItem struct {
		Name     string
		Link     string
		Episodes string
		Desc     string
		Info     *scraper.DramaInfo
	}

	items := make([]DramaItem, len(dramas))
	for i := 0; i < len(dramas); i++ {
		d := &dramas[i]
		items[i] = DramaItem{
			Name: d.Name,
			Link: d.FullURL,
			Info: d,
		}
	}

	// how the prompt should be displayed
	templates := &promptui.SelectTemplates{
		Label:    "{{ . | white | bold }}",
		Active:   "\U0001F449 {{ .Name | green | bold }}",
		Inactive: "  {{ .Name | red }}",
		Selected: "\U0001F449 {{ .Name | green }}",
		Details: `

{{ .Name | blue | bold}}
{{ "----------------------------" | white }}
{{ "Link:" | faint }}	{{ .Link | yellow }}`, // TODO: add more properties
	}

	// for using the SEARCH feature in the prompt
	searcher := func(input string, index int) bool {
		item := items[index]
		properties := []string{
			item.Name,
		}

		return len(fuzzy.FindNormalizedFold(input, properties)) > 0
	}

	// settings for the prompt
	prompt := promptui.Select{
		Label:     "Select a Drama:",
		Items:     items,
		Templates: templates,
		Size:      9,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	return items[i].Info, err
}

// DirSelect - Prompt for choosing a file/folder
func DirSelect(label string, files []os.FileInfo, foldersOnly bool) (os.FileInfo, error) {
	var displayed = make(map[string]os.FileInfo)
	for _, f := range files {
		if !foldersOnly || f.IsDir() {
			displayed[f.Name()] = f
		}
	}
	names := make([]string, 0, len(displayed))
	for k := range displayed {
		names = append(names, k)
	}
	// TODO: sort files/folders by date
	p := promptui.Select{
		Label: label,
		Items: names,
	}
	_, res, err := p.Run()
	ProcessPromptError(err)

	file, exists := displayed[res]
	if !exists {
		panic(errors.New("'" + res + "' could not be found"))
	}
	return file, err
}

// Confirm - Prompt for comfirmation
func Confirm(label string) bool {
	p := promptui.Prompt{
		Label: label + " (y/N)",
		Validate: func(input string) error {
			// if len(input) > 0 {
			// }
			return nil
		},
	}
	res, err := p.Run()
	ProcessPromptError(err)

	return strings.EqualFold(res, "Y")
}

// String - Prompt for a string input
func String(label string) (string, error) {
	p := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			if len(strings.TrimSpace(input)) == 0 {
				return errors.New("invalid empty imput")
			}
			return nil
		},
	}
	res, err := p.Run()
	ProcessPromptError(err)
	return res, err
}

// PositiveInteger - Prompt for a positive integer
func PositiveInteger(label string) (string, error) {
	p := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			trimmed := strings.TrimSpace(input)
			isNum, _ := scraper.CheckNumber(trimmed)
			if !isNum {
				return errors.New(trimmed + " is not valid")
			}
			return nil
		},
	}
	res, err := p.Run()
	ProcessPromptError(err)
	return res, err
}

func ProcessPromptError(err error) {
	if err == promptui.ErrInterrupt {
		os.Exit(0)
	} else if err != nil {
		panic(err)
	}
}
