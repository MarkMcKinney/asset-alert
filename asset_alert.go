package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	godotenv "github.com/joho/godotenv"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

func getEnvVars(key string) string {

	// load .env file and check for error
	// For MAC
	err := godotenv.Load("app.env")

	// check for error
	if err != nil {
		log.Fatalf(err.Error())
	}

	return os.Getenv(key)
}

func getAssetAction(asset string) [3]string {

	// Get dates for chart ranges
	today_date := time.Now()
	today_weekday := today_date.Weekday()
	date_delta := -1
	// If today is Monday, make yester_date find Friday
	if today_weekday == 1 {
		date_delta = -3
	}
	fmt.Println("Date delta:", date_delta)
	yester_date := today_date.AddDate(0, 0, date_delta)

	// Set chart parameters
	p := &chart.Params{
		Symbol:   asset,
		Start:    &datetime.Datetime{Month: int(yester_date.Month()), Day: int(yester_date.Day()), Year: int(yester_date.Year())},
		End:      &datetime.Datetime{Month: int(today_date.Month()), Day: int(today_date.Day()), Year: int(today_date.Year())},
		Interval: datetime.FiveDay,
	}

	// Get list of chart values as pointer
	iter := chart.Get(p)

	// Define array of price closes
	var priceCloses []string

	// Iterate over results. Will exit upon any error.
	for iter.Next() {
		// Convert decimal type to String
		close := iter.Bar().Close.String()
		// Append the closing price to the priceCloses array
		priceCloses = append(priceCloses, close)
	}

	// Catch an error, if there was one.
	if iter.Err() != nil {
		// Uh-oh!
		panic(iter.Err())
	}

	yesterCloseFlt, _ := strconv.ParseFloat(priceCloses[0], 64)
	todayCloseFlt, _ := strconv.ParseFloat(priceCloses[1], 64)
	percDelta := 100 * ((todayCloseFlt - yesterCloseFlt) / yesterCloseFlt)
	percDeltaStr := strconv.FormatFloat(percDelta, 'f', 2, 64)

	return [3]string{priceCloses[0], priceCloses[1], percDeltaStr}

}

func main() {
	// Get Telegram API key
	bot, err := tgbotapi.NewBotAPI(getEnvVars("TELEGRAM_BOT_API_KEY"))
	if err != nil {
		panic(err)
	}

	// Set message string
	var activityUpdate string

	// Read assets to lookup as an array
	assetList := strings.Split(getEnvVars("ASSETS"), ",")

	// Get Telegram receiver ID
	receiverID, err := strconv.ParseInt(getEnvVars("RECEIVER"), 10, 64)
	if err != nil {
		panic(err)
	}

	// Setup update message unless there was an error
	if err != nil {
		// Update message to error from fetching price
		activityUpdate = err.Error()
		panic(activityUpdate)
	} else {
		// Update message to current price of asset
		for _, asset := range assetList {
			// If asset is empty
			if asset != "" {
				assetPriceInfo := getAssetAction(asset)
				activityUpdate += asset + "\nToday's price: $" + assetPriceInfo[1] + "\n"
				activityUpdate += "Yesterday's price: $" + assetPriceInfo[0] + "\n"
				activityUpdate += "Delta: " + assetPriceInfo[2] + "%\n\n"
				fmt.Println(activityUpdate)
			}

		}
	}

	// Setup telegram message
	msg := tgbotapi.NewMessage(receiverID, activityUpdate)

	// Send message
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}

}
