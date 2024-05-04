# ip
 Simplest cli tool to get your IP on Mac (local, external, gateway)

# install
`brew install ip`

# usage
local, default option, also `-l`
```
❯ ip
Local IPs:
  192.168.1.38 - ine
  192.168.1.50 - ine
```

gateway
```
❯ ip -g
Gateway IP:
  192.168.1.1
```
external
```
❯ ip -e
External IP:
  88.7.14.121
```
all
```
❯ ip -a
Local IPs:
  192.168.1.38 - ine
  192.168.1.50 - ine
Gateway IP:
  192.168.1.1
External IP:
  88.7.14.121
```
