class Ip < Formula
  desc "Simplest cli tool to get your IP (local, external, gateway)"
  homepage "https://github.com/adriangalilea/homebrew-ip"
  version "2.0.0"
  
  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.0.0/ip-darwin-arm64"
      sha256 "0613db2d47ca6331ee1724414751e26ca9ed7e064e242ac8901596407f1323e4"
    else
      url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.0.0/ip-darwin-amd64"
      sha256 "3a0c33ff5a1b6765f79ec7602d5832507cdf9a1911901da030a5201df7cdd554"
    end
  elsif OS.linux?
    url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.0.0/ip-linux-amd64"
    sha256 "ebb8625a64fc75e147c4d3619214b452df82b578a380bddab6b4bc3eb8301724"
  end

  def install
    bin.install "ip-darwin-#{Hardware::CPU.arm? ? "arm64" : "amd64"}" => "ip" if OS.mac?
    bin.install "ip-linux-amd64" => "ip" if OS.linux?
  end

  test do
    output = shell_output("#{bin}/ip")
    assert_match(/\d+\.\d+\.\d+\.\d+/, output)
  end
end
