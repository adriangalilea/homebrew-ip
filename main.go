package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/lipgloss"
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

var Version = "2.1.1"

type CLI struct {
	Local      bool             `short:"l" help:"Show local non-loopback IPv4 addresses"`
	Gateway    bool             `short:"g" help:"Show gateway IP"`
	External   bool             `short:"e" help:"Show external IP address"`
	All        bool             `short:"a" help:"When combined with other flags, shows all interfaces"`
	NoHeaders  bool             `short:"n" help:"Don't show headers (for scripting)"`
	ShowBridge bool             `short:"b" help:"Include bridge interfaces in local IPs"`
	Version    kong.VersionFlag `short:"v" help:"Show version"`
}

type IPEntry struct {
	Addr      string
	Interface string
}

type Section struct {
	Label   string
	Entries []IPEntry
	Err     error
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
	output, err := exec.Command("netstat", "-rn").Output()
	if err != nil {
		return "", fmt.Errorf("netstat failed: %w", err)
	}

	for _, line := range strings.Split(string(output), "\n") {
		if !strings.Contains(line, "default") || strings.Contains(line, "::") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if net.ParseIP(fields[1]) != nil {
			return fields[1], nil
		}
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

func renderSections(sections []Section, pretty bool) {
	first := true
	for _, s := range sections {
		if s.Err != nil {
			if pretty && !first {
				fmt.Println()
			}
			fmt.Fprintf(os.Stderr, "%s %s\n",
				errorStyle.Render(s.Label+":"),
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
			fmt.Println(labelStyle.Render(s.Label + ":"))
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
		kong.Vars{"version": Version},
	)

	showAll := !cli.Local && !cli.Gateway && !cli.External
	pretty := showAll && !cli.NoHeaders
	includeBridge := showAll || cli.ShowBridge || (cli.Local && cli.All)

	var sections []Section

	if cli.Local || showAll {
		entries, err := getLocalIPs(includeBridge)
		sections = append(sections, Section{Label: "Local IPs", Entries: entries, Err: err})
	}

	if cli.Gateway || showAll {
		ip, err := getGatewayIP()
		var entries []IPEntry
		if err == nil {
			entries = []IPEntry{{Addr: ip}}
		}
		sections = append(sections, Section{Label: "Gateway IP", Entries: entries, Err: err})
	}

	if cli.External || showAll {
		ip, err := getExternalIP()
		var entries []IPEntry
		if err == nil {
			entries = []IPEntry{{Addr: ip}}
		}
		sections = append(sections, Section{Label: "External IP", Entries: entries, Err: err})
	}

	renderSections(sections, pretty)
}
