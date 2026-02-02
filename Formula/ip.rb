class Ip < Formula
  desc "Simplest cli tool to get your IP (local, external, gateway)"
  homepage "https://github.com/adriangalilea/homebrew-ip"
  version "2.1.0"
  
  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.1.0/ip-darwin-arm64"
      sha256 "dede5f9e02499c24be1add7fe508968e6fe0f3f3c57a81cf08d799130e7b75dc"
    else
      url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.1.0/ip-darwin-amd64"
      sha256 "2bae10c5e07ba9eebd5dc3cc1e68b71c164664ba452cada4c8b2f758f2ad1745"
    end
  elsif OS.linux?
    url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.1.0/ip-linux-amd64"
    sha256 "2a49d5b5185d212fbfe02a0809ea732b580be9ca6d41623166053e8d0a22f90a"
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
