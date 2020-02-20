class ServerGo < Formula
  desc "OpenIoTHub Server Client"
  homepage "https://github.com/OpenIoTHub/server-go"
  url "https://github.com/OpenIoTHub/server-go.git",
      :tag      => "v1.1.29",
      :revision => "68f2e56990e4ffb3318efa709fcdb7d4ddeb6e9a"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags",
             "-s -w -X main.version=#{version} -X main.commit=#{stable.specs[:revision]} -X main.builtBy=homebrew",
             "-o", bin/"server-go"
  end

  test do
    system "#{bin}/server-go", "-v"
  end
end
