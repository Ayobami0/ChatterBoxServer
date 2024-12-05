package server

import "fmt"

type AppConfig struct {
	databases map[string]interface{}
}

func NewConfig(databases map[string]interface{}) AppConfig {
	return AppConfig{
		databases: databases,
	}
}

func (a AppConfig) DB(name string) (interface{}, error) {
	db, ok := a.databases[name]

	if !ok {
		return nil, fmt.Errorf("database '%s' not registered", name)
	}
	return db, nil
}
