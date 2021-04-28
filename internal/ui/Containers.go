package ui

import (
	"database/sql"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/averageflow/nederlands-lidwoord-spelletje/internal/words"
)

func GetLidwoordContainer(db *sql.DB) *fyne.Container {
	var currentWord *words.GameWord

	currentWord, err := words.GetRandomWord(db)
	if err != nil {
		panic(err.Error())
	}
	gameWordLabel := widget.NewLabel(currentWord.Content)
	gameWordLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	gameStatusLabel := widget.NewLabel("")
	gameStatusLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	reloadIcon, _ := fyne.LoadResourceFromURLString("")
	lidwoordTitle := widget.NewLabel("Het lidwoord-spelletje!")
	lidwoordTitle.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	return container.NewVBox(
		container.NewCenter(container.NewHBox(
			lidwoordTitle,
			widget.NewButtonWithIcon("skip >>", reloadIcon, func() {
				currentWord, err = words.GetRandomWord(db)
				if err != nil {
					gameStatusLabel.Text = err.Error()
				}
				gameWordLabel.Text = currentWord.Content
			}),
		)),
		container.NewCenter(gameWordLabel),
		container.NewCenter(container.NewHBox(
			widget.NewButton("de", func() {
				if currentWord.Lidwoord != "de" {
					gameStatusLabel.Text = "FAILED!"
				} else {
					currentWord, err = words.GetRandomWord(db)
					if err != nil {
						gameStatusLabel.Text = err.Error()
					}
					gameWordLabel.Text = currentWord.Content
					gameStatusLabel.Text = ""
				}
			}),
			widget.NewButton("het", func() {
				if currentWord.Lidwoord != "het" {
					gameStatusLabel.Text = "FAILED!"
				} else {
					currentWord, err = words.GetRandomWord(db)
					if err != nil {
						gameStatusLabel.Text = err.Error()
					}
					gameWordLabel.Text = currentWord.Content
					gameStatusLabel.Text = ""
				}
			}),
		)),
		container.NewCenter(container.NewVBox(gameStatusLabel)),
	)
}

func GetPluralContainer(db *sql.DB) *fyne.Container {
	var currentWord *words.GameWordWithPlural

	currentWord, err := words.GetRandomWordWithPlural(db)
	if err != nil {
		panic(err.Error())
	}
	gameWordLabel := widget.NewLabel(fmt.Sprintf("%s %s -> de ?", currentWord.Lidwoord, currentWord.Content))
	gameWordLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	gameStatusLabel := widget.NewLabel("")
	gameStatusLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	reloadIcon, _ := fyne.LoadResourceFromURLString("")
	lidwoordTitle := widget.NewLabel("Het plural-spelletje!")
	lidwoordTitle.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	userInput := widget.NewEntry()

	userInput.OnSubmitted = func(s string) {
		if !strings.EqualFold(currentWord.Plural, s) {
			gameStatusLabel.Text = "FAILED!"
		} else {
			currentWord, err = words.GetRandomWordWithPlural(
				db)
			if err != nil {
				gameStatusLabel.Text = err.Error()
			}
			gameStatusLabel.Text = ""
			userInput.SetText("")
			gameWordLabel.Text = fmt.Sprintf("%s %s -> de ?", currentWord.Lidwoord, currentWord.Content)
		}
	}

	return container.NewVBox(
		container.NewCenter(container.NewHBox(
			lidwoordTitle,
			widget.NewButtonWithIcon("skip >>", reloadIcon, func() {
				currentWord, err = words.GetRandomWordWithPlural(
					db)
				if err != nil {
					gameStatusLabel.Text = err.Error()
				}
				gameWordLabel.Text = fmt.Sprintf("%s %s -> de ?", currentWord.Lidwoord, currentWord.Content)
			}),
		)),
		container.NewCenter(gameWordLabel),
		container.NewMax(userInput),
		container.NewCenter(container.NewHBox(
			widget.NewButton("Clear", func() {
				gameStatusLabel.Text = ""
				userInput.SetText("")
			}),
			widget.NewButton("Submit", func() {
				if !strings.EqualFold(currentWord.Plural, userInput.Text) {
					gameStatusLabel.Text = "FAILED!"
				} else {
					currentWord, err = words.GetRandomWordWithPlural(
						db)
					if err != nil {
						gameStatusLabel.Text = err.Error()
					}
					gameStatusLabel.Text = ""
					userInput.SetText("")
					gameWordLabel.Text = fmt.Sprintf("%s %s -> de ?", currentWord.Lidwoord, currentWord.Content)
				}

			}),
		)),
		container.NewCenter(container.NewVBox(gameStatusLabel)),
	)
}
