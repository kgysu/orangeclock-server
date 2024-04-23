package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printer = message.NewPrinter(language.English)

// Mempool

func LoadMempoolData() string {
	data := ""

	price, err := getCurrentPrice()
	if err != nil {
		data += fmt.Sprintf("Error prices: %s\n", err.Error())
	} else {
		data += fmt.Sprintf("PRICE (%s)\n", price.TimeS())
		data += printer.Sprintf(" $%d   +%d\n", price.Usd, price.Chf)
		data += printer.Sprintf("1$: %d\n", price.MoscowTime())
		data += "\n"
	}

	tipHeight, err := getCurrentTipHeightInternal()
	if err != nil {
		data += fmt.Sprintf("Error tip: %s\n", err.Error())
	} else {
		data += printer.Sprintf("      %d @\n", tipHeight)
		data += "\n"
	}

	halving, err := getNextHalving(tipHeight)
	if err != nil {
		data += fmt.Sprintf("Error halving: %s\n", err.Error())
	} else {
		data += "HALVING\n"
		data += halving + "\n"
		data += "\n"
	}

	fees, err := getFees()
	if err != nil {
		data += fmt.Sprintf("Error fees: %s\n", err.Error())
	} else {
		data += "FEES\n"
		data += fees.String() + "\n"
		data += "\n"
	}

	return data
}

// Prices

type PriceResponse struct {
	Time int64 `json:"time"`
	Usd  int   `json:"USD"`
	Eur  int   `json:"EUR"`
	Chf  int   `json:"CHF"`
}

func (pr PriceResponse) String() string {
	return printer.Sprintf("$%d +%d (%s)", pr.Usd, pr.Chf, pr.TimeS())
}

func (pr PriceResponse) TimeS() string {
	t := time.Unix(pr.Time, 0)
	return t.Format(time.TimeOnly)
}

func (pr PriceResponse) MoscowTime() int {
	return 100_000_000 / pr.Usd
}

func getCurrentPrice() (*PriceResponse, error) {
	resp, err := http.Get("https://mempool.space/api/v1/prices")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("no valid status code returned: %d %s\n", resp.StatusCode, resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var prices PriceResponse
	err = json.Unmarshal(data, &prices)
	if err != nil {
		return nil, err
	}
	return &prices, nil
}

// Latest Block (tip height)

func getCurrentTipHeightInternal() (int, error) {
	out, err := exec.Command("bitcoin-cli", "getblockcount").Output()
	if err != nil {
		return 0, err
	}
	val := strings.Replace(string(out), "\n", "", 1)
	res, err := strconv.Atoi(val)
	if err != nil {
		return 0, nil
	}
	return res, nil
}

func getCurrentTipHeight() (int, error) {
	client := http.DefaultClient
	resp, err := client.Get("https://mempool.space/api/blocks/tip/height")
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("no valid status code returned: %d %s\n", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	res, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return res, nil
}

func getNextHalving(currentHeight int) (string, error) {
	nextHalvingBlock := ((currentHeight / 210_000) + 1) * 210_000
	blocksToGo := nextHalvingBlock - currentHeight
	percentage := (100 / nextHalvingBlock) * blocksToGo
	return printer.Sprintf(" %d (%d%s)", blocksToGo, percentage, "%"), nil
}

// Fees

type FeeResponse struct {
	FastestFee  int `json:"fastestFee"`
	HalfHourFee int `json:"halfHourFee"`
	HourFee     int `json:"hourFee"`
	EconomyFee  int `json:"economyFee"`
	MinimumFee  int `json:"minimumFee"`
}

func (fr FeeResponse) String() string {
	return printer.Sprintf(" %d - %d - %d", fr.HourFee, fr.HalfHourFee, fr.FastestFee)
}

func getFees() (*FeeResponse, error) {
	client := http.DefaultClient
	resp, err := client.Get("https://mempool.space/api/v1/fees/recommended")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("no valid status code returned: %d %s\n", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var prices FeeResponse
	err = json.Unmarshal(data, &prices)
	if err != nil {
		return nil, err
	}
	return &prices, nil
}
