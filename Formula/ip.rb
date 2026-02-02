class Ip < Formula
  desc "Simplest cli tool to get your IP (local, external, gateway)"
  homepage "https://github.com/adriangalilea/homebrew-ip"
  version "2.1.1"
  
  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.1.1/ip-darwin-arm64"
      sha256 "aab16c5d942112206ef2055efc96b1ccc0265d1f99737f6dd5ca3356805a651a"
    else
      url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.1.1/ip-darwin-amd64"
      sha256 "25fb05329c7feecd22959fb4d322fe45cc57a7943f186bad2591a0a435b0f3f0"
    end
  elsif OS.linux?
    url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.1.1/ip-linux-amd64"
    sha256 "94f4e0339b702edf75601191d30c4368eb387a001957dd91af64cefa281114bb"
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
