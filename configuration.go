package main

import (
	"github.com/emicklei/goproperties"
	//"github.com/emicklei/hopwatch"
	"errors"
	"strings"
)

var configurationMap = map[string]properties.Properties{}

func initConfiguration(props properties.Properties) {
	aliases := props.SelectProperties("mongod.*")
	for k, v := range aliases {
		parts := strings.Split(k, ".")
		alias := parts[1]
		config := configurationMap[alias]
		if config == nil {
			config = properties.Properties{}
			config["alias"] = alias
			configurationMap[alias] = config
		}
		config[parts[2]] = v
	}
	//hopwatch.Dump(configurationMap)
}

func configuration(alias string) (properties.Properties, error) {
	config := configurationMap[alias]
	if config == nil {
		return nil, errors.New("Unknown alias:" + alias)
	} else {
		return config, nil
	}
}
