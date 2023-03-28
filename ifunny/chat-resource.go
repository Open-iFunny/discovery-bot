package ifunny

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

func (chat *Chat) call(desc call, output interface{}) error {
	traceID := uuid.New().String()
	log := chat.client.log.WithFields(logrus.Fields{
		"trace_id":  traceID,
		"type":      "CALL",
		"procedure": desc.procedure,
		"kwargs":    desc.kwargs,
	})

	log.Trace("exec call")
	result, err := chat.ws.Call(desc.procedure, desc.options, desc.args, desc.kwargs)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Trace(fmt.Sprintf("call OK recv: %+v\n", result.ArgumentsKw))
	return mapstructure.Decode(result.ArgumentsKw, output)
}
