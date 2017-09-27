# Kubernetes Reference Docs

Tools to build reference documentation for Kubernetes APIs and CLIs.

# Api Docs

## Generate new api docs

1. From the kubernetes/kubernetes repo, copy the file `k8s.io/kubernetes/api/openapi-spec/swagger.json` to `gen_open_api/openapi-spec/swagger.json` in the reference-docs repo.

2. Update the file `gen_open_api/config.yaml`, adding any new resource definitions or operation types not already present
  - TODO: Write more on this

3. Run `make api` to build the doc html and javascript

4. Html files will be written to `gen_open_api/build`.  Copy these to where they will be hosted.

# Cli

## Generate new kubectl docs

1. Regenerate the `kubectl.yaml` file with the command metadata
  - Update `gen_kubectl/kubectl.yaml` by running `k8s.io/kubernetes/cmd/genslateyaml/gen_slate_yaml.go` from the kubernetes/kuberentes repo
  - Make sure you are on the kubuernetes branch that you want to generate docs for
  - Copy the file to the [correct directory](https://github.com/kubernetes-incubator/reference-docs/blob/master/gen-kubectldocs/generators/v1_6/kubectl.yaml) for the kubernetes version

2. Update the Makefile for your environment
  - Set `K8SIOROOT` to the kubernetes/kubernetes.io root directory
  - Set `K8SROOT` to the kubernetes/kubernetes root directory
  - Update the kubernetes version to match the kubernetes branch
    - `v1_6`
    - `v1.6`

3. Build the cli docs
  - Run `make cli`
  - Files will be copied under `gen-kubectldocs/generators/build/`
  - Open up `index.html` in a browser and make sure it looks good

4. Copy the cli docs to kubernetes/kubernetes.io
  - Run `make copycli`
  - Files will be copied to the appropriate directory under `K8SIOROOT`
  - You may need to create a new directory for new kubernetes versions

# Updating brodocs version

*May need to change the image repo to one you have write access to.*

1. Update Dockerfile so it will re-clone the repo

2. Run `make brodocs`
