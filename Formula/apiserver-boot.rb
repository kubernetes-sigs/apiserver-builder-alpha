class ApiserverBoot < Formula
    desc "Scaffolding Toolkit for APIServer Aggregation"
    homepage "https://github.com/kubernetes-sigs/apiserver-builder-alpha"
    url "https://github.com/kubernetes-sigs/apiserver-builder-alpha.git",
      :using => :git,
      :tag => "v1.18.0",
      :revision => "95dca1d34e91d6e76c50fa4f272a77f573fd7558"
    depends_on "bazel" => :build

    def install
      if OS.mac?
        system "bazel","build","--platforms=@io_bazel_rules_go//go/toolchain:darwin_amd64","cmd:apiserver-builder"
      elsif OS.linux?
        system "bazel","build","--platforms=@io_bazel_rules_go//go/toolchain:linux_amd64","cmd:apiserver-builder"
      end
      system "bash","-c","tar xzvf bazel-bin/cmd/apiserver-builder.tar.gz"
      bin.install "bin/apiserver-boot"
    end

    test do
      assert_match version.to_s, shell_output("#{bin}/apiserver-boot version")
    end
  end
