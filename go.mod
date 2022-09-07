module github.com/jellyqwq/Paimon

go 1.19

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
    github.com/jellyqwq/Paimon/news v0.0.0
)

replace (
    github.com/jellyqwq/Paimon => ../Paimon
    github.com/jellyqwq/Paimon/news => ./news
)
