class Ip < Formula
  desc "Simplest cli tool to get your IP (local, external, gateway)"
  homepage "https://github.com/adriangalilea/homebrew-ip"
  version "2.2.0"
  
  if Hardware::CPU.arm?
    url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.2.0/ip-darwin-arm64"
    sha256 "8f006266fca24d9a636871f0f5ea566e30c9ba818dc3484a6b604d83a960680e"
  else
    url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.2.0/ip-darwin-amd64"
    sha256 "05136981b7afbee4747e344c75d61aaf1a24ce77ea3957e5f8c3ba12bccec53d"
  end

  def install
    bin.install "ip-darwin-#{Hardware::CPU.arm? ? "arm64" : "amd64"}" => "ip"
  end

  test do
    output = shell_output("#{bin}/ip")
    assert_match(/\d+\.\d+\.\d+\.\d+/, output)
  end
end
