package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type TemperatureScope struct {
	Lo int
	Hi int
}

type Sun struct {
	Rise string
	Set  string
}

type WeatherInfo struct {
	CityEN             string
	CityCN             string
	Date               string
	Week               string
	MoonCalendar       string
	CurrentTemperature int
	Weather            string
	WeatherPictureUrl  string
	TemperatureScope   TemperatureScope
	AQ                 string
	AQI                int
	Sun                Sun
}

func (this *WeatherInfo) FillCity(selection *goquery.Selection) {
	this.CityCN = selection.Find(".name").Find("h2").Text()
}

func (this *WeatherInfo) FillDate(selection *goquery.Selection) {
	dateInfo := selection.Find(".week").Text()
	infoes := strings.Split(dateInfo, "　")
	this.Date = infoes[0]
	this.Week = infoes[1]
	this.MoonCalendar = strings.TrimSpace(infoes[2])
}

func (this *WeatherInfo) FillWeather(selection *goquery.Selection) {
	dateInfo := selection.Find(".weather")
	this.WeatherPictureUrl = dateInfo.Find("i").Find("img").AttrOr("src", "")
	this.CurrentTemperature, _ = strconv.Atoi(dateInfo.Find(".now").Find("b").Text())
	this.Weather = dateInfo.Find("span").Find("b").Text()
	scopeString := dateInfo.Find("span").Text()
	scopeString = strings.TrimLeft(scopeString, this.Weather)
	scopeString = strings.TrimRight(scopeString, "℃")
	scope := strings.Split(scopeString, " ~ ")
	lo, _ := strconv.Atoi(scope[0])
	hi, _ := strconv.Atoi(scope[1])
	this.TemperatureScope = TemperatureScope{
		Lo: lo,
		Hi: hi,
	}
}

func (this *WeatherInfo) FillAQI(selection *goquery.Selection) {
	aqiInfo := selection.Find(".kongqi")
	this.AQ = strings.TrimLeft(aqiInfo.Find("h5").Text(), "空气质量：")
	this.AQI, _ = strconv.Atoi(strings.TrimLeft(aqiInfo.Find("h6").Text(), "PM: "))
	sunInfo := strings.Split(aqiInfo.Find("span").Text(), "日落: ")
	sunRise := strings.TrimLeft(sunInfo[0], "日出: ")
	sunSet := sunInfo[1]
	this.Sun = Sun{
		Rise: sunRise,
		Set:  sunSet,
	}
}

func GetWeatherWebPage(city string) WeatherInfo {
	url := fmt.Sprintf("https://www.tianqi.com/%s/", city)
	//var client http.Client

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	resp, err := client.Do(req)
	//resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var weatherInfo WeatherInfo
	weatherInfo.CityEN = city
	if resp.StatusCode == http.StatusOK {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		wi := doc.Find(".weather_info")
		weatherInfo.FillAQI(wi)
		weatherInfo.FillCity(wi)
		weatherInfo.FillDate(wi)
		weatherInfo.FillWeather(wi)
	}
	//bs, _ := json.Marshal(weatherInfo)
	//fmt.Println(string(bs))
	return weatherInfo
}
