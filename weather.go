package main

const sunny = 0
const cloudy = 1
const rainy = 2
const stormy = 3
const max = 4

type weather struct {
	weather  int
	nbOfDays int
}

func (w *weather) changeWeather(d Dice) {
	w.applyWeatherChange(d.roll(100))
}

func (w *weather) applyWeatherChange(roll int) {
	if roll <= 25 {
		w.weather = (w.weather - 1 + max) % max
		w.nbOfDays = 0
	} else if roll >= 75 {
		w.weather = (w.weather + 1 + max) % max
		w.nbOfDays = 0
	} else {
		w.nbOfDays++
	}
}

func (w *weather) foodGatheringBonus() int {
	switch w.weather {
	case sunny:
		return 2
	case stormy:
		return -2
	}
	return 0
}

func (w *weather) weatherAsString() string {
	switch w.weather {
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

func newWeather() weather {
	return weather{weather: sunny, nbOfDays: 0}
}
