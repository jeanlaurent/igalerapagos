package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWeather(t *testing.T) {
	weather := newWeather()
	assert.Equal(t, sunny, weather.weather)
	assert.Equal(t, 0, weather.nbOfDays)
}

func TestWeatherChangeFromSunnyToStormy(t *testing.T) {
	weather := newWeather()
	weather.applyWeatherChange(0)
	assert.Equal(t, stormy, weather.weather)
	assert.Equal(t, 0, weather.nbOfDays)
}

func TestWeatherChangeFromSunnyToCloudy(t *testing.T) {
	weather := newWeather()
	weather.applyWeatherChange(100)
	assert.Equal(t, cloudy, weather.weather)
	assert.Equal(t, 0, weather.nbOfDays)
}
func TestWeatherChangeFromCloudyToRainy(t *testing.T) {
	weather := newWeather()
	weather.weather = cloudy
	weather.applyWeatherChange(100)
	assert.Equal(t, rainy, weather.weather)
	assert.Equal(t, 0, weather.nbOfDays)
}

func TestWeatherChangeFromStormyToSunny(t *testing.T) {
	weather := newWeather()
	weather.weather = stormy
	weather.applyWeatherChange(100)
	assert.Equal(t, sunny, weather.weather)
	assert.Equal(t, 0, weather.nbOfDays)
}

func TestWeatherString(t *testing.T) {
	weather := newWeather()
	assert.Equal(t, weather.weatherAsString(), "sunny")
	weather.weather = cloudy
	assert.Equal(t, weather.weatherAsString(), "cloudy")
	weather.weather = rainy
	assert.Equal(t, weather.weatherAsString(), "rainy")
	weather.weather = stormy
	assert.Equal(t, weather.weatherAsString(), "stormy")
	weather.weather = -1
	assert.Equal(t, weather.weatherAsString(), "unknown")
}

func TestStormyWeatherHinderFoodGathering(t *testing.T) {
	weather := newWeather()
	weather.weather = stormy
	assert.Less(t, weather.foodGatheringBonus(), 0)
}

func TestSunnyWeatherBoostFoodGathering(t *testing.T) {
	weather := newWeather()
	weather.weather = sunny
	assert.Greater(t, weather.foodGatheringBonus(), 0)
}

func TestCloudyAndRainyWeatherDoNotBoostFoodGathering(t *testing.T) {
	weather := newWeather()
	weather.weather = cloudy
	assert.Equal(t, weather.foodGatheringBonus(), 0)
	weather.weather = rainy
	assert.Equal(t, weather.foodGatheringBonus(), 0)
}
