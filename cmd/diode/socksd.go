// Diode Network Client
// Copyright 2019 IoT Blockchain Technology Corporation LLC (IBTC)
// Licensed under the Diode License, Version 1.0
package main

import (
	"fmt"

	"github.com/diodechain/diode_go_client/command"
	"github.com/diodechain/diode_go_client/config"
	"github.com/diodechain/diode_go_client/rpc"
)

var (
	socksdCmd = &command.Command{
		Name:        "socksd",
		HelpText:    `  Enable a socks proxy for use with browsers and other apps.`,
		ExampleText: `  diode socksd -socksd_port 8082 -socksd_host 127.0.0.1`,
		Run:         socksdHandler,
	}
)

func init() {
	cfg := config.AppConfig
	socksdCmd.Flag.StringVar(&cfg.SocksServerHost, "socksd_host", "127.0.0.1", "host of socks server listening to")
	socksdCmd.Flag.IntVar(&cfg.SocksServerPort, "socksd_port", 1080, "port of socks server listening to")
	socksdCmd.Flag.StringVar(&cfg.SocksFallback, "fallback", "localhost", "how to resolve web2 addresses")
}

func socksdHandler() (err error) {
	err = app.Start()
	if err != nil {
		return
	}
	cfg := config.AppConfig
	client := app.datapool.GetClientByOrder(1)
	cfg.EnableSocksServer = true
	cfg.EnableProxyServer = true
	cfg.ProxyServerPort = 8080
	if cfg.EnableAPIServer {
		configAPIServer := NewConfigAPIServer(cfg)
		configAPIServer.SetAddr(cfg.APIServerAddr)
		configAPIServer.ListenAndServe()
		app.SetConfigAPIServer(configAPIServer)
	}
	socksServer := client.NewSocksServer(app.datapool)
	if len(cfg.Binds) > 0 {
		socksServer.SetBinds(cfg.Binds)
		printInfo("")
		printLabel("Bind      <name>", "<mode>     <remote>")
		for _, bind := range cfg.Binds {
			printLabel(fmt.Sprintf("Port      %5d", bind.LocalPort), fmt.Sprintf("%5s     %11s:%d", config.ProtocolName(bind.Protocol), bind.To, bind.ToPort))
		}
	}
	socksServer.SetConfig(&rpc.Config{
		Addr:            cfg.SocksServerAddr(),
		FleetAddr:       cfg.FleetAddr,
		Blocklists:      cfg.Blocklists,
		Allowlists:      cfg.Allowlists,
		EnableProxy:     false,
		ProxyServerAddr: cfg.ProxyServerAddr(),
		Fallback:        cfg.SocksFallback,
	})
	if err = socksServer.Start(); err != nil {
		cfg.Logger.Error(err.Error())
		return
	}
	app.SetSocksServer(socksServer)
	app.Wait()
	return
}
