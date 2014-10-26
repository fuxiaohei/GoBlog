package gof

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type ConfigInterface interface {
	Set(key string, value interface{})
	String(key string, def ...string) string
	Int(key string, def ...int64) int64
	Float(key string, def ...float64) float64
	Bool(key string) bool
	ToFile(file string) error
	FromFile(file string) error
}

var _ ConfigInterface = new(Config)

type Config struct {
	data map[string]interface{}
}

func NewConfig(file string) (*Config, error) {
	c := new(Config)
	c.data = make(map[string]interface{})
	if file == "" {
		return c, nil
	}
	err := c.FromFile(file)
	return c, err
}

func (c *Config) FromFile(file string) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	dataMap := make(map[string]interface{})
	if err = json.Unmarshal(bytes, &dataMap); err != nil {
		return err
	}
	// add to current data map, not re-assign
	for k, v := range dataMap {
		c.data[k] = v
	}
	return nil
}

func (c *Config) ToFile(file string) error {
	bytes, err := json.MarshalIndent(c.data, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, bytes, os.ModePerm)
	return err
}

func (c *Config) Set(key string, value interface{}) {
	mapPtr, keyName := c.getDeepMap(key)
	(*mapPtr)[keyName] = value
}

func (c *Config) getDeepMap(key string) (*map[string]interface{}, string) {
	keySlice := strings.Split(key, ".")
	keyName := keySlice[len(keySlice)-1]
	var mapPtr *map[string]interface{}
	for _, k := range keySlice[:len(keySlice)-1] {
		if mapPtr == nil {
			if _, ok := c.data[k]; !ok {
				c.data[k] = make(map[string]interface{})
			}
			m := c.data[k].(map[string]interface{})
			mapPtr = &m
			continue
		}
		m := *mapPtr
		if _, ok := m[k]; !ok {
			m[k] = make(map[string]interface{})
		}
		m = m[k].(map[string]interface{})
		mapPtr = &m
	}
	return mapPtr, keyName
}

func (c *Config) String(key string, def ...string) string {
	mapPtr, keyName := c.getDeepMap(key)
	v := fmt.Sprint((*mapPtr)[keyName])
	if v == "" && len(def) > 0 {
		return def[0]
	}
	return v
}

func (c *Config) Int(key string, def ...int64) int64 {
	s := c.String(key)
	i, _ := strconv.ParseInt(s, 10, 64)
	if i == 0 && len(def) > 0 {
		return def[0]
	}
	return i
}

func (c *Config) Float(key string, def ...float64) float64 {
	s := c.String(key)
	f, _ := strconv.ParseFloat(s, 10)
	if f == 0 && len(def) > 0 {
		return def[0]
	}
	return f
}

func (c *Config) Bool(key string) bool {
	s := c.String(key)
	b, _ := strconv.ParseBool(s)
	return b
}
