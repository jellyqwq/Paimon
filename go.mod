module github.com/jellyqwq/Paimon

go 1.19

replace (
	github.com/jellyqwq/Paimon => ../Paimon
	github.com/jellyqwq/Paimon/news => ./news
	github.com/jellyqwq/Paimon/what => ./what
)

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/jellyqwq/Paimon/news v0.0.0-00010101000000-000000000000
	github.com/jellyqwq/Paimon/what v0.0.0
)
