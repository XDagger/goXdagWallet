package components

import (
	"compress/gzip"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/buger/jsonparser"
	"goXdagWallet/i18n"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type HistoryItem struct {
	Direction string
	Address   string
	Amount    string
	Time      string
	Remark    string
}

var historyStatus = binding.NewString()
var HistoryContainer *fyne.Container
var HistoryProgressContainer *fyne.Container
var HistoryRefreshContainer *fyne.Container
var HistoryData = make([]HistoryItem, 0, 10)
var HistoryTable *widget.Table

func HistoryPage(w fyne.Window) *fyne.Container {
	refreshBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		historyStatus.Set(i18n.GetString("WalletWindow_HistoryBusy"))
		HistoryProgressContainer.Show()
		HistoryRefreshContainer.Hide()
		HistoryContainer.Remove(HistoryTable)
		go refreshTable()
	})
	refreshBtn.Importance = widget.HighImportance

	HistoryProgressContainer = container.New(layout.NewPaddedLayout(), widget.NewProgressBarInfinite())
	HistoryRefreshContainer = container.New(layout.NewHBoxLayout(), layout.NewSpacer(), refreshBtn)
	HistoryRefreshContainer.Hide()
	HistoryContainer = container.New(layout.NewMaxLayout())
	HistoryTitle := widget.NewLabelWithData(historyStatus)
	historyStatus.Set(i18n.GetString("WalletWindow_HistoryBusy"))

	go refreshTable()
	top := container.NewVBox(
		container.NewHBox(layout.NewSpacer(), HistoryTitle, layout.NewSpacer()),
		HistoryRefreshContainer,
		HistoryProgressContainer)
	return container.New(
		layout.NewBorderLayout(top, nil, nil, nil),
		top,
		HistoryContainer)
}

func refreshTable() {
	var body []byte
	err := getUrl(Address, "https://explorer.xdag.io/api/block", &body)
	if err != nil {
		historyStatus.Set(i18n.GetString("WalletWindow_HistoryBusy"))
		return
	}

	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		direction, _ := jsonparser.GetString(value, "direction")
		address, _ := jsonparser.GetString(value, "address")
		amount, _ := jsonparser.GetString(value, "amount")
		times, _ := jsonparser.GetString(value, "time")
		remark, _ := jsonparser.GetString(value, "remark")
		HistoryData = append(HistoryData, HistoryItem{
			Direction: direction,
			Address:   address,
			Amount:    amount,
			Time:      times,
			Remark:    remark,
		})
	}, "block_as_address")
	HistoryProgressContainer.Hide()
	HistoryRefreshContainer.Show()
	if err != nil {
		historyStatus.Set(i18n.GetString("WalletWindow_HistoryBusy"))
		return
	}

	historyStatus.Set(i18n.GetString("WalletWindow_HistoryColumns_BlockAddress") + " : " + Address)

	HistoryTable = widget.NewTable(
		func() (int, int) { return len(HistoryData) + 1, 5 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row == 0 {
				switch id.Col {
				case 0:
					label.SetText(i18n.GetString("WalletWindow_HistoryColumns_Direction"))
				case 1:
					label.SetText(i18n.GetString("WalletWindow_HistoryColumns_Amount"))
				case 2:
					label.SetText(i18n.GetString("WalletWindow_HistoryColumns_PartnerAddress"))
				case 3:
					label.SetText(i18n.GetString("WalletWindow_HistoryColumns_TimeStamp"))
				case 4:
					label.SetText(i18n.GetString("WalletWindow_HistoryColumns_Remark"))
				default:
					label.SetText("cell")
				}
			} else {
				switch id.Col {
				case 0:
					label.SetText(HistoryData[id.Row-1].Direction)
				case 1:
					label.SetText(HistoryData[id.Row-1].Amount)
				case 2:
					label.SetText(HistoryData[id.Row-1].Address)
				case 3:
					label.SetText(HistoryData[id.Row-1].Time)
				case 4:
					label.SetText(HistoryData[id.Row-1].Remark)
				default:
					label.SetText("cell")
				}
			}

		})
	HistoryTable.SetColumnWidth(0, 82)
	HistoryTable.SetColumnWidth(1, 142)
	HistoryTable.SetColumnWidth(2, 362)
	HistoryTable.SetColumnWidth(3, 222)
	HistoryTable.SetColumnWidth(4, 152)

	HistoryContainer.Add(HistoryTable)
}
func getUrl(params, apiUrl string, body *[]byte) error {
	urlString := apiUrl + "/" + params
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return err
	}

	// 表单方式(必须)
	//req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Apache-HttpClient/4.3.1")

	client := &http.Client{
		Timeout: 8 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	*body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}
