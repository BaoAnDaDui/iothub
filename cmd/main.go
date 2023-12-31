package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"iothub"
	"iothub/api"
	"iothub/connector/mqtt/client"
	"iothub/connector/mqtt/embed"

	"iothub/auth/password"

	"iothub/config"
	mq "iothub/connector/mqtt"
	"iothub/db/mysql"
	"iothub/db/sqlite"
	"iothub/pkg/log"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/kpango/glg"
	"gorm.io/gorm"

	"iothub/shadow"
	shadowWire "iothub/shadow/wire"

	shadowApi "iothub/shadow/api"
	"iothub/thing"
	thingApi "iothub/thing/api"
	thingWire "iothub/thing/wire"
)

var (
	Version   = ""
	GitCommit = ""
)

const (
	stopWaitTime = time.Second * 1
)

func main() {
	config.Version = Version
	config.GitCommit = GitCommit
	log.Infof("Version: %s GitCommit: %s", Version, GitCommit)

	// load config
	cfg := config.ReadConfig()
	cfgJ, _ := json.Marshal(cfg)
	log.Infof("Config: %s", cfgJ)

	// set log level
	if cfg.Log.Level != "" {
		ll := glg.Atol(cfg.Log.Level)
		if ll == glg.DEBG || ll == glg.INFO ||
			ll == glg.WARN || ll == glg.ERR || ll == glg.FATAL {
			glg.Get().SetLevel(ll)
		} else {
			log.Fatalf("Wrong log level %q", cfg.Log.Level)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if sig := signalHandler(ctx); sig != nil {
			cancel()
			log.Info(fmt.Sprintf("iothub shutdown by signal: %s", sig))
		}
	}()

	dbConn := newDb(cfg)
	autoMigrate(dbConn)

	// mqtt client and connector for iothub interacts with message broker

	mqttClient := client.NewClient(cfg.Connector.MqttClient)
	connector := mq.InitConnector(cfg.Connector, mqttClient)

	// services
	shadowSvc := shadowWire.InitSvc(dbConn, connector)
	thingSvc := thingWire.InitSvc(ctx, dbConn, shadowSvc, connector)

	// embedded mqtt broker
	if cfg.Connector.Typ == config.ConnectorMqttEmbed {
		authzFn := password.AuthzMqttClient(ctx, cfg.Connector.MqttBroker.SuperUsers, thingSvc)
		startMqttBroker(ctx, cfg.Connector.MqttBroker, authzFn)
	}

	// init
	if err := connector.Start(ctx); err != nil {
		log.Fatalf("Mqtt connector start error: %v", err)
	}
	if err := connector.InitMethodHandler(ctx); err != nil {
		log.Fatalf("Connector init method handler error: %v", err)
	}
	if err := connector.InitNtpHandler(ctx); err != nil {
		log.Fatalf("Connector init ntp handler error: %v", err)
	}
	if err := shadow.Link(ctx, connector, shadowSvc); err != nil {
		log.Fatalf("Link shadow service to connector error %v", err)
	}
	if err := mqttClient.Connect(ctx); err != nil {
		log.Fatalf("Mqtt client start error: %v", err)
	}

	// htt api

	iothub.RouteSwagger()
	iothub.RouteWeb()
	azf := api.BasicAuthMiddleware(cfg.API.BasicAuth.Name, cfg.API.BasicAuth.Password)
	thingWs := thingApi.Service(ctx, thingSvc).
		Filter(api.LoggingMiddleware).
		Filter(azf)
	shadowApi.Service(ctx, thingWs, shadowSvc, thingSvc, connector)

	restful.DefaultContainer.Add(thingWs)
	restful.DefaultContainer.Add(mq.Service(ctx, connector).Filter(azf))
	restful.DefaultContainer.Add(thingApi.ServiceForEmqxIntegration())
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(api.OpenapiConfig()))
	startHttpSvr(ctx, cfg.API.Port, nil)

	// wait some seconds before shutting down
	time.Sleep(1 * time.Second)
}

func startHttpSvr(ctx context.Context, apiPort int, handler http.Handler) {
	addr := fmt.Sprintf(":%d", apiPort)
	server := &http.Server{Addr: addr, Handler: handler}
	errCh := make(chan error)
	go func() {
		log.Infof("Http listening on %s", addr)
		errCh <- server.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		ctxShutdown, cancelShutdown := context.WithTimeout(ctx, stopWaitTime)
		defer cancelShutdown()
		if err := server.Shutdown(ctxShutdown); err != nil {
			log.Errorf("Http server error occurred during shutdown at %s: %s", addr, err)
		}
		log.Info(fmt.Sprintf("Http server shutdown of http at %s", addr))
	case err := <-errCh:
		log.Errorf("Http server exit cause: %v", err)
	}
}

func autoMigrate(conn *gorm.DB) {
	err := conn.AutoMigrate(&thing.Entity{}, &shadow.Entity{})
	if err != nil {
		log.Fatalf("auto migrate db error: %v", err)
	}
	time.Sleep(time.Millisecond * 100)
}

func startMqttBroker(ctx context.Context, cfg config.InnerMqttBroker, authzFn embed.AuthzFn) embed.Broker {
	return embed.InitBroker(embed.MochiConfig{
		TcpPort:    cfg.TcpPort,
		TcpSslPort: cfg.TcpSslPort,
		WsPort:     cfg.WsPort,
		WssPort:    cfg.WssPort,
		KeyFile:    cfg.KeyFile,
		CertFile:   cfg.CertFile,
		AuthzFn:    authzFn,
		AclFn: func(user string, topic string, write bool) bool {
			return thing.TopicAcl(cfg.SuperUsers, user, topic, write)
		},
		SuperUsers: cfg.SuperUsers,
	})
}

func newDb(cfg config.Config) *gorm.DB {
	switch cfg.DB.Typ {
	case config.DBMySQL:
		return newMysqlDB(cfg.DB.Mysql)
	case config.DBSqlite:
		return newSqliteDB(cfg.DB.Sqlite)
	default:
		log.Fatal("Unknown database type: ", cfg.DB.Typ)
	}
	return nil
}

func newSqliteDB(cfg sqlite.Config) *gorm.DB {
	db, err := sqlite.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func newMysqlDB(cfg mysql.Config) *gorm.DB {
	conn, err := mysql.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func signalHandler(ctx context.Context) error {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	select {
	case sig := <-c:
		return fmt.Errorf("%s", sig)
	case <-ctx.Done():
		return nil
	}
}
