package imageextractor

import (
	"fmt"
)

type ParseResponse struct {
	ImagePath     string `json:"image_path"`
	Title         string `json:"title"`
	Score         string `json:"score"`
	PlayTimestamp string `json:"play_timestamp"`
	Rate          string `json:"rate"`
}

func ParseImage(imagePath string, language string) (ParseResponse, error) {

	Ie, err := New(imagePath)
	if err != nil {
		return ParseResponse{}, err
	}
	//extract player stats
	moneyScore, _ := Ie.Extractor(MoneyScore)
	playDate, _ := Ie.Extractor(PlayDate)
	playTime, _ := Ie.Extractor(PlayTime)

	title := ""
	switch language {
	case "eng":
		title, _ = Ie.Extractor(TitleEng)
		break
	case "jpn":
		title, _ = Ie.Extractor(TitleJpn)
		break
	default:
		title = "LANGUAGE NOT SUPPORTED!"
	}
	rate, _ := Ie.Extractor(AchievementRate)

	fmt.Println("Extracted Title:", title)
	fmt.Println("Extracted Score:", moneyScore)
	fmt.Println("Extracted Play Date:", playDate, playTime)
	fmt.Println("Extracted Achievement Rate:", rate, "%")

	return ParseResponse{
		ImagePath:     imagePath,
		Title:         title,
		Score:         moneyScore,
		PlayTimestamp: fmt.Sprintf("%s %s", playDate, playTime),
		Rate:          fmt.Sprintf("%s%%", rate),
	}, nil
}
