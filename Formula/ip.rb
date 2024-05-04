class Ip < Formula
  desc "Simplest cli tool to get your IP (local, external, gateway)"
  homepage "https://github.com/adriangalilea/ip"
  url "https://github.com/adriangalilea/ip/archive/refs/tags/v1.0.tar.gz" 
  sha256 "bf5817134faa4b90a2aa99cdd8a61e708a680cf563f5e076df9678e36e22e622" 

  def install
    bin.install "ip.sh" => "ip"
  end

  test do
    system "#{bin}/ip", "-h"
  end
end
