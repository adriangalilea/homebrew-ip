# ip

Simplest CLI tool to get your IP addresses on macOS (local, gateway, external).

## Installation

```bash
brew install adriangalilea/tap/ip
```

or `go install github.com/adriangalilea/ip@latest`. If you installed it from the old `adriangalilea/ip` tap, `brew update` moves it to `adriangalilea/tap` for you.

## Usage

### Default (show all with pretty formatting)
```bash
$ ip
Local IPs:
  192.168.1.95 (en0)
  192.168.64.1 (bridge100)

Gateway IP:
  192.168.1.1

External IP:
  88.6.43.97
```

### Specific IPs (plain output for scripting)
```bash
$ ip -l  # Local IP
192.168.1.95

$ ip -g  # Gateway IP
192.168.1.1

$ ip -e  # External IP
88.6.43.97
```

### Combinations
```bash
$ ip -lg   # Local + Gateway
192.168.1.95
192.168.1.1

$ ip -lge  # All IPs (plain output)
192.168.1.95
192.168.1.1
88.6.43.97

$ ip -la   # Local with all interfaces (including bridges)
192.168.1.95
192.168.64.1
```

### JSON output
```bash
$ ip -j
{
  "local": [
    { "addr": "192.168.1.95", "interface": "en0" },
    { "addr": "192.168.64.1", "interface": "bridge100" }
  ],
  "gateway": "192.168.1.1",
  "external": "88.6.43.97"
}

$ ip -jl  # Just local as JSON
{
  "local": [
    { "addr": "192.168.1.95", "interface": "en0" }
  ]
}
```

### Options
- `-l` - Show local IP
- `-g` - Show gateway IP
- `-e` - Show external IP
- `-j` - Output as JSON
- `-a` - When combined with other flags, shows all interfaces
- `-b` - Include bridge interfaces
- `-n` - No headers (force plain output)
- `-h` - Help

### Scripting Examples
```bash
# Copy external IP to clipboard
ip -e | pbcopy

# SSH to a machine using your local IP
ssh user@$(ip -l)

# Check if VPN is connected (external IP changes)
[ "$(ip -e)" != "YOUR_NORMAL_IP" ] && echo "VPN is active"

# Get local IP for a specific interface
ip -j | jq '.local[] | select(.interface == "en0") | .addr'
```

## Release

The tag is the version: `git tag v2.3.0 && git push origin v2.3.0` (the tag on its own, or no tag event fires). CI verifies, publishes the GitHub release with notes, and points `Formula/ip.rb` in [adriangalilea/homebrew-tap](https://github.com/adriangalilea/homebrew-tap) at the tag, through the shared workflows in [adriangalilea/ci](https://github.com/adriangalilea/ci).

