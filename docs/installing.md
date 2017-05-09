# Installing the apiserver build tools

- Download the latest [release](https://github.com/kubernetes-incubator/apiserver-builder/releases)
- Create a new directory for the binaries:
  `sudo mkdir /usr/local/apiserver-builder/`
- Unpack the release `tar -xzvf <apiserver-builder.tar>`
- Move release into the directory
  `sudo mv <apiserver-buidler/*> /usr/local/apiserver-builder/`
- Add the binaries to your *PATH*
  `export PATH=$PATH:/usr/local/apiserver-builder/bin`
- Test things are working by running `apiserver-boot -h`
