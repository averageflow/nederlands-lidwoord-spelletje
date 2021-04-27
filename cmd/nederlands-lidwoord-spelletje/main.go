package main

import (
	"database/sql"

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
	gameWordLabel := widget.NewLabel(words.GetRandomWord(db))
	gameWordLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	lidwoordContainer := container.NewVBox(
		container.NewCenter(container.NewHBox(
			lidwoordTitle,
			widget.NewButtonWithIcon("RELOAD", reloadIcon, func() {
				gameWordLabel.Text = words.GetRandomWord(db)
			}),
		)),
		container.NewCenter(gameWordLabel),
		container.NewCenter(container.NewHBox(
			widget.NewButton("de", func() {}),
			widget.NewButton("het", func() {}),
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
