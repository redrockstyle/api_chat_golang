package httpserver

import (
	"api_chat/api/layers/base/logx"
	httprouter "api_chat/api/layers/controller/server/http_router"
	config "api_chat/api/layers/domain/cfg"
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type HttpApiServer struct {
	//handler *httphandler.HttpApiHandler
	apr  *httprouter.ApiRouter
	cfg  *config.Configuration
	logx logx.Logger
}

func NewRestApiServer(apr *httprouter.ApiRouter, cfg *config.Configuration, logx logx.Logger) *HttpApiServer {
	return &HttpApiServer{apr: apr, cfg: cfg, logx: logx}
}

func (has *HttpApiServer) Runtime() {
	has.logx.Infof("Server starting on %v", has.cfg.Server.Port)

	has.apr.InitRouter(has.cfg.Server.PrefixPath)
	if has.cfg.SSL.Active {
		has.logx.Info("SSL enable")
		if has.cfg.SSL.Domain == "" || has.cfg.SSL.Domain == "localhost" {
			go func() {
				has.logx.Infof("Mode development is active: use %v and %v", has.cfg.SSL.CrtPath, has.cfg.SSL.KeyPath)
				if err := fasthttp.ListenAndServeTLS(has.cfg.Server.Port, has.cfg.SSL.CrtPath, has.cfg.SSL.KeyPath, has.apr.GetHandler()); err != nil {
					has.logx.Warnf("Server returned err: %v", err)
					return
				}
			}()
		} else {
			has.logx.Infof("Set hostname %v", has.cfg.Server.Domain)
			var cfg *tls.Config
			var lnTls net.Listener
			cfg = GenerateTlsConfig(has.cfg.Server.Domain)
			has.logx.Info("Config port ignored - set port 443")
			ln, err := net.Listen("tcp4", "0.0.0.0:443")
			if err != nil {
				has.logx.Errorf("Error listen 443 port: %v", err)
				return
			}
			lnTls = tls.NewListener(ln, cfg)
			if err := fasthttp.Serve(lnTls, has.apr.GetHandler()); err != nil {
				has.logx.Warnf("Server returned err: %v", err)
				return
			}
		}

	} else {
		go func() {
			has.logx.Info("SSL disable")
			//if err := fasthttp.ListenAndServe(has.cfg.Server.Port, has.handler.RequestHadler); err != nil {
			if err := fasthttp.ListenAndServe(has.cfg.Server.Port, has.apr.GetHandler()); err != nil {
				has.logx.Warnf("Server returned err: %v", err)
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	has.logx.Info("Service stopped")
}

func GenerateTlsConfig(domain string) *tls.Config {
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache("./certs"),
	}

	return &tls.Config{
		GetCertificate: m.GetCertificate,
		NextProtos: []string{
			"http/1.1", acme.ALPNProto,
		},
	}
}
