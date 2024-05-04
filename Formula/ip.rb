class Ip < Formula
  desc "Simplest cli tool to get your IP (local, external, gateway)"
  homepage "https://github.com/adriangalilea/homebrew-ip"
  url "https://github.com/adriangalilea/homebrew-ip/archive/refs/tags/v1.0.1.tar.gz" 
  sha256 "e2fd2d182c88a4fa8faeb8628c517cedeb606692c128c26c021e2fd95e48d205" 

  def install
    bin.install "ip.sh" => "ip"
  end

  test do
    system "#{bin}/ip", "-h"
  end
end
