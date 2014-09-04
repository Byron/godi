// Package cli implements the 'web' subcommand to host the web-frontend.
package cli

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Byron/godi/web/server"

	"github.com/Byron/open-golang/open"
	gcli "github.com/codegangsta/cli"
)

const (
	Name         = "web"
	httpProtocol = "http://"
	usage        = `
	Launch a web-based front-end to control all godi operations in your browser.

`
	addressFlagName = "address"
	noOpenFlagName  = "no-show"
)

var (
	defaultAddress = "localhost:9078"
	hostUsage      = fmt.Sprintf(`The address to serve on, defaults to '%s'.
Use 'localhost' to serve for users on this host only, '0.0.0.0' will serve
to everyone who can reach this host on any network interface.
You may also specify any IP assigned to an interface to restrict availability.`, defaultAddress)
)

type serverInfo struct {
	mux     *http.ServeMux
	addr    string
	mayOpen bool
}

// return subcommands for our particular area of algorithms
func SubCommands() []gcli.Command {
	out := make([]gcli.Command, 1)

	var info serverInfo

	web := gcli.Command{
		Name:      Name,
		ShortName: "",
		Usage:     usage,
		Before:    func(c *gcli.Context) error { return before(c, &info) },
		Action:    func(c *gcli.Context) { action(c, &info) },
		Flags: []gcli.Flag{
			gcli.BoolFlag{
				Name:  noOpenFlagName,
				Usage: "If set, the hosted web-site will not be opened automatically",
			},
			gcli.StringFlag{
				Name:  fmt.Sprintf("%s, a", addressFlagName),
				Value: defaultAddress,
				Usage: hostUsage,
			},
		},
	}

	out[0] = web
	return out
}

func before(c *gcli.Context, info *serverInfo) error {
	info.addr = c.String(addressFlagName)
	if info.addr == "" {
		return errors.New("Server address must not be empty")
	}

	info.mayOpen = !c.Bool(noOpenFlagName)
	info.mux = server.NewHandler()
	return nil
}

func action(c *gcli.Context, info *serverInfo) {
	addr := info.addr
	if !strings.HasPrefix(addr, httpProtocol) {
		addr = httpProtocol + info.addr
	}

	fmt.Println("About to listen on ", addr)
	if info.mayOpen {
		if err := open.Start(addr); err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}

	s := http.Server{
		Addr:    info.addr,
		Handler: info.mux,
	}

	// Respond to abort requests
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		os.Exit(2)
	}()

	err := s.ListenAndServe()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(3)
	}
}
