package config

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

func extractStringValue(key string) string {
	checkPresenceOf(key)
	return viper.GetString(key)
}

func extractIntValue(key string) int {
	checkPresenceOf(key)
	v, err := strconv.Atoi(viper.GetString(key))
	if err != nil {
		panic(fmt.Sprintf("key %s is not a valid Integer value", key))
	}

	return v
}

func checkPresenceOf(key string) {
	if !viper.IsSet(key) {
		panic(fmt.Sprintf("key %s is not set", key))
	}
}
