package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"go.elastic.co/apm/module/apmgrpc/v2"
	"go.elastic.co/apm/module/apmhttp/v2"
	"google.golang.org/grpc"

	"fmt"
	"net"

	apigrpc "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/api/grpc"
	apihttp "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/api/http"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/constant"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/db"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/interceptors"
	pb "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/pb"
	svc "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/service"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/utils"

	servicelogger "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/service-logger"

	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/config"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/urfave/cli"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/spf13/viper"
)

var s *grpc.Server
var appConfig *config.Config
var logger *servicelogger.AddonsLogrus
var database *db.Db

func init() {

	appConfig = config.InitConfig()
	logger = servicelogger.New(&servicelogger.LoggerConfig{
		ServiceName:   appConfig.AppName,
		LogOutput:     appConfig.LoggerOutput,
		LogLevel:      appConfig.LoggerLevel,
		FluentbitHost: appConfig.FluentbitHost,
		FluentbitPort: appConfig.FluentbitPort,
	})

	location, locErr := utils.SetTimeLocation(appConfig.TimeLocation)
	if locErr != nil {
		panic(locErr)
	}

	time.Local = location
	database = db.NewDatabase(appConfig, logger)
}

func main() {
	app := cli.NewApp()
	app.Name = ""
	app.Commands = []cli.Command{
		grpcServerCmd(),
		gatewayServerCmd(),
		grpcGatewayServerCmd(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		os.Exit(1)
	}

}

func grpcServerCmd() cli.Command {

	return cli.Command{
		Name:  "grpc-server",
		Usage: "starts a gRPC server",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port",
				Value: appConfig.AppPort,
			},
		},
		Action: func(c *cli.Context) error {

			port := c.Int("port")

			// Wait for Control C to exit
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			database.StartDBConnection()

			go func() {
				if err := grpcServer(port); err != nil {
					logger.Fatalf("failed RPC serve: %v", err)
				}
			}()

			// Block until a signal is received
			<-ch

			database.CloseDBConnections()

			logger.Println("Stopping RPC server")
			s.Stop()
			logger.Println("RPC server stopped")
			return nil

		},
	}

}

func gatewayServerCmd() cli.Command {

	return cli.Command{
		Name:  "gw-server",
		Usage: "starts a Gateway server",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port",
				Value: 3000,
			},
			cli.StringFlag{
				Name:  "grpc-endpoint",
				Value: ":" + fmt.Sprint(appConfig.AppPort),
				Usage: "the address of the running gRPC server to transcode to",
			},
		},
		Action: func(c *cli.Context) error {

			port, grpcEndpoint := c.Int("port"), c.String("grpc-endpoint")

			go func() {
				if err := httpGatewayServer(port, grpcEndpoint); err != nil {
					logger.Fatalf("failed JSON Gateway serve: %v", err)
				}
			}()

			// Wait for Control C to exit
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			// Block until a signal is received
			<-ch

			logger.Println("JSON Gateway server stopped")

			return nil

		},
	}

}

func grpcGatewayServerCmd() cli.Command {

	return cli.Command{
		Name:  "grpc-gw-server",
		Usage: "Starts gRPC and Gateway server",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port1",
				Value: appConfig.AppPort,
			},
			cli.IntFlag{
				Name:  "port2",
				Value: 3000,
			},
			cli.StringFlag{
				Name:  "grpc-endpoint",
				Value: ":" + fmt.Sprint(appConfig.AppPort),
				Usage: "the address of the running gRPC server to transcode to",
			},
		},
		Action: func(c *cli.Context) error {

			rpcPort, httpPort, grpcEndpoint := c.Int("port1"), c.Int("port2"), c.String("grpc-endpoint")

			// Wait for Control C to exit
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			database.StartDBConnection()

			go func() {
				if err := grpcServer(rpcPort); err != nil {
					logger.Fatalf("failed RPC serve: %v", err)
				}
			}()

			go func() {
				if err := httpGatewayServer(httpPort, grpcEndpoint); err != nil {
					logger.Fatalf("failed JSON Gateway serve: %v", err)
				}
			}()

			// Block until a signal is received
			<-ch

			logger.Println("Stopping RPC server")
			s.GracefulStop()
			database.CloseDBConnections()
			logger.Println("RPC server stopped")
			logger.Println("JSON Gateway server stopped")

			return nil

		},
	}

}

func grpcServer(port int) error {

	logger.Printf("Starting %s Service ................", constant.ServiceName)
	logger.Printf("Starting RPC server on port %d...", port)

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	logger.Println("[service - connection] Initialize Other Service Connection...")
	svcConn := svc.InitServicesConn(
		logger,
		"",
		appConfig.AuthService,
	)
	defer svcConn.CloseAllServicesConn()

	dbProvider := db.NewProvider(database.DbSql)

	apiServer, apiErr := apigrpc.New(
		appConfig,
		logger,
		dbProvider,
		svcConn,
	)
	if apiErr != nil {
		return apiErr
	}


	kliringQueueName := utils.GetEnv("queue-status-kliring", "status-kliring-queue-local")
	err = apiServer.SetupRabbitMQConn("rabbit-conn-publisher-kliring", kliringQueueName)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	swiftQueueName := utils.GetEnv("queue-status-swift", "status-swift-queue-local")
	err = apiServer.SetupRabbitMQConn("rabbit-conn-publisher-swift", swiftQueueName)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	go apiServer.JobPending()

	interceptor := interceptors.NewInterceptor(logger)

	authInterceptor := interceptors.NewAuthInterceptor(appConfig.ApiServicePath, svcConn)
	unaryInterceptorOpt := grpc.UnaryInterceptor(interceptor.UnaryInterceptors(authInterceptor))
	streamInterceptorOpt := grpc.StreamInterceptor(interceptor.StreamInterceptors(authInterceptor))

	s = grpc.NewServer(unaryInterceptorOpt, streamInterceptorOpt)

	pb.RegisterTransactionStatusServiceServer(s, apiServer)

	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	return s.Serve(list)

}

func httpGatewayServer(port int, grpcEndpoint string) error {

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Connect to the GRPC server
	conn, err := grpc.Dial(
		grpcEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()),
		grpc.WithStreamInterceptor(apmgrpc.NewStreamClientInterceptor()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	httpHandler := &apihttp.HttpHandler{
		Logger: logger,
	}

	rmux := runtime.NewServeMux(
		runtime.WithErrorHandler(httpHandler.CustomHTTPError),
		runtime.WithForwardResponseOption(httpHandler.HttpResponseModifier),
	)

	client := pb.NewTransactionStatusServiceClient(conn)
	err = pb.RegisterTransactionStatusServiceHandlerClient(ctx, rmux, client)
	if err != nil {
		return err
	}

	// Serve the swagger-ui and swagger file
	mux := http.NewServeMux()
	mux.Handle("/", originMiddleware(rmux))

	mux.HandleFunc("/api/template/docs/swagger.json", serveSwagger)
	fs := http.FileServer(http.Dir(appConfig.SwaggerPath + "swagger-ui"))
	mux.Handle("/api/transaction-status/docs/", http.StripPrefix("/api/transaction-status/docs/", fs))

	// Start
	logger.Printf("Starting JSON Gateway server on port %d...", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), apmhttp.Wrap(setHeaders(mux)))

}

func serveSwagger(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, appConfig.SwaggerPath+"swagger.json")

}

func originMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")
		referer := r.Header.Get("Referer")

		method := r.Method
		path := r.URL.Path

		logger.Printf("[MIDDLEWARE] Request Start... Method: %s, Path: %s", method, path)

		if utils.GetEnv("ENV", "DEV") == "PROD" {
			envOrigin := r.Header.Get("ENV-Allow-Origin")
			if len([]rune(envOrigin)) > 0 {
				envOrigins := []string{}
				if envOrigin != "" {
					envOrigins = strings.Split(envOrigin, ",")
				}
				if len(envOrigins) > 0 && envOrigin != "" {
					envOriginsHeader := []string{}
					for _, v := range envOriginsHeader {
						envOriginsHeader = append(envOriginsHeader, strings.TrimSpace(v))
						//envOriginsHeader[i] = strings.TrimSpace(v)
					}
					if len(envOriginsHeader) > 0 {
						logger.Infof("Origin: %v - Ref: %v - ENV: %v", origin, referer, envOriginsHeader)
						pass := false
						if origin != "" {
							if len(envOriginsHeader) > 0 {
								for _, v := range envOriginsHeader {
									if origin == v {
										pass = true
									}
								}
							}
						}
						if referer != "" {
							if len(envOriginsHeader) > 0 {
								for _, v := range envOriginsHeader {
									if strings.Contains(referer, v) {
										pass = true
									}
								}
							}
						}
						if !pass {
							http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
							return
						}
					}
				}
			}
		}

		next.ServeHTTP(w, r)

	})

}

func allowedOrigin(origin string) bool {

	if utils.StringInSlice(viper.GetString("cors"), appConfig.CorsAllowedOrigins) {
		return true
	}

	if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
		return true
	}

	return false

}

func setHeaders(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Strict-Transport-Security", "max-age=31536000")

		if allowedOrigin(r.Header.Get("Origin")) {
			if utils.GetEnv("ENV", "DEV") != "PROD" {
				w.Header().Set("Content-Security-Policy", "object-src 'none'; child-src 'none'; script-src 'unsafe-inline' https: http: ")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.Header().Set("X-Frame-Options", "DENY")
				w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
				w.Header().Set("X-XSS-Protection", "1; mode=block")
				w.Header().Set("Permissions-Policy", "geolocation=()")
				w.Header().Set("Referrer-Policy", "no-referrer")

				w.Header().Set("Access-Control-Allow-Origin", strings.Join(appConfig.CorsAllowedOrigins, ", "))
			}
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(appConfig.CorsAllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(appConfig.CorsAllowedHeaders, ", "))
			w.Header().Add("Access-Control-Expose-Headers", strings.Join(appConfig.ExposedHeaders, ", "))
		}

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)

	})

}
