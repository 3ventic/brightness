package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/mqu/go-notify"
)

const (
	brightnessPath    = "/sys/class/backlight/intel_backlight/brightness"
	maxBrightnessPath = "/sys/class/backlight/intel_backlight/max_brightness"
)

func main() {
	incFlag := flag.Bool("inc", false, "increase brightness by 2%")
	decFlag := flag.Bool("dec", false, "decrease brightness by 2%")
	notifyFlag := flag.Bool("notify", false, "send a desktop notification with new brightness")
	flag.Parse()

	inc := incFlag != nil && *incFlag
	dec := decFlag != nil && *decFlag
	desktopNotify := notifyFlag != nil && *notifyFlag
	oldBrightness := readInt64(brightnessPath)
	maxBrightness := readInt64(maxBrightnessPath)
	if !inc && !dec {
		fmt.Println(percentageString(float64(oldBrightness), float64(maxBrightness)))
		return
	} else if inc && dec {
		panic("both -inc and -dec specified. Cannot determine what to do")
	}

	step := maxBrightness / 50 // 2% steps
	min := step
	if dec {
		step *= -1
	}
	newBrightness := oldBrightness + step
	if newBrightness < min {
		newBrightness = min
	}
	if newBrightness > maxBrightness {
		newBrightness = maxBrightness
	}
	ps := percentageString(float64(newBrightness), float64(maxBrightness))
	fmt.Println(ps)
	if desktopNotify && notify.Init("Brightness") {
		n := notify.NotificationNew("Brightness", fmt.Sprintf("Changed: %s", ps), "")
		if n == nil {
			panic("unable to create notification")
		}
		notify.NotificationSetTimeout(n, 1500)
		err := notify.NotificationShow(n)
		if err != nil {
			panic(err)
		}

		notify.NotificationClose(n)
		notify.UnInit()
	}
	updateBrightness(newBrightness)
}

func updateBrightness(i int64) {
	s := strconv.FormatInt(i, 10)
	b := []byte(fmt.Sprintf("%s\n", s))
	err := ioutil.WriteFile(brightnessPath, b, 0644)
	if err != nil {
		panic(fmt.Errorf("%s write: %v", brightnessPath, err))
	}
}

func readInt64(path string) int64 {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("%s read: %v", path, err))
	}

	i, err := strconv.ParseInt(strings.Trim(string(d), " \r\n"), 10, 64)
	if err != nil {
		panic(fmt.Errorf("%s parse: %v", path, err))
	}
	return i
}

func percentageString(current, max float64) string {
	return fmt.Sprintf("%.2f%%\n", 100*current/max)
}
