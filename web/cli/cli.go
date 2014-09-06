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

	gcli "github.com/codegangsta/cli"
	"github.com/skratchdot/open-golang/open"
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

type ServerInfo struct {
	Mux     *http.ServeMux
	Addr    string
	MayOpen bool
}

// return subcommands for our particular area of algorithms
func SubCommands() []gcli.Command {
	out := make([]gcli.Command, 1)

	var info ServerInfo

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

func before(c *gcli.Context, info *ServerInfo) error {
	info.Addr = c.String(addressFlagName)
	if info.Addr == "" {
		return errors.New("Server address must not be empty")
	}

	info.MayOpen = !c.Bool(noOpenFlagName)
	return nil
}

func action(c *gcli.Context, info *ServerInfo) {
	RunWebServer(info)
}

func RunWebServer(info *ServerInfo) {
	addr := info.Addr
	info.Mux = server.NewHandler()
	if !strings.HasPrefix(addr, httpProtocol) {
		addr = httpProtocol + info.Addr
	}

	fmt.Println("About to listen on ", addr)
	fmt.Println("Hit CTRL+C to close")
	if info.MayOpen {
		if err := open.Start(addr); err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}

	s := http.Server{
		Addr:    info.Addr,
		Handler: info.Mux,
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
