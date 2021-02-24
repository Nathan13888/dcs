package prompt

import (
	"dcs/scraper"
	"errors"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

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
	if err != nil {
		panic(err)
	}
	file, exists := displayed[res]
	if !exists {
		panic(errors.New("'" + res + "' could not be found"))
	}
	return file, err
}

// Confirm - Prompt for comfirmation
func Confirm(label string) (string, error) {
	p := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}
	res, err := p.Run()
	if err != nil {
		panic(err)
	}
	return res, err
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
	if err != nil {
		panic(err)
	}
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
	if err != nil {
		panic(err)
	}
	return res, err
}

// LimitedPositiveInteger - Prompt for a positive integer with an upper bound
func LimitedPositiveInteger(label string, upper int) (string, error) {
	p := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			trimmed := strings.TrimSpace(input)
			isNum, num := scraper.CheckNumber(trimmed)
			if !isNum || num > upper {
				return errors.New(trimmed + " is not valid")
			}
			return nil
		},
	}
	res, err := p.Run()
	if err != nil {
		panic(err)
	}
	return res, err
}
