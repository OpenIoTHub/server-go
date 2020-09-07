class ServerGo < Formula
  desc "Server for OpenIoTHub"
  homepage "https://github.com/OpenIoTHub/server-go"
  url "https://github.com/OpenIoTHub/server-go.git",
      :tag      => "v1.1.58",
      :revision => "4c4f27fa71996d186ad10e0dee89c0619f9b978f"

  depends_on "go" => :build

  def install
    (etc/"server-go").mkpath
    system "go", "build", "-mod=vendor", "-ldflags",
             "-s -w -X main.version=#{version} -X main.commit=#{stable.specs[:revision]} -X main.builtBy=homebrew",
             *std_go_args
    etc.install "server-go.yaml" => "server-go/server-go.yaml"
  end

  plist_options manual: "server-go -c #{HOMEBREW_PREFIX}/etc/server-go/server-go.yaml"

  def plist
    <<~EOS
      <?xml version="1.0" encoding="UTF-8"?>
      <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
      <plist version="1.0">
        <dict>
          <key>Label</key>
          <string>#{plist_name}</string>
          <key>KeepAlive</key>
          <true/>
          <key>ProgramArguments</key>
          <array>
            <string>#{opt_bin}/server-go</string>
            <string>-c</string>
            <string>#{etc}/server-go/server-go.yaml</string>
          </array>
          <key>StandardErrorPath</key>
          <string>#{var}/log/server-go.log</string>
          <key>StandardOutPath</key>
          <string>#{var}/log/server-go.log</string>
        </dict>
      </plist>
    EOS
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/server-go -v 2>&1")
    assert_match "config created", shell_output("#{bin}/server-go init --config=server.yml 2>&1")
    assert_predicate testpath/"server.yml", :exist?
  end
end
