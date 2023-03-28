package ifunny

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

const (
	CHANNEL_MESSAGE eventType = 200
)

type eventType int

type Event interface {
	Type() eventType
	Decode(target interface{}) error
}

func makeEvent(eType int, data map[string]interface{}) Event {
	switch eventType(eType) {
	case CHANNEL_MESSAGE:
		return keySerialize{CHANNEL_MESSAGE, data, "message"}
	default:
		return noSerialize{eventType(eType), data}
	}
}

type keySerialize struct {
	eType eventType
	data  map[string]interface{}
	key   string
}

func (key keySerialize) Type() eventType {
	return key.eType
}

func (key keySerialize) Decode(target interface{}) error {
	if key.data[key.key] == nil {
		return fmt.Errorf("key %s doesn't exist on data type %d", key.key, key.eType)
	}

	return mapstructure.Decode(key.data[key.key], target)
}

type noSerialize struct {
	eType eventType
	data  map[string]interface{}
}

func (no noSerialize) Type() eventType {
	return no.eType
}

func (no noSerialize) Decode(interface{}) error {
	return fmt.Errorf("no decoder for data type %d: %+v", no.eType, no.data)
}
