package password

import (
	"context"

	"iothub/connector/mqtt/embed"

	"iothub/config"
	"iothub/pkg/log"
	"iothub/thing"
)

func AuthzMqttClient(ctx context.Context, superUsers []config.UserPassword, thingSvc thing.Service) embed.AuthzFn {
	return func(user, password string) bool {
		for _, u := range superUsers {
			if user == u.Name && password == u.Password {
				log.Infof("Mqtt client user %s is authorized by default users", u.Name)
				return true
			}
		}
		th, err := thingSvc.Get(ctx, user)
		if err != nil {
			log.Infof("Mqtt client user %s authz error: %v", user, err)
			return false
		}
		if th.AuthValue == password {
			return true
		} else {
			log.Infof("Mqtt client user %s password %s is not authorized", user, password)
			return false
		}
	}
}
