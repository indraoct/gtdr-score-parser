package imageextractor

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract"
)

type ImageExtractor struct {
	image image.Image
	path  string
}

type goserractConfig struct {
	Image     string
	Whitelist string
	Blacklist string
	Language  string
}

type imageExtractConfig struct {
	x0            int
	y0            int
	x1            int
	y1            int
	CleanAllSpace bool
}

const langEng = "eng"
const langJpn = "jpn"

const (
	MoneyScore = iota
	PlayDate
	PlayTime
	TitleEng
	TitleJpn
	AchievementRate
)

func New(imagePath string) (ImageExtractor, error) {
	//open image
	rawImage, err := imaging.Open(imagePath)
	if err != nil {
		log.Fatal("Error when loading image", err)
	}

	validErr := validateImage(rawImage)
	if validErr != nil {
		return ImageExtractor{}, validErr
	}

	//get image size
	width := rawImage.Bounds().Dx()
	height := rawImage.Bounds().Dy()

	//double image size and increase the contrast
	rawImage = imaging.Resize(rawImage, width*2, height*2, imaging.Lanczos)
	rawImage = imaging.AdjustContrast(rawImage, 100)
	rawImage = imaging.AdjustBrightness(rawImage, -80)

	return ImageExtractor{
		image: rawImage,
		path:  imagePath,
	}, nil
}

func (ie *ImageExtractor) Extractor(extractType int) (string, error) {
	gosseractConfig, imageConfig, err := ie.extractSelector(extractType)
	if err != nil {
		fmt.Println("Error Loading Config:", err)
		return "", err
	}

	//crop image
	image := imaging.Crop(ie.image, image.Rect(imageConfig.x0, imageConfig.y0, imageConfig.x1, imageConfig.y1))

	//generate image to path
	out, err := os.Create(gosseractConfig.Image)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	jpeg.Encode(out, image, nil)
	defer func() {
		err := os.Remove(gosseractConfig.Image)
		if err != nil {
			log.Printf("Error in deleting file", gosseractConfig.Image, ":", err)
		}
	}()

	//goserract
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(gosseractConfig.Image)
	client.SetLanguage(gosseractConfig.Language)
	if len(gosseractConfig.Whitelist) > 0 {
		client.SetWhitelist(gosseractConfig.Whitelist)
	}
	if len(gosseractConfig.Blacklist) > 0 {
		client.SetBlacklist(gosseractConfig.Blacklist)
	}

	text, _ := client.Text()
	if imageConfig.CleanAllSpace {
		text = strings.Replace(text, " ", "", -1)
	}

	return text, nil
}

func (ie *ImageExtractor) extractSelector(extractType int) (goserractConfig, imageExtractConfig, error) {
	switch extractType {
	case MoneyScore:
		return goserractConfig{
				Image:     ie.path + "_score",
				Whitelist: "0123456789",
				Language:  langEng,
			},
			imageExtractConfig{
				x0: 640,
				y0: 450,
				x1: 740,
				y1: 480,
			},
			nil
	case PlayDate:
		return goserractConfig{
				Image:     ie.path + "_playdate",
				Whitelist: "0123456789-",
				Language:  langEng,
			},
			imageExtractConfig{
				x0:            940,
				y0:            40,
				x1:            1100,
				y1:            70,
				CleanAllSpace: true,
			},
			nil
	case PlayTime:
		return goserractConfig{
				Image:     ie.path + "_playtime",
				Whitelist: "0123456789:",
				Language:  langEng,
			},
			imageExtractConfig{
				x0:            1100,
				y0:            40,
				x1:            1200,
				y1:            70,
				CleanAllSpace: true,
			},
			nil
	case TitleEng:
		return goserractConfig{
				Image:     ie.path + "_title",
				Whitelist: "",
				Language:  langEng,
			},
			imageExtractConfig{
				x0: 180,
				y0: 300,
				x1: 1200,
				y1: 340,
			},
			nil

	case TitleJpn:
		return goserractConfig{
				Image:     ie.path + "_title",
				Whitelist: "",
				Language:  langJpn,
			},
			imageExtractConfig{
				x0: 180,
				y0: 300,
				x1: 1200,
				y1: 340,
			},
			nil
	case AchievementRate:
		return goserractConfig{
				Image:     ie.path + "_rate",
				Whitelist: "0123456789.",
				Language:  langEng,
			},
			imageExtractConfig{
				x0:            620,
				y0:            790,
				x1:            720,
				y1:            820,
				CleanAllSpace: true,
			},
			nil
	default:
		return goserractConfig{}, imageExtractConfig{}, fmt.Errorf("Extract type not found")
	}
}

func validateImage(input image.Image) error {
	if input.Bounds().Dx() != 600 || input.Bounds().Dy() != 480 {
		return fmt.Errorf("seems like the picture is not taken from eamuse app")
	}
	return nil
}
