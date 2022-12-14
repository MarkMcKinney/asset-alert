# Telegram Asset Alert Bot (Go)

I've been wanting to learn Go and I just needed a good project to get inspired by. 
Recently, I've been wanting custom asset movement alerts for free so I decided to build that in Go.

Currently a continual work-in-progress.

### Features

- [x] Send ETF or stock price movements (current price from the previous day's closing price) updates for ETFs and stocks over Telegram.
- [x] If today is Monday, get the closing price from the previous Friday.
- [] Use Go channels to get all assets at the same time. This will reduce execution time.

### Get Started

1. Create a file called app.env
2. You'll need a value for the following environmental variables:
    a. TELEGRAM_BOT_API_KEY=XXXXXXXXX
        - Your Telegram bot API key.
    b. ASSETS=GOOG,AAPL,AMD
        - List of assets that you want to get notifications for (no max).
3. Currently, I have this running via a cronjob that runs once a day on weekdays.