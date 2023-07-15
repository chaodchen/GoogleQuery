package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type MyLogs struct {
	Logs []string
	box *fyne.Container
	scroll *container.Scroll
}

func NewMyLogs() *MyLogs {
	m := &MyLogs{
		box: container.NewVBox(),
	}
	return m
}

func (l *MyLogs) GetView() *container.Scroll {
	l.scroll = container.NewScroll(l.box)
	l.scroll.SetMinSize(fyne.NewSize(0, 80))
	return l.scroll
}

func (l *MyLogs) Info(str string) {
	t := canvas.NewText(str, nil)
	t.TextSize = 16
	l.box.Add(t)
	l.Logs = append(l.Logs, str)
	l.scroll.ScrollToBottom()
	// l.box.Refresh()
}

func (l *MyLogs) Cancel() {
	l.box.RemoveAll()
}