package ifunny

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

const (
	CHANNEL_MESSAGE resourceType = 200
)

type resourceType int

type WSResource interface {
	Type() resourceType
	Decode(target interface{}) error
}

func makeEvent(eType int, data map[string]interface{}) WSResource {
	switch resourceType(eType) {
	case CHANNEL_MESSAGE:
		return keySerialize{CHANNEL_MESSAGE, data, "message"}
	default:
		return noSerialize{resourceType(eType), data}
	}
}

type keySerialize struct {
	eType resourceType
	data  map[string]interface{}
	key   string
}

func (key keySerialize) Type() resourceType {
	return key.eType
}

func (key keySerialize) Decode(target interface{}) error {
	if key.data[key.key] == nil {
		return fmt.Errorf("key %s doesn't exist on data type %d", key.key, key.eType)
	}

	return mapstructure.Decode(key.data[key.key], target)
}

type noSerialize struct {
	eType resourceType
	data  map[string]interface{}
}

func (no noSerialize) Type() resourceType {
	return no.eType
}

func (no noSerialize) Decode(interface{}) error {
	return fmt.Errorf("no decoder for data type %d: %+v", no.eType, no.data)
}
