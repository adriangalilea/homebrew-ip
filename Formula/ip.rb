class Ip < Formula
  desc "Simplest cli tool to get your IP (local, external, gateway)"
  homepage "https://github.com/adriangalilea/ip"
  url "https://github.com/adriangalilea/ip/archive/refs/tags/v1.0.tar.gz" 
  sha256 "fc8a1b3e9c0ee97ad3409431ce1a26f90423be70ee69f62915e4ffa7f5fd6f39" 

  def install
    bin.install "ip.sh" => "ip"
  end

  test do
    system "#{bin}/ip", "-h"
  end
end
