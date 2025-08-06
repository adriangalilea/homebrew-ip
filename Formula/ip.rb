class Ip < Formula
  desc "Simplest cli tool to get your IP (local, external, gateway)"
  homepage "https://github.com/adriangalilea/homebrew-ip"
  version "2.0.1"
  
  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.0.1/ip-darwin-arm64"
      sha256 "23e9537ec52b6c1cb1809ab2bbbbee76387ec738443bcbcda13b95f93d51ebe6"
    else
      url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.0.1/ip-darwin-amd64"
      sha256 "f96731cf85b16f632f2182f1a060eee26e1e464c4d4914ce91d330d22ee80774"
    end
  elsif OS.linux?
    url "https://github.com/adriangalilea/homebrew-ip/releases/download/v2.0.1/ip-linux-amd64"
    sha256 "0659743720c2992b0443905713028be75496d69e21bdb71e4e96aa44b5378468"
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
