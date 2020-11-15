package gamesys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func LoadValidConfig() *Configuration {
	// Loads up and returns a known valid configuration for testing.
	newconfig := NewConfiguration()
	newconfig.Load("test_assets/system.config", false)
	return newconfig
}

func TestNewConfiguration(t *testing.T) {
	// Check that our options are present
	newconfig := NewConfiguration()
	assert.NotNil(t, newconfig)
	assert.NotNil(t, newconfig.options)
}

func TestNewSystemConfiguration(t *testing.T) {
	// I think we can set our vars from here and check them.
	NewSystemConfiguration("test_assets/system.config")

	// Check our configuration
	assert.NotNil(t, Config)
	assert.NotNil(t, Config.options)
}

func TestLoadConfiguration(t *testing.T) {
	// Our valid info
	newconfig := LoadValidConfig()

	// Ensure it's there
	assert.NotNil(t, newconfig)
	assert.NotNil(t, newconfig.options)

	// Check we loaded ALL options
	assert.Len(t, newconfig.options, 7, "We should have 7 options")

	// Length should double if we append.
	newconfig.Load("test_assets/systemextra.config", true)
	assert.Len(t, newconfig.options, 14, "We should have 14 options")

	// Now for a missing file
	newconfig = NewConfiguration()
	err = newconfig.Load("filenotfound", false)
	assert.Error(t, err)

}

func TestConfigurationAccess(t *testing.T) {
	newconfig := NewConfiguration()

	newconfig.Set("myint", 5)
	newconfig.Set("mybigint", int64(5))
	newconfig.Set("mystring", "hello")
	newconfig.Set("myfloat", 5.5)
	newconfig.Set("mytruth", true)
	newconfig.Set("myfalsy", false)

	// This should have set 6 options
	assert.Len(t, newconfig.options, 6, "We should have 6 options")

	// Ensure we can get things the way we expect
	assert.Equal(t, 5, newconfig.Value("myint"), "We are expecting an int")
	assert.Equal(t, int64(5), newconfig.BigValue("mybigint"), "We are expecting a 64bit int")
	assert.Equal(t, "hello", newconfig.String("mystring"), "We are expecting a string")
	assert.Equal(t, 5.5, newconfig.Float("myfloat"), "We are expecting a float64")
	assert.Equal(t, true, newconfig.Bool("mytruth"), "We are expecting a true bool")
	assert.Equal(t, false, newconfig.Bool("myfalsy"), "We are expecting a false bool")

	// What if we break anything?
	assert.NotNil(t, newconfig.String("myint"))

}
