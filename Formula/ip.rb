class Ip < Formula
  desc "Simplest cli tool to get your IP (local, external, gateway)"
  homepage "https://github.com/adriangalilea/homebrew-ip"
  url "https://github.com/adriangalilea/homebrew-ip/archive/refs/heads/main.tar.gz"
  version "2.0.0"
  head "https://github.com/adriangalilea/homebrew-ip.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"ip", "."
  end

  test do
    output = shell_output("#{bin}/ip")
    assert_match(/\d+\.\d+\.\d+\.\d+/, output)
  end
end