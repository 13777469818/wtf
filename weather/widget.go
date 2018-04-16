package weather

import (
	"fmt"
	"strings"
	"time"

	owm "github.com/briandowns/openweathermap"
	"github.com/olebedev/config"
	"github.com/senorprogrammer/wtf/wtf"
)

var Config *config.Config

type Widget struct {
	wtf.TextWidget

	Current int
}

func NewWidget() *Widget {
	widget := Widget{
		TextWidget: wtf.NewTextWidget(" Weather ", "weather"),
		Current:    0,
	}

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	if widget.Disabled() {
		return
	}

	data := Fetch(wtf.ToInts(Config.UList("wtf.mods.weather.cityids", widget.defaultCityCodes())))

	widget.View.Clear()
	widget.contentFrom(data)
	widget.RefreshedAt = time.Now()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) contentFrom(data []*owm.CurrentWeatherData) {
	cityData := widget.currentCityData(data)

	if len(cityData.Weather) == 0 {
		fmt.Fprintf(widget.View, "%s", " Weather data is unavailable.")
	}

	widget.View.SetTitle(widget.contentTitle(cityData))

	str := widget.contentTickMarks(data) + "\n"
	str = str + widget.contentDescription(cityData) + "\n\n"
	str = str + widget.contentTemperatures(cityData) + "\n"
	str = str + widget.contentSunInfo(cityData)

	fmt.Fprintf(widget.View, "%s", str)
}

func (widget *Widget) contentTickMarks(data []*owm.CurrentWeatherData) string {
	str := ""

	if len(data) > 1 {
		tickMarks := strings.Repeat("*", len(data))
		tickMarks = tickMarks[:widget.Current] + "_" + tickMarks[widget.Current+1:]

		str = "[lightblue]" + fmt.Sprintf(wtf.RightAlignFormat(widget.View), tickMarks) + "[white]"
	}

	return str
}

func (widget *Widget) contentTitle(cityData *owm.CurrentWeatherData) string {
	return fmt.Sprintf(" %s %s ", widget.icon(cityData), cityData.Name)
}

func (widget *Widget) contentDescription(cityData *owm.CurrentWeatherData) string {
	descs := []string{}
	for _, weather := range cityData.Weather {
		descs = append(descs, fmt.Sprintf(" %s", weather.Description))
	}

	return strings.Join(descs, ",")
}

func (widget *Widget) contentTemperatures(cityData *owm.CurrentWeatherData) string {
	tempUnit := Config.UString("wtf.mods.weather.tempUnit", "C")

	str := fmt.Sprintf("%8s: %4.1f° %s\n", "High", cityData.Main.TempMax, tempUnit)
	str = str + fmt.Sprintf("%8s: [green]%4.1f° %s[white]\n", "Current", cityData.Main.Temp, tempUnit)
	str = str + fmt.Sprintf("%8s: %4.1f° %s\n", "Low", cityData.Main.TempMin, tempUnit)

	return str
}

func (widget *Widget) contentSunInfo(cityData *owm.CurrentWeatherData) string {
	return fmt.Sprintf(
		" Rise: %s    Set: %s\n",
		wtf.UnixTime(int64(cityData.Sys.Sunrise)).Format("15:04 MST"),
		wtf.UnixTime(int64(cityData.Sys.Sunset)).Format("15:04 MST"),
	)
}

func (widget *Widget) currentCityData(data []*owm.CurrentWeatherData) *owm.CurrentWeatherData {
	return data[widget.Current]
}

func (widget *Widget) defaultCityCodes() []interface{} {
	defaultArr := []int{6176823, 360630, 3413829}

	var defaults []interface{} = make([]interface{}, len(defaultArr))
	for i, d := range defaultArr {
		defaults[i] = d
	}

	return defaults
}

// icon returns an emoji for the current weather
// src: https://github.com/chubin/wttr.in/blob/master/share/translations/en.txt
// Note: these only work for English weather status. Sorry about that
func (widget *Widget) icon(data *owm.CurrentWeatherData) string {
	var icon string

	if len(data.Weather) == 0 {
		return ""
	}

	switch data.Weather[0].Description {
	case "broken clouds":
		icon = "☁️"
	case "clear":
		icon = "☀️"
	case "clear sky":
		icon = "☀️ "
	case "cloudy":
		icon = "⛅️"
	case "few clouds":
		icon = "🌤"
	case "fog":
		icon = "🌫"
	case "haze":
		icon = "🌫"
	case "heavy rain":
		icon = "💦"
	case "heavy snow":
		icon = "⛄️"
	case "light intensity shower rain":
		icon = "☔️"
	case "light rain":
		icon = "🌦"
	case "light snow":
		icon = "🌨"
	case "mist":
		icon = "🌬"
	case "moderate rain":
		icon = "🌧"
	case "moderate snow":
		icon = "🌨"
	case "overcast":
		icon = "🌥"
	case "overcast clouds":
		icon = "🌥"
	case "partly cloudy":
		icon = "🌤"
	case "scattered clouds":
		icon = "☁️"
	case "shower rain":
		icon = "☔️"
	case "snow":
		icon = "❄️"
	case "sunny":
		icon = "☀️"
	default:
		icon = "💥"
	}

	return icon
}

func (widget *Widget) refreshedAt() string {
	return widget.RefreshedAt.Format("15:04:05")
}
