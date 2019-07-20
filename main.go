package main

import (
	"fmt"

	"github.com/wicoady1/gtdr-score-parser/imageextractor"
)

func main() {
	const IMAGE_PATH = "asset/S__9953338.jpg"

	Ie, _ := imageextractor.New(IMAGE_PATH)
	//extract player stats
	moneyScore, _ := Ie.Extractor(imageextractor.MoneyScore)
	playDate, _ := Ie.Extractor(imageextractor.PlayDate)
	playTime, _ := Ie.Extractor(imageextractor.PlayTime)
	title, _ := Ie.Extractor(imageextractor.TitleJpn)
	rate, _ := Ie.Extractor(imageextractor.AchievementRate)

	fmt.Println("Extracted Title:", title)
	fmt.Println("Extracted Score:", moneyScore)
	fmt.Println("Extracted Play Date:", playDate, playTime)
	fmt.Println("Extracted Achievement Rate:", rate, "%")
}
