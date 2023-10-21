package YearProgress

import (
	"fmt"
	"time"
)

type YearProgressConfig struct {
	FilledFlag string
	BlankFlag string
	Length int64
	barNow string
	bar string
	ChatID int64
	percentage float64
}

func NewYearProgressConfig() *YearProgressConfig {
	return &YearProgressConfig{
		FilledFlag: "▓",
		BlankFlag: "░",
		Length: 20,
		barNow: "",
		bar: "",
		ChatID: int64(0),
		percentage: float64(0),
	}
}

func (ypc *YearProgressConfig) GetYearProgress() string {
	ypc.display()
	if ypc.barNow != ypc.bar || ypc.barNow == "" || ypc.bar == "" {
		ypc.barNow = ypc.bar
		return fmt.Sprintf("%s %.f%%", ypc.bar, round(ypc.percentage, 1))
	} else {
		return ""
	}
}

func leapYear(date time.Time) int64 {
	if date.Year()%400 == 0 {
		return 366
	} else if date.Year()%4 == 0 && date.Year()%100 != 0 {
		return 366
	} else {
		return 365
	}
}

func getPercentage(date time.Time) float64 {
	CURRENT_DAY := date.YearDay()
	TOTAL_DAYS := leapYear(date)
	return float64(CURRENT_DAY*100) / float64(TOTAL_DAYS)
}


// Return string of year progress bar.
func (ypc *YearProgressConfig) display() string {
	ypc.percentage = getPercentage(time.Now())
	FILLED := float64(ypc.Length) * ypc.percentage / 100
	BLANK := float64(ypc.Length) - FILLED
	ypc.bar = ""
	for i := float64(0); i < FILLED; i++ {
		ypc.bar += ypc.FilledFlag
	}

	for i := float64(0); i < BLANK; i++ {
		ypc.bar += ypc.BlankFlag
	}
	// return 
	return ypc.bar
}

func round(x, unit float64) float64 {
	return float64(int64(x*unit+0.5)) / unit
}
