package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/lipgloss"
)

var (
	labelStyle = lipgloss.NewStyle() // Default (white/black)

	ipStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "214", Dark: "220"}) // Gold/Yellow

	interfaceStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "241", Dark: "245"}) // Gray

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "196", Dark: "204"}) // Red/Pink
)

var Version = "2.0.0"

type CLI struct {
	Local      bool `short:"l" help:"Show local non-loopback IPv4 addresses"`
	Gateway    bool `short:"g" help:"Show gateway IP"`
	External   bool `short:"e" help:"Show external IP address"`
	All        bool `short:"a" help:"When combined with other flags, shows all interfaces"`
	NoHeaders  bool `short:"n" help:"Don't show headers (for scripting)"`
	ShowBridge bool `short:"b" help:"Include bridge interfaces in local IPs"`
	Version    kong.VersionFlag `short:"v" help:"Show version"`
}

func getLocalIPs(includeBridge bool) []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var ips []string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		if !includeBridge && strings.HasPrefix(iface.Name, "bridge") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				ips = append(ips, fmt.Sprintf("%s - %s", ipnet.IP.String(), iface.Name))
			}
		}
	}

	return ips
}

func getGatewayIP() string {
	cmd := exec.Command("netstat", "-rn")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "default") && !strings.Contains(line, "::") {
			fields := strings.Fields(line)
			// netstat -rn format: destination gateway flags refs use interface
			// We expect at least: "default" "192.168.1.1" ...
			if len(fields) < 2 {
				continue
			}

			gateway := fields[1]
			// Assert it looks like an IP
			if parts := strings.Split(gateway, "."); len(parts) != 4 {
				continue
			}

			return gateway
		}
	}

	return ""
}

func getExternalIP() string {
	cmd := exec.Command("dig", "+short", "myip.opendns.com", "@resolver1.opendns.com")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	ip := strings.TrimSpace(string(output))

	// Assert it looks like an IP (dig should return a single IP address)
	if parts := strings.Split(ip, "."); len(parts) != 4 {
		return ""
	}

	return ip
}

func main() {
	var cli CLI
	kong.Parse(&cli,
		kong.Name("ip"),
		kong.Description("Simple tool to get your IP addresses"),
		kong.UsageOnError(),
		kong.Vars{"version": Version},
	)

	// No flags = show all
	showAll := !cli.Local && !cli.Gateway && !cli.External

	// Pretty output only for default view, plain output for specific flags (piping)
	prettyOutput := showAll && !cli.NoHeaders

	// Bridge interfaces: show in default view, with -b, or with -la
	includeBridge := showAll || cli.ShowBridge || (cli.Local && cli.All)

	// Print what was requested
	needSpacer := false

	if cli.Local || showAll {
		if ips := getLocalIPs(includeBridge); len(ips) > 0 {
			if prettyOutput {
				fmt.Println(labelStyle.Render("Local IPs:"))
				for _, ip := range ips {
					parts := strings.Split(ip, " - ")
					fmt.Printf("  %s %s\n",
						ipStyle.Render(parts[0]),
						interfaceStyle.Render("("+parts[1]+")"))
				}
			} else {
				for _, ip := range ips {
					parts := strings.Split(ip, " - ")
					fmt.Println(parts[0])
				}
			}
			needSpacer = true
		}
	}

	if cli.Gateway || showAll {
		if ip := getGatewayIP(); ip != "" {
			if prettyOutput && needSpacer {
				fmt.Println()
			}
			if prettyOutput {
				fmt.Println(labelStyle.Render("Gateway IP:"))
				fmt.Printf("  %s\n", ipStyle.Render(ip))
			} else {
				fmt.Println(ip)
			}
			needSpacer = true
		}
	}

	if cli.External || showAll {
		if ip := getExternalIP(); ip != "" {
			if prettyOutput && needSpacer {
				fmt.Println()
			}
			if prettyOutput {
				fmt.Println(labelStyle.Render("External IP:"))
				fmt.Printf("  %s\n", ipStyle.Render(ip))
			} else {
				fmt.Println(ip)
			}
		}
	}
}
