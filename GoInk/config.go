package GoInk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Config map[string]map[string]interface{}

func (cfg *Config) String(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return ""
	}
	str, ok := (*cfg)[keys[0]][keys[1]]
	if !ok {
		return ""
	}
	return fmt.Sprint(str)
}

func (cfg *Config) StringOr(key string, def string) string {
	value := cfg.String(key)
	if value == "" {
		cfg.Set(key, def)
		return def
	}
	return value
}

func (cfg *Config) Int(key string) int {
	str := cfg.String(key)
	i, _ := strconv.Atoi(str)
	return i
}

func (cfg *Config) IntOr(key string, def int) int {
	i := cfg.Int(key)
	if i == 0 {
		cfg.Set(key, def)
		return def
	}
	return i
}

func (cfg *Config) Float(key string) float64 {
	str := cfg.String(key)
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

func (cfg *Config) FloatOr(key string, def float64) float64 {
	f := cfg.Float(key)
	if f == 0.0 {
		cfg.Set(key, def)
		return def
	}
	return f
}

func (cfg *Config) Bool(key string) bool {
	str := cfg.String(key)
	b, _ := strconv.ParseBool(str)
	return b
}

func (cfg *Config) Set(key string, value interface{}) {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return
	}
	if (*cfg) == nil {
		(*cfg) = make(map[string]map[string]interface{})
	}
	if _, ok := (*cfg)[keys[0]]; !ok {
		(*cfg)[keys[0]] = make(map[string]interface{})
	}
	(*cfg)[keys[0]][keys[1]] = value
}

func NewConfig(fileAbsPath string) (*Config, error) {
	cfg := new(Config)
	bytes, e := ioutil.ReadFile(fileAbsPath)
	if e != nil {
		return cfg, e
	}
	e = json.Unmarshal(bytes, cfg)
	return cfg, e
}
