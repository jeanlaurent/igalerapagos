package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWeather(t *testing.T) {
	weather := newWeather()
	assert.Equal(t, sunny, weather.state)
	assert.Equal(t, 0, weather.nbOfDays)
}

func TestWeatherChangeFromSunnyToStormy(t *testing.T) {
	weather := newWeather()
	weather.applyWeatherChange(0)
	assert.Equal(t, stormy, weather.state)
	assert.Equal(t, 0, weather.nbOfDays)
}

func TestWeatherChangeFromSunnyToCloudy(t *testing.T) {
	weather := newWeather()
	weather.applyWeatherChange(100)
	assert.Equal(t, cloudy, weather.state)
	assert.Equal(t, 0, weather.nbOfDays)
}
func TestWeatherChangeFromCloudyToRainy(t *testing.T) {
	weather := newWeather()
	weather.state = cloudy
	weather.applyWeatherChange(100)
	assert.Equal(t, rainy, weather.state)
	assert.Equal(t, 0, weather.nbOfDays)
}

func TestWeatherChangeFromStormyToSunny(t *testing.T) {
	weather := newWeather()
	weather.state = stormy
	weather.applyWeatherChange(100)
	assert.Equal(t, sunny, weather.state)
	assert.Equal(t, 0, weather.nbOfDays)
}

func TestWeatherString(t *testing.T) {
	weather := newWeather()
	assert.Equal(t, "sunny", weather.state.String())
	weather.state = cloudy
	assert.Equal(t, "cloudy", weather.state.String())
	weather.state = rainy
	assert.Equal(t, "rainy", weather.state.String())
	weather.state = stormy
	assert.Equal(t, "stormy", weather.state.String())
	weather.state = -1
	assert.Equal(t, "unknown", weather.state.String())
}

func TestStormyWeatherHinderFoodGathering(t *testing.T) {
	weather := newWeather()
	weather.state = stormy
	assert.Less(t, weather.foodGatheringBonus(), 0)
}

func TestSunnyWeatherBoostFoodGathering(t *testing.T) {
	weather := newWeather()
	weather.state = sunny
	assert.Greater(t, weather.foodGatheringBonus(), 0)
}

func TestCloudyAndRainyWeatherDoNotBoostFoodGathering(t *testing.T) {
	weather := newWeather()
	weather.state = cloudy
	assert.Equal(t, weather.foodGatheringBonus(), 0)
	weather.state = rainy
	assert.Equal(t, weather.foodGatheringBonus(), 0)
}
