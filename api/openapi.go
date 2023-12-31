package api

import (
	"fmt"
	"reflect"
	"strings"

	"iothub/config"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
)

func OpenapiConfig() restfulspec.Config {
	c := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject,
		ModelTypeNameHandler: func(t reflect.Type) (string, bool) {
			key := t.String()
			key = strings.ReplaceAll(key, "/", ".")
			return key, true
		},
	}

	return c
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Security = []map[string][]string{{"basic": {}}}
	swo.SecurityDefinitions = map[string]*spec.SecurityScheme{
		"basic": spec.BasicAuth(),
	}
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title: "iothub",
			Description: "A tiny iothub core\n\n" +
				fmt.Sprintf("Build Info:\n- Version: %s\n- GitCommit: %s\n", config.Version, config.GitCommit),
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{
		{
			TagProps: spec.TagProps{
				Name:        "things",
				Description: ""},
		},
		{
			TagProps: spec.TagProps{
				Name:        "shadows",
				Description: ""},
		},
	}
}
