package rpg

import (
	"bufio"
	"os"
	"strings"
)

// Configuration can hold a collection of values, useful to setup
// starting situations or other system uses.
type Configuration struct {
	options map[string]interface{}
}

// NewConfiguration will create a new configuration and
// initialize it.
func NewConfiguration() *Configuration {
	return &Configuration{options: make(map[string]interface{})}
}

// NewSystemConfiguration will start the game system configuration.
func NewSystemConfiguration(file string) {
	Config = NewConfiguration()
	Config.Load(file, false)
}

// Load will open up the requested configuration file and read
// these options into the config. Option to append or overwrite.
func (c *Configuration) Load(file string, appendScript bool) error {
	// Open our file, return on error
	configfile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer configfile.Close()

	// If append is false, we need to clear the script.
	if !appendScript {
		c.options = make(map[string]interface{})
	}

	// Create our reader
	configReader := bufio.NewReader(configfile)

	// And now our scanner
	configScanner := bufio.NewScanner(configReader)

	configLine := ""
	// Loop while we are scanning
	for configScanner.Scan() {
		// Get our next line
		configLine = configScanner.Text()

		// Split our line up
		results := strings.Split(configLine, ":")

		results[0] = strings.TrimSpace(results[0])
		results[1] = strings.TrimSpace(results[1])

		// We should only have 2 fields for now, so should be easy.
		c.options[results[0]] = results[1]
	}

	return nil
}

// Set will set an option in our configuration.
func (c *Configuration) Set(s string, v interface{}) {
	c.options[s] = v
}

// Get will return an option from our configuration. Use this
// and be prepared to typecast otherwise it can get ugly. Generally
// using the premade typed return values is better. If you are
// requesting a configuration value, you should know what
// you are expecting.
func (c *Configuration) Get(s string) interface{} {
	return c.options[s]
}

// Value will return an option as an int.
func (c *Configuration) Value(s string) int {
	return c.Get(s).(int)
}

// BigValue will return an option as an int64.
func (c *Configuration) BigValue(s string) int64 {
	return int64(c.Get(s).(int))
}

// String will return an option as a string.
func (c *Configuration) String(s string) string {
	return c.Get(s).(string)
}

// Float will return an option as a float64.
func (c *Configuration) Float(s string) float64 {
	return StrFloat(c.Get(s))
}

// Bool will return an option as a bool.
func (c *Configuration) Bool(s string) bool {
	return StrBool(c.Get(s))
}
