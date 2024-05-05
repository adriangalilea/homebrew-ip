#!/bin/bash

# TODO change style so that current one is under `-v` "verbose" and running ip simply outputs "192.168.1.X" so it can be used in scripting easily
# TODO add `test do` in order to try and merge this in brew
# TODO set up and verify actions to release homebrew automatically
# TODO add colored output
# labels: style

# Function to display help
display_help() {
    echo "Usage: $0 [option...]"
    echo "Displays IP information"
    echo ""
    echo "   -l     Show local non-loopback IPv4 addresses"
    echo "   -g     Show gateway IP"
    echo "   -e     Show external IP address"
    echo "   -a     Show all of the above information"
    echo "   -h     Display this help and exit"
    echo ""
}

# Initialize variables
show_local=0
show_gateway=0
show_external=0
show_all=0

# Function to display local non-loopback IPv4 addresses along with their interface names
display_local_ips() {
    echo "Local IPs:"
    ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print "  " $2 " - " substr($1, 1, length($1)-1)}'
}

# Function to display the default gateway IP (IPv4 only)
display_gateway_ip() {
    gateway_ip=$(netstat -rn | grep 'default' | grep -v '::' | awk '{print $2}' | head -n 1)
    echo "Gateway IP:"
    echo "  $gateway_ip"
}

# Function to display the external IP address
display_external_ip() {
    external_ip=$(dig +short myip.opendns.com @resolver1.opendns.com)
    echo "External IP:"
    echo "  $external_ip"
}

# Parse command-line options
while getopts "lgeah" opt; do
  case $opt in
    l)
      show_local=1
      ;;
    g)
      show_gateway=1
      ;;
    e)
      show_external=1
      ;;
    a)
      show_all=1
      ;;
    h)
      display_help
      exit 0
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      display_help
      exit 1
      ;;
  esac
done

# Determine what to display based on flags
if [ $show_all -eq 1 ]; then
    display_local_ips
    display_gateway_ip
    display_external_ip
elif [ $show_local -eq 1 ]; then
    display_local_ips
fi
if [ $show_gateway -eq 1 ]; then
    display_gateway_ip
fi
if [ $show_external -eq 1 ]; then
    display_external_ip
fi

# Default behavior when no flags are provided
if [ $OPTIND -eq 1 ]; then
    display_local_ips
fi

