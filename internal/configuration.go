package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DefaultConfigFileName = "defaults.yaml"
	appDataVar            = "APPDATA"
	fkKey                 = "key"
	fkType                = "type"
	fkValue               = "value"
)

// ReadConfigurationFile reads defaults.yaml from the specified path and returns
// a pointer to a cooked Configuration instance
func ReadConfigurationFile(o OutputBus) (c *Configuration, ok bool) {
	var appDataValue string
	var appDataSet bool
	if appDataValue, appDataSet = appData(o); !appDataSet {
		c = EmptyConfiguration()
		ok = true
		return
	}
	path := CreateAppSpecificPath(appDataValue)
	configFile := filepath.Join(path, DefaultConfigFileName)
	var err error
	var exists bool
	if exists, err = verifyFileExists(o, configFile); err != nil {
		return
	}
	if !exists {
		c = EmptyConfiguration()
		ok = true
		return
	}
	yfile, _ := ioutil.ReadFile(configFile) // only probable error circumvented by verifyFileExists failure
	data := make(map[string]interface{})
	if err = yaml.Unmarshal(yfile, &data); err != nil {
		o.LogWriter().Warn(LW_CANNOT_UNMARSHAL_YAML, map[string]interface{}{
			FK_DIRECTORY: path,
			FK_FILE_NAME: DefaultConfigFileName,
			FK_ERROR:     err,
		})
		o.WriteError(USER_CONFIGURATION_FILE_GARBLED, configFile, err)
		return
	}
	c = CreateConfiguration(o, data)
	ok = true
	o.LogWriter().Info(LI_CONFIGURATION_FILE_READ, map[string]interface{}{
		FK_DIRECTORY: path,
		FK_FILE_NAME: DefaultConfigFileName,
		fkValue:      c,
	})
	return
}

func appData(o OutputBus) (string, bool) {
	if value, ok := os.LookupEnv(appDataVar); ok {
		return value, ok
	}
	o.LogWriter().Info(LI_NOT_SET, map[string]interface{}{
		fkVarName: appDataVar,
	})
	return "", false
}

func verifyFileExists(o OutputBus, path string) (ok bool, err error) {
	f, err := os.Stat(path)
	if err == nil {
		if f.IsDir() {
			o.LogWriter().Error(LE_FILE_IS_DIR, map[string]interface{}{
				FK_DIRECTORY: filepath.Dir(path),
				FK_FILE_NAME: filepath.Base(path),
			})
			o.WriteError(USER_CONFIGURATION_FILE_IS_DIR, path)
			err = fmt.Errorf(ERROR_FILE_IS_DIR)
			return
		}
		ok = true
		return
	}
	if errors.Is(err, os.ErrNotExist) {
		o.LogWriter().Info(LI_NO_SUCH_FILE, map[string]interface{}{
			FK_DIRECTORY: filepath.Dir(path),
			FK_FILE_NAME: filepath.Base(path),
		})
		err = nil
	}
	return
}

func (c *Configuration) String() string {
	var output []string
	if len(c.bMap) != 0 {
		output = append(output, fmt.Sprintf("%v", c.bMap))
	}
	if len(c.iMap) != 0 {
		output = append(output, fmt.Sprintf("%v", c.iMap))
	}
	if len(c.sMap) != 0 {
		output = append(output, fmt.Sprintf("%v", c.sMap))
	}
	if len(c.cMap) != 0 {
		output = append(output, fmt.Sprintf("%v", c.cMap))
	}
	return strings.Join(output, ", ")
}

func EmptyConfiguration() *Configuration {
	return &Configuration{
		bMap: make(map[string]bool),
		iMap: make(map[string]int),
		sMap: make(map[string]string),
		cMap: make(map[string]*Configuration),
	}
}

// SubConfiguration returns a specified sub-configuration
func (c *Configuration) SubConfiguration(key string) *Configuration {
	if configuration, ok := c.cMap[key]; ok {
		return configuration
	}
	return EmptyConfiguration()
}

// BoolDefault returns a boolean value for a specified key
func (c *Configuration) BoolDefault(key string, defaultValue bool) (b bool, err error) {
	b = defaultValue
	if value, ok := c.bMap[key]; ok {
		b = value
	} else {
		if value, ok := c.iMap[key]; ok {
			switch value {
			case 0:
				b = false
			case 1:
				b = true
			default:
				err = fmt.Errorf("invalid boolean value \"%d\" for -%s: parse error", value, key)
			}
		} else {
			// True values may be specified as "t", "T", "true", "TRUE", or "True"
			// False values may be specified as "f", "F", "false", "FALSE", or "False"_."
			if value, ok := c.sMap[key]; ok {
				rawValue := InterpretEnvVarReferences(value)
				if cookedValue, e := strconv.ParseBool(rawValue); e == nil {
					b = cookedValue
				} else {
					err = fmt.Errorf("invalid boolean value %q for -%s: parse error", value, key)
				}
			}
		}
	}
	return
}

type IntBounds struct {
	minValue     int
	defaultValue int
	maxValue     int
}

func NewIntBounds(v1, v2, v3 int) *IntBounds {
	is := []int{v1, v2, v3}
	sort.Ints(is)
	return &IntBounds{
		minValue:     is[0],
		defaultValue: is[1],
		maxValue:     is[2],
	}
}

func (c *Configuration) IntDefault(key string, sortedBounds *IntBounds) (i int) {
	i = sortedBounds.defaultValue
	if value, ok := c.iMap[key]; ok {
		if value < sortedBounds.minValue {
			i = sortedBounds.minValue
		} else if value > sortedBounds.maxValue {
			i = sortedBounds.maxValue
		} else {
			i = value
		}
	} else {
		if value, ok := c.sMap[key]; ok {
			rawValue := InterpretEnvVarReferences(value)
			if cookedValue, e := strconv.Atoi(rawValue); e == nil {
				if cookedValue < sortedBounds.minValue {
					i = sortedBounds.minValue
				} else if cookedValue > sortedBounds.maxValue {
					i = sortedBounds.maxValue
				} else {
					i = cookedValue
				}
			} // else 'invalid value "flopsy" for flag -timeout: parse error'
		}
	}
	return
}

// StringDefault returns a string value for a specified key
func (c *Configuration) StringDefault(key string, defaultValue string) (s string) {
	s = InterpretEnvVarReferences(defaultValue)
	if value, ok := c.sMap[key]; ok {
		s = InterpretEnvVarReferences(value)
	}
	return
}

// StringValue returns the definition of the specified key and ok if the value
// is defined
func (c *Configuration) StringValue(key string) (value string, ok bool) {
	value, ok = c.sMap[key]
	return
}

type Configuration struct {
	sMap map[string]string
	bMap map[string]bool
	iMap map[string]int
	cMap map[string]*Configuration
}

func CreateConfiguration(o OutputBus, data map[string]interface{}) *Configuration {
	c := EmptyConfiguration()
	for key, v := range data {
		switch t := v.(type) {
		case string:
			c.sMap[key] = t
		case bool:
			c.bMap[key] = t
		case int:
			c.iMap[key] = t
		case map[string]interface{}:
			c.cMap[key] = CreateConfiguration(o, t)
		default:
			o.LogWriter().Warn(LW_UNEXPECTED_VALUE_TYPE, map[string]interface{}{
				fkKey:   key,
				fkValue: v,
				fkType:  fmt.Sprintf("%T", v),
			})
			c.sMap[key] = fmt.Sprintf("%v", v)
		}
	}
	return c
}
