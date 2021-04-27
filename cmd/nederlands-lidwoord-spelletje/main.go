package main

import (
	"database/sql"
	"flag"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/averageflow/nederlands-lidwoord-spelletje/internal/words"
)

func main() {
	db, err := sql.Open("sqlite3", "../../storage/lidwoord.sqlite")
	if err != nil {
		panic(err.Error())
	}

	woordFlag := flag.String("woord", "", "Woord")
	lidwoordFlag := flag.String("lidwoord", "", "Lidwoord")
	pluralFlag := flag.String("plural", "", "Plural")

	flag.Parse()

	if woordFlag != nil && lidwoordFlag != nil && *woordFlag != "" && *lidwoordFlag != "" && pluralFlag != nil && *pluralFlag != "" {
		err := words.InsertNewWord(db, *woordFlag, *lidwoordFlag, *pluralFlag)
		if err != nil {
			panic(err.Error())
		}
		return
	}

	a := app.New()
	w := a.NewWindow("Nederlands Lidwoorden")
	w.Resize(fyne.Size{Height: 600, Width: 800})

	//lidwoordTitle,
	//	widget.NewButton("Hi!", func() {
	//		lidwoordTitle.SetText("Welcome :)")
	//	}),)

	lidwoordTitle := widget.NewLabel("Welkom in het lidwoord-spelletje!")
	lidwoordTitle.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	//pluralTitle := widget.NewLabel("Welkom in het plural-spelletje!")
	//pluralTitle.TextStyle = fyne.TextStyle{
	//	Bold: true,
	//}

	reloadIcon, _ := fyne.LoadResourceFromURLString("")

	var currentWord *words.GameWord

	currentWord, err = words.GetRandomWord(db)
	if err != nil {
		panic(err.Error())
	}
	gameWordLabel := widget.NewLabel(currentWord.Content)
	gameWordLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	lidwoordContainer := container.NewVBox(
		container.NewCenter(container.NewHBox(
			lidwoordTitle,
			widget.NewButtonWithIcon("RELOAD", reloadIcon, func() {
				currentWord, err = words.GetRandomWord(db)
				if err != nil {
					panic(err.Error())
				}
				gameWordLabel.Text = currentWord.Content
			}),
		)),
		container.NewCenter(gameWordLabel),
		container.NewCenter(container.NewHBox(
			widget.NewButton("de", func() {
				if currentWord.Lidwoord != "de" {
					fmt.Println("FAILED")
				} else {
					fmt.Println("SUCCESS")
					currentWord, err = words.GetRandomWord(db)
					if err != nil {
						panic(err.Error())
					}
					gameWordLabel.Text = currentWord.Content
				}
			}),
			widget.NewButton("het", func() {
				if currentWord.Lidwoord != "het" {
					fmt.Println("FAILED")
				} else {
					fmt.Println("SUCCESS")
					currentWord, err = words.GetRandomWord(db)
					if err != nil {
						panic(err.Error())
					}
					gameWordLabel.Text = currentWord.Content
				}

			}),
		)),
	)

	w.SetContent(container.NewAppTabs(
		&container.TabItem{
			Text:    "Lidwoord",
			Icon:    nil,
			Content: lidwoordContainer,
		},
		//&container.TabItem{
		//	Text: "Plural",
		//	Icon: nil,
		//	Content: container.NewVBox(
		//		pluralTitle,
		//	),
		//},
	))

	w.ShowAndRun()
}
