package types

import (
	"github.com/MaxPolarfox/goTools/mongoDB"
)

type Options struct {
	Port        int         `json:"port"`
	ServiceName string      `json:"serviceName"`
	DB          Collections `json:"db"`
}

type Collections struct {
	Blog mongoDB.Options `json:"tasks"`
}
