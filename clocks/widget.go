package clocks

import (
	"sort"
	"time"

	"github.com/olebedev/config"
	"github.com/senorprogrammer/wtf/wtf"
)

const TimeFormat = "15:04 MST"
const DateFormat = "Jan 2"

// Config is a pointer to the global config object
var Config *config.Config

type Widget struct {
	wtf.TextWidget
}

func NewWidget() *Widget {
	widget := Widget{
		TextWidget: wtf.NewTextWidget(" 🕗 World Clocks ", "clocks"),
	}

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	if widget.Disabled() {
		return
	}

	widget.View.Clear()
	widget.display()
	widget.RefreshedAt = time.Now()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) colorFor(idx int) string {
	rowColor := Config.UString("wtf.mods.clocks.rowcolors.even", "lightblue")

	if idx%2 == 0 {
		rowColor = Config.UString("wtf.mods.clocks.rowcolors.odd", "white")
	}

	return rowColor
}

func (widget *Widget) locations(locs map[string]interface{}) map[string]time.Time {
	times := make(map[string]time.Time)

	for label, loc := range locs {
		tzloc, err := time.LoadLocation(loc.(string))

		if err != nil {
			continue
		}

		times[label] = time.Now().In(tzloc)
	}

	return times
}

func (widget *Widget) sortedLabels(locs map[string]time.Time) []string {
	labels := []string{}

	for label, _ := range locs {
		labels = append(labels, label)
	}

	sort.Strings(labels)

	return labels
}
