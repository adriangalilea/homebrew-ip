package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/net/route"
)

var (
	labelStyle = lipgloss.NewStyle()

	ipStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "214", Dark: "220"})

	interfaceStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "241", Dark: "245"})

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "196", Dark: "204"})
)

// The git tag is the version. The Homebrew formula builds from the tag's
// source tarball, which carries no VCS data, so it passes the version with
// -ldflags "-X main.version=..."; `go install ...@vX` records it in the build
// info instead. A build that has neither is a development build.
var version = ""

func resolvedVersion() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return strings.TrimPrefix(info.Main.Version, "v")
	}
	return "dev"
}

type CLI struct {
	Local      bool             `short:"l" help:"Show local non-loopback IPv4 addresses"`
	Gateway    bool             `short:"g" help:"Show gateway IP"`
	External   bool             `short:"e" help:"Show external IP address"`
	All        bool             `short:"a" help:"When combined with other flags, shows all interfaces"`
	NoHeaders  bool             `short:"n" help:"Don't show headers (for scripting)"`
	ShowBridge bool             `short:"b" help:"Include bridge interfaces in local IPs"`
	JSON       bool             `short:"j" help:"Output as JSON"`
	Version    kong.VersionFlag `short:"v" help:"Show version"`
}

type IPEntry struct {
	Addr      string `json:"addr"`
	Interface string `json:"interface,omitempty"`
}

type Kind int

const (
	Local Kind = iota
	Gateway
	External
)

func (k Kind) Label() string {
	switch k {
	case Local:
		return "Local IPs"
	case Gateway:
		return "Gateway IP"
	case External:
		return "External IP"
	}
	panic(fmt.Sprintf("unknown kind %d", k))
}

type Section struct {
	Kind    Kind
	Entries []IPEntry
	Err     error
}

type JSONOutput struct {
	Local         []IPEntry `json:"local,omitempty"`
	LocalError    string    `json:"local_error,omitempty"`
	Gateway       string    `json:"gateway,omitempty"`
	GatewayError  string    `json:"gateway_error,omitempty"`
	External      string    `json:"external,omitempty"`
	ExternalError string    `json:"external_error,omitempty"`
}

func getLocalIPs(includeBridge bool) ([]IPEntry, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("listing interfaces: %w", err)
	}

	var entries []IPEntry
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if !includeBridge && strings.HasPrefix(iface.Name, "bridge") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("reading addrs for %s: %w", iface.Name, err)
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.To4() == nil {
				continue
			}
			entries = append(entries, IPEntry{Addr: ipnet.IP.String(), Interface: iface.Name})
		}
	}

	return entries, nil
}

func getGatewayIP() (string, error) {
	rib, err := route.FetchRIB(syscall.AF_INET, syscall.NET_RT_DUMP, 0)
	if err != nil {
		return "", fmt.Errorf("fetching routing table: %w", err)
	}

	msgs, err := route.ParseRIB(syscall.NET_RT_DUMP, rib)
	if err != nil {
		return "", fmt.Errorf("parsing routing table: %w", err)
	}

	for _, msg := range msgs {
		rm, ok := msg.(*route.RouteMessage)
		if !ok || len(rm.Addrs) <= syscall.RTAX_GATEWAY {
			continue
		}

		dst, ok := rm.Addrs[syscall.RTAX_DST].(*route.Inet4Addr)
		if !ok || dst.IP != [4]byte{0, 0, 0, 0} {
			continue
		}

		gw, ok := rm.Addrs[syscall.RTAX_GATEWAY].(*route.Inet4Addr)
		if !ok {
			continue
		}

		return net.IPv4(gw.IP[0], gw.IP[1], gw.IP[2], gw.IP[3]).String(), nil
	}

	return "", fmt.Errorf("no default gateway found")
}

func getExternalIP() (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	ip := strings.TrimSpace(string(body))
	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("invalid IP returned: %q", ip)
	}

	return ip, nil
}

func renderJSON(sections []Section) error {
	var out JSONOutput
	for _, s := range sections {
		switch s.Kind {
		case Local:
			if s.Err != nil {
				out.LocalError = s.Err.Error()
			} else {
				out.Local = s.Entries
			}
		case Gateway:
			if s.Err != nil {
				out.GatewayError = s.Err.Error()
			} else {
				out.Gateway = s.Entries[0].Addr
			}
		case External:
			if s.Err != nil {
				out.ExternalError = s.Err.Error()
			} else {
				out.External = s.Entries[0].Addr
			}
		}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func renderSections(sections []Section, pretty bool) {
	first := true
	for _, s := range sections {
		if s.Err != nil {
			if pretty && !first {
				fmt.Println()
			}
			fmt.Fprintf(os.Stderr, "%s %s\n",
				errorStyle.Render(s.Kind.Label()+":"),
				errorStyle.Render(s.Err.Error()))
			first = false
			continue
		}
		if len(s.Entries) == 0 {
			continue
		}
		if pretty && !first {
			fmt.Println()
		}
		if pretty {
			fmt.Println(labelStyle.Render(s.Kind.Label() + ":"))
			for _, e := range s.Entries {
				if e.Interface != "" {
					fmt.Printf("  %s %s\n", ipStyle.Render(e.Addr), interfaceStyle.Render("("+e.Interface+")"))
				} else {
					fmt.Printf("  %s\n", ipStyle.Render(e.Addr))
				}
			}
		} else {
			for _, e := range s.Entries {
				fmt.Println(e.Addr)
			}
		}
		first = false
	}
}

func main() {
	var cli CLI
	kong.Parse(&cli,
		kong.Name("ip"),
		kong.Description("Simple tool to get your IP addresses"),
		kong.UsageOnError(),
		kong.Vars{"version": resolvedVersion()},
	)

	showAll := !cli.Local && !cli.Gateway && !cli.External
	pretty := showAll && !cli.NoHeaders
	includeBridge := showAll || cli.ShowBridge || (cli.Local && cli.All)

	var sections []Section

	if cli.Local || showAll {
		entries, err := getLocalIPs(includeBridge)
		sections = append(sections, Section{Kind: Local, Entries: entries, Err: err})
	}

	if cli.Gateway || showAll {
		sections = append(sections, single(Gateway, getGatewayIP))
	}

	if cli.External || showAll {
		sections = append(sections, single(External, getExternalIP))
	}

	if cli.JSON {
		if err := renderJSON(sections); err != nil {
			fmt.Fprintln(os.Stderr, errorStyle.Render("writing JSON: "+err.Error()))
			os.Exit(1)
		}
	} else {
		renderSections(sections, pretty)
	}

	// A lookup that failed is a failed run, so `[ "$(ip -e)" != ... ]` in a
	// script can tell "no answer" from "a different answer".
	for _, s := range sections {
		if s.Err != nil {
			os.Exit(1)
		}
	}
}

func single(kind Kind, lookup func() (string, error)) Section {
	ip, err := lookup()
	if err != nil {
		return Section{Kind: kind, Err: err}
	}
	return Section{Kind: kind, Entries: []IPEntry{{Addr: ip}}}
}
