package mqtt

import (
	"sync"

	"iothub/connector/mqtt/client"
	"iothub/connector/mqtt/embed"
	"iothub/ntp"

	"iothub/config"
	"iothub/connector/mqtt/emqx"
	"iothub/pkg/log"
	"iothub/shadow"
)

const (
	DefaultQos = byte(1)
)

type mqttConnector struct {
	client client.Client
	shadow.Connectivity
	shadow.StateHandler
	shadow.MethodHandler
	ntp.Handler
}

var onceNewConnector sync.Once
var connectorSingleton shadow.Connector

func InitConnector(cfg config.Connector, cl client.Client) shadow.Connector {
	var c shadow.Connectivity
	typ := cfg.Typ
	if typ == config.ConnectorMqttEmbed {
		c = embed.NewEmbedAdapter()
		log.Infof("Use embed connector")
	} else if typ == config.ConnectorEmqx {
		c = emqx.NewEmqxAdapter(cfg.Emqx, cl)
		log.Infof("Use emqx connector")
	} else {
		log.Fatalf("Unsupported connector type %s", typ)
	}

	onceNewConnector.Do(func() {
		s := NewShadowHandler(cl)
		m := NewMethodHandler(cl, c)
		n := NewNtpHandler(cl)
		connectorSingleton = &mqttConnector{cl, c, s, m, n}
	})
	return connectorSingleton
}

func Connector() shadow.Connector { return connectorSingleton }

var _ shadow.Connector = (*mqttConnector)(nil)
