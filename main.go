package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"image/color"
	"os"
	"reflect"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"fish666/api"
	// "fish666/theme"
	"fish666/tool"
	"fish666/views"
)

var proxyList api.ProxyList
var myapp fyne.App
var mywindow fyne.Window

func init() {
	proxyList = api.ProxyList{
		NoProxy: "",
		Local:   fmt.Sprintf("%s:%s", tool.ReadIni("local_proxy", "host"), tool.ReadIni("local_proxy", "port")),
		Alpha:   fmt.Sprintf("%s:%s", tool.ReadIni("alpha_proxy", "host"), tool.ReadIni("alpha_proxy", "port")),
		Beta:    fmt.Sprintf("%s:%s", tool.ReadIni("beta_proxy", "host"), tool.ReadIni("beta_proxy", "port")),
		Gamma:   fmt.Sprintf("%s:%s", tool.ReadIni("gamma_proxy", "host"), tool.ReadIni("gamma_proxy", "port")),
	}
	fmt.Println(proxyList)
}

func main() {
	os.Setenv("FYNE_FONT", "./fonts/"+tool.ReadIni("ui", "font"))
	myapp = app.New()
	// myapp.Settings().SetTheme(&theme.MyTheme{})
	mywindow = myapp.NewWindow("GUI Program")

	mywindow.SetContent(UI(mywindow))
	mywindow.Resize(fyne.NewSize(1280, 600))
	mywindow.ShowAndRun()
}

func Alert(str string) {
	dialog.ShowInformation("Alert", str, mywindow)
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

	table := api.TableData{
		TopTableItems: [][7]string{{"website", "all", "hour", "day", "week", "month", "year"}},
		TableItems:    [][7]string{},
	}

	searchTypeArr := []string{"any", "url"}
	searchTimeArr := []string{"all", "hour", "day", "week", "month", "year"}

	wordEntry := widget.NewEntry()
	wordEntry.SetPlaceHolder("Plase entry keywords")
	wordEntry.SetText(tool.ReadIni("ui", "default_word"))

	websiteEntry := widget.NewEntry()
	websiteEntry.SetPlaceHolder("Plase entry the website address")
	websiteEntry.SetText(tool.ReadIni("ui", "default_web"))

	searchTypeSelect := widget.NewSelect(searchTypeArr, nil)
	searchTypeSelect.SetSelectedIndex(tool.GetKeyIndex(searchTypeArr, tool.ReadIni("ui", "default_type")))
	searchTimeSelect := widget.NewSelect(searchTimeArr, nil)
	searchTimeSelect.SetSelectedIndex(tool.GetKeyIndex(searchTimeArr, tool.ReadIni("ui", "default_time")))

	proxyArray := make([]string, 0)
	t := reflect.TypeOf(proxyList)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			proxyArray = append(proxyArray, field.Name)
		}
	}
	proxySelect := widget.NewSelect(proxyArray, nil)
	proxySelect.SetSelectedIndex(tool.GetKeyIndex(proxyArray, tool.ReadIni("ui", "default_proxy")))

	test := canvas.NewText("aa", color.White)
	test.TextSize = 2

	dosort := true
	var tableBody *widget.Table
	tableTop := widget.NewTable(
		func() (int, int) { return len(table.TopTableItems), len(table.TopTableItems[0]) },
		func() fyne.CanvasObject {
			button := widget.NewButton("TestTestTestTestTest", nil)
			return button
		},
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			button := co.(*widget.Button)
			button.SetText(table.TopTableItems[0][tci.Col])
			button.OnTapped = func() {
				if tci.Row == 0 && tci.Col > 0 {
					// 开始排序
					fmt.Printf("开始排序:%d", tci.Col)
					dosort = !dosort
					if dosort {
						table.Ascend(tci.Col)
					} else {
						table.Descend(tci.Col)
					}
					tableBody.Refresh()
				}
			}
		},
	)
	fmt.Printf("minisize: %.2f, %.2f\n", tableTop.MinSize().Width, tableTop.MinSize().Height)

	// tableTop.SetRowHeight(0, 35)
	tableTop.SetColumnWidth(0, 200)

	tableBody = widget.NewTable(
		func() (int, int) { return len(table.TableItems), len(table.TableItems[0]) },
		func() fyne.CanvasObject {
			return widget.NewLabel("TestTestTestTestTestTest")
		},
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			var itemData string
			if itemData = table.TableItems[tci.Row][tci.Col]; itemData == "" {
				itemData = "0"
			}
			co.(*widget.Label).SetText(itemData)
		},
	)

	tableBody.SetColumnWidth(0, 200)
	tableBody.SetColumnWidth(1, tableTop.MinSize().Width)
	tableBody.SetColumnWidth(2, tableTop.MinSize().Width)
	tableBody.SetColumnWidth(3, tableTop.MinSize().Width)
	tableBody.SetColumnWidth(4, tableTop.MinSize().Width)
	tableBody.SetColumnWidth(5, tableTop.MinSize().Width)
	tableBody.SetColumnWidth(6, tableTop.MinSize().Width)

	// tableBox := container.NewVBox(tableTop, tableBody)
	tableBox := container.NewVScroll(tableBody)
	tableBox.SetMinSize(fyne.NewSize(800, 300))

	ml := views.NewMyLogs()

	openFileButton := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		if err != nil || uc == nil {
			fmt.Println("Read file failed.")
			return
		}
		filePath := uc.URI().Path()
		fmt.Printf("Path: %s", filePath)
		if filePath[len(filePath)-4:] != ".txt" {
			fmt.Println("no txt")
			return
		}
		fmt.Println("selected: " + filePath)
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Open file failed.")
			return
		}
		defer file.Close()
		// 清空列表
		table.Clean()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			table.TableItems = append(table.TableItems, [7]string{line})
		}
		tableBody.Refresh()
		if scanner.Err() != nil {
			fmt.Println("Read line filed.")
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
			{Text: "", Widget: tableTop},
			{Text: "", Widget: tableBox},

			// {Text: "", Widget: tableBody},
			{Text: "", Widget: widget.NewButton("Export", func() {
				file, err := os.Create("output.csv")
				if err != nil {
					fmt.Println("output.csv create failed.")
					return
				}
				defer file.Close()
				writer := csv.NewWriter(file)
				err = writer.WriteAll(GetCsvSlices(table.TableItems))
				if err != nil {
					Alert("output.csv writeall failed.")
					return
				}
				writer.Flush()
				if err := writer.Error(); err != nil {
					Alert("output.csv flush failed.")
					return
				}
				ml.Info("Export out.csv success.")
				Alert("Export out.csv success.")

			})},
			{Text: "Logs", Widget: ml.GetView()},
		},
		OnSubmit: func() {
			// 反射获取代理映射值
			ml.Info("Start Search")
			proxyText := reflect.ValueOf(proxyList).FieldByName(proxySelect.Selected).Interface().(string)

			if websiteEntry.Text != "" && strings.Contains(websiteEntry.Text, ".") {
				table.Clean()
				table.TableItems = append(table.TableItems, [7]string{websiteEntry.Text})
				tableBody.Refresh()
			}

			go func() {
				for index := 0; index < len(table.TableItems); index++ {
					curweb := table.TableItems[index][0]
					ml.Info(fmt.Sprintf("[+] index: [%d];web: [%s]", index, curweb))

					para := api.UIParameter{
						Word:  wordEntry.Text,
						Time:  searchTimeSelect.Selected,
						Web:   curweb,
						Type:  searchTypeSelect.Selected,
						Proxy: proxyText,
					}

					api.GetSearchRet(para, func(s string, err error) {
						if err != nil || s == "" {
							ml.Info(fmt.Sprintf("[!] index: [%d];web: [%s]; search failed.", index, curweb))
							ml.Info(fmt.Sprintf("[!] error: %s", err.Error()))
							return
						}
						s = strings.ReplaceAll(s, ",", "")
						fmt.Printf("i: %d\n", index)
						switch para.Time {
						case "all":
							table.TableItems[index][1] = s
						case "hour":
							table.TableItems[index][2] = s
						case "day":
							table.TableItems[index][3] = s
						case "week":
							table.TableItems[index][4] = s
						case "month":
							table.TableItems[index][5] = s
						case "year":
							table.TableItems[index][6] = s
						}
						tableBody.Refresh()
						fmt.Println("[*] 搜索完成.")
						ml.Info(fmt.Sprintf("[+] index: [%d];web: [%s]; search success.", index, curweb))
					})
				}
				Alert("Searched.")
			}()
		},
		OnCancel: func() {
			wordEntry.SetText("")
			websiteEntry.SetText("")
			ml.Cancel()
			// table.Clean()
			// table.TableItems = append(table.TableItems, [7]string{})
			// tableBody.Refresh()
			ml.Info("Cancel all")
			// 清除表格

		},
	}

	return container.NewVBox(
		form,
	)
}
