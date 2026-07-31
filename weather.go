package main

// WeatherState represents the current weather condition.
type WeatherState int

const (
	sunny       WeatherState = iota // 0
	cloudy                          // 1
	rainy                           // 2
	stormy                          // 3
	numWeathers                     // 4 — sentinel used for modular wrap-around
)

// String implements fmt.Stringer for readable output and debugging.
func (w WeatherState) String() string {
	switch w {
	case sunny:
		return "sunny"
	case cloudy:
		return "cloudy"
	case rainy:
		return "rainy"
	case stormy:
		return "stormy"
	}
	return "unknown"
}

type weather struct {
	state    WeatherState
	nbOfDays int
}

func (w *weather) changeWeather(d Dice) {
	w.applyWeatherChange(d.roll(100))
}

func (w *weather) applyWeatherChange(roll int) {
	if roll <= 25 {
		w.state = (w.state - 1 + numWeathers) % numWeathers
		w.nbOfDays = 0
	} else if roll >= 75 {
		w.state = (w.state + 1 + numWeathers) % numWeathers
		w.nbOfDays = 0
	} else {
		w.nbOfDays++
	}
}

func (w *weather) foodGatheringBonus() int {
	switch w.state {
	case sunny:
		return 2
	case stormy:
		return -2
	}
	return 0
}

func newWeather() weather {
	return weather{state: sunny, nbOfDays: 0}
}
