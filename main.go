package main

import (
	"bufio"
	"fish666/api"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"encoding/csv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var proxyList api.ProxyList

func init() {
	username := "devtest"
	password := "123123aaaA"
	node := "dc.jp-pr.oxylabs.io:12000"
	proxyList = api.ProxyList{
		NoProxy: "",
		Local:   "http://127.0.0.1:4780",
		Test:    fmt.Sprintf("http://%s:%s@%s", username, password, node),
	}
}

func main() {
	// os.Setenv("FYNE_FONT", "Songti.ttc")
	a := app.New()
	w := a.NewWindow("GUI Program")

	w.SetContent(UI(w))
	// 设置全局字体
	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}

func CleanTable(tables *[][7]string) {
	if len(*tables) > 0 {
		(*tables) = (*tables)[:1]
	}
}

func GetCsvSlices(old [][7]string) [][]string {
	ret := make([][]string, len(old))
	for i := range old {
		ret[i] = make([]string, len(old[i]))
		for j := range old[i] {
			if old[i][j] == "" {
				ret[i][j] = "0"
			} else {
				ret[i][j] = old[i][j]
			}
		}
	}
	return ret
}

func UI(window fyne.Window) *fyne.Container {
	tableItems := [][7]string{
		{"website", "all", "hour", "day", "week", "month", "year"},
	}
	// var wg sync.WaitGroup
	var mu sync.Mutex

	wordEntry := widget.NewEntry()
	wordEntry.SetPlaceHolder("Plase entry keywords")
	websiteEntry := widget.NewEntry()
	websiteEntry.SetPlaceHolder("Plase entry the website address")
	searchTypeSelect := widget.NewSelect([]string{"any", "url"}, nil)
	searchTypeSelect.SetSelectedIndex(0)
	searchTimeSelect := widget.NewSelect([]string{"all", "hour", "day", "week", "month", "year"}, nil)
	searchTimeSelect.SetSelectedIndex(0)

	proxyArray := make([]string, 0)
	t := reflect.TypeOf(proxyList)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			proxyArray = append(proxyArray, field.Name)
		}
	}
	proxySelect := widget.NewSelect(proxyArray, nil)
	proxySelect.SetSelectedIndex(0)

	logEntry := widget.NewMultiLineEntry()
	logEntry.Wrapping = fyne.TextWrapWord
	logEntry.Disable()

	table := widget.NewTable(
		func() (int, int) { return len(tableItems), len(tableItems[0]) },
		func() fyne.CanvasObject {
			label := widget.NewLabel("Item")
			return label
		},
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			var itemData string
			if itemData = tableItems[tci.Row][tci.Col]; itemData == "" {
				itemData = "0"
			}
			co.(*widget.Label).SetText(itemData)
			// fmt.Println(tableItems[tci.Row][tci.Col])
		},
	)
	table.SetColumnWidth(0, 120)
	table.SetColumnWidth(1, 120)
	table.SetColumnWidth(2, 120)
	table.SetColumnWidth(3, 120)
	table.SetColumnWidth(4, 120)
	table.SetColumnWidth(5, 120)
	table.SetColumnWidth(6, 120)

	// table.SetRowHeight(-1, 20)

	tableBox := container.NewVScroll(table)
	tableBox.SetMinSize(fyne.NewSize(650, 300))

	openFileButton := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		if err != nil || uc == nil {
			logEntry.SetText("Read file failed.")
			return
		}
		filePath := uc.URI().Path()
		fmt.Printf("Path: %s", filePath)
		if filePath[len(filePath)-4:] != ".txt" {
			logEntry.SetText("no txt")
			return
		}
		logEntry.SetText("selected: " + filePath)
		file, err := os.Open(filePath)
		if err != nil {
			logEntry.SetText("Open file failed.")
			return
		}
		defer file.Close()
		CleanTable(&tableItems)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			tableItems = append(tableItems, [7]string{line})
		}
		table.Refresh()
		if scanner.Err() != nil {
			logEntry.SetText("Read line filed.")
		}
		websiteEntry.SetText("")
	}, window)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Word", Widget: wordEntry},
			{Text: "Time", Widget: searchTimeSelect},
			{Text: "Web", Widget: websiteEntry},
			{Text: "", Widget: widget.NewButton("Import", func() {
				openFileButton.Show()
			})},
			{Text: "Type", Widget: searchTypeSelect},

			{Text: "Proxy", Widget: proxySelect},
			{Text: "", Widget: tableBox},
			{Text: "", Widget: widget.NewButton("Export", func() {
				file, err := os.Create("output.csv")
				if err != nil {
					logEntry.SetText("output.csv create failed.")
					return
				}
				defer file.Close()
				writer := csv.NewWriter(file)
				err = writer.WriteAll(GetCsvSlices(tableItems))
				if err != nil {
					logEntry.SetText("output.csv writeall failed.")
					return
				}
				writer.Flush()
				if err := writer.Error(); err != nil {
					logEntry.SetText("output.csv flush failed.")
					return
				}
			})},
			{Text: "Logs", Widget: logEntry},
		},
		OnSubmit: func() {
			// 反射获取代理映射值
			proxyText := reflect.ValueOf(proxyList).FieldByName(proxySelect.Selected).Interface().(string)

			if websiteEntry.Text != "" && strings.Contains(websiteEntry.Text, ".") {
				CleanTable(&tableItems)
				tableItems = append(tableItems, [7]string{websiteEntry.Text})
				table.Refresh()
			}

			for index := 1; index < len(tableItems); index++ {
				curweb := tableItems[index][0]
				// fmt.Printf("ste: %s\n", curweb)
				para := api.UIParameter{
					Word:  wordEntry.Text,
					Time:  searchTimeSelect.Selected,
					Web:   curweb,
					Type:  searchTypeSelect.Selected,
					Proxy: proxyText,
				}

				i := index
				go api.GetSearchRet(para, func(s string, err error) {
					mu.Lock()
					if err != nil || s == "" {return}
					fmt.Printf("i: %d\n", i)
					switch para.Time {
					case "all":
						tableItems[i][1] = s
					case "hour":
						tableItems[i][2] = s
					case "day":
						tableItems[i][3] = s
					case "week":
						tableItems[i][4] = s
					case "month":
						tableItems[i][5] = s
					case "year":
						tableItems[i][6] = s
					}
					table.Refresh()
					mu.Unlock()
					// fmt.Printf("GetSearchRet: %s", s)
				})
			}

		},
		OnCancel: func() {
			wordEntry.SetText("")
			websiteEntry.SetText("")
			searchTypeSelect.SetSelectedIndex(0)
			searchTimeSelect.SetSelectedIndex(0)
			proxySelect.SetSelectedIndex(0)
			logEntry.SetText("Cancel\n")
		},
	}

	return container.NewVBox(
		form,
	)
}