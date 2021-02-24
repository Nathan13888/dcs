package prompt

import (
	"dcs/scraper"
	"errors"
	"strings"

	"github.com/manifoldco/promptui"
)

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
