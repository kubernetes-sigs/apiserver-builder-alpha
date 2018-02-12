# Contributing

## Project structure
 
### cmd/kubebuilder

Convenience commands for bootstrapping a project.  Totally optional, but helpful.
 
- Initialize a new repo with vendored deps and project structure
- Bootstrap types.go, types_test.go, controllers and sample yamls
- Run code generators
- Compile and run controlplane locally

### cmd/kubebuilder-gen

Generate lots of boilerplate and wiring one would have to do by hand otherwise

### pkg/

Libraries for defining common implementations of resources
