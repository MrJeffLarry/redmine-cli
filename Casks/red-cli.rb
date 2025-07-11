# This file was generated by GoReleaser. DO NOT EDIT.
cask "red-cli" do
  desc "Redmine CLI"
  homepage "https://github.com/mrjefflarry/redmine-cli"
  version "0.1.8-beta2"

  livecheck do
    skip "Auto-generated on release."
  end

  binary "red-cli"

  on_macos do
    on_intel do
      url "https://github.com/MrJeffLarry/redmine-cli/releases/download/v0.1.8-beta2/red-cli_0.1.8-beta2_darwin_amd64.tar.gz"
      sha256 "fa94b520f3248a18bc8701237b02a8d93f0c9ed2c4f79beebddb53d016f6ca8f"
    end
    on_arm do
      url "https://github.com/MrJeffLarry/redmine-cli/releases/download/v0.1.8-beta2/red-cli_0.1.8-beta2_darwin_arm64.tar.gz"
      sha256 "ff8fa569d497279e65f3a9b7b9ed931ae3538247bfe69aec242595214fd96101"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/MrJeffLarry/redmine-cli/releases/download/v0.1.8-beta2/red-cli_0.1.8-beta2_linux_amd64.tar.gz"
      sha256 "e2273f80034e9b9f655464569380cec34ca63dfb9040118edaa1287324630568"
    end
    on_arm do
      url "https://github.com/MrJeffLarry/redmine-cli/releases/download/v0.1.8-beta2/red-cli_0.1.8-beta2_linux_arm64.tar.gz"
      sha256 "bc23af27ad218d406b5b7ccf97b51737a20397cbe14575ce3c09f778150bfc21"
    end
  end

  postflight do
    if system_command("/usr/bin/xattr", args: ["-h"]).exit_status == 0
      # replace 'red-cli' with the actual binary name
      system_command "/usr/bin/xattr", args: ["-dr", "com.apple.quarantine", "#{staged_path}/red-cli"]
    end
  end

  # No zap stanza required
end
