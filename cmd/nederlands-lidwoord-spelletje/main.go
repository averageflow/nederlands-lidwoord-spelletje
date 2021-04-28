package main

import (
	"database/sql"
	"flag"

	"github.com/averageflow/nederlands-lidwoord-spelletje/internal/ui"

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

	if woordFlag != nil && lidwoordFlag != nil && *woordFlag != "" && *lidwoordFlag != "" && pluralFlag != nil {
		err := words.InsertNewWord(db, *woordFlag, *lidwoordFlag, *pluralFlag)
		if err != nil {
			panic(err.Error())
		}
		return
	}

	a := app.New()
	w := a.NewWindow("Nederlands Lidwoorden")
	w.Resize(fyne.Size{Height: 600, Width: 800})

	pluralTitle := widget.NewLabel("Het plural-spelletje!")
	pluralTitle.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	lidwoordContainer := ui.GetLidwoordContainer(db)
	pluralContainer := ui.GetPluralContainer(db)

	w.SetContent(container.NewAppTabs(
		&container.TabItem{
			Text:    "Lidwoord",
			Icon:    nil,
			Content: lidwoordContainer,
		},
		&container.TabItem{
			Text:    "Plural",
			Icon:    nil,
			Content: pluralContainer,
		},
	))

	w.ShowAndRun()
}
