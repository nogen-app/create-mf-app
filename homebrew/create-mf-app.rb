require 'rbconfig'
class Create-Mf-App < Formula
  desc ""
  homepage "https://github.com/nogen/create-mf-app"
  version "{{.Version}}"

  if Hardware::CPU.is_64_bit?
    case RbConfig::CONFIG['host_os']
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      url "https://github.com/nogen/create-mf-app/releases/download/v{{.Version}}/{{.Mac64.FileName}}"
      sha256 "{{.Mac64.Hash}}"
    when /linux/
      url "https://github.com/nogen/create-mf-app/releases/download/v{{.Version}}/{{.Linux64.FileName}}"
      sha256 "{{.Linux64.Hash}}"
    when /solaris|bsd/
      :unix
    else
      :unknown
    end
  else
    case RbConfig::CONFIG['host_os']
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      :mac
    when /linux/
      url "https://github.com/nogen/create-mf-app/releases/download/v{{.Version}}/{{.Linux386.FileName}}"
      sha256 "{{.Linux386.Hash}}"
    when /solaris|bsd/
      :unix
    else
      :unknown
    end
  end

  def install
    bin.install "create-mf-app"
  end

  test do
    system "create-mf-app"
  end

end
