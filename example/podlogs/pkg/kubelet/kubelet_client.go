package kubelet

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
)

type NodeGetter func(name string) (*v1.Node, error)

// KubeletClientConfig defines config parameters for the kubelet client
type KubeletClientConfig struct {
	// Port specifies the default port - used if no information about Kubelet port can be found in Node.NodeStatus.DaemonEndpoints.
	Port uint

	// ReadOnlyPort specifies the Port for ReadOnly communications.
	ReadOnlyPort uint

	// EnableHTTPs specifies if traffic should be encrypted.
	EnableHTTPS bool

	// PreferredAddressTypes - used to select an address from Node.NodeStatus.Addresses
	PreferredAddressTypes []string

	// TLSClientConfig contains settings to enable transport layer security
	restclient.TLSClientConfig

	// Server requires Bearer authentication
	BearerToken string

	// HTTPTimeout is used by the client to timeout http requests to Kubelet.
	HTTPTimeout time.Duration

	// Dial is a custom dialer used for the client
	Dial utilnet.DialFunc
}

// transportConfig converts a client config to an appropriate transport config.
func (c *KubeletClientConfig) transportConfig() *transport.Config {
	cfg := &transport.Config{
		TLS: transport.TLSConfig{
			CAFile:   c.CAFile,
			CAData:   c.CAData,
			CertFile: c.CertFile,
			CertData: c.CertData,
			KeyFile:  c.KeyFile,
			KeyData:  c.KeyData,
		},
		BearerToken: c.BearerToken,
	}
	if c.EnableHTTPS && !cfg.HasCA() {
		cfg.TLS.Insecure = true
	}
	return cfg
}

// NewNodeConnectionInfoGetter creates a new NodeConnectionInfoGetter.
func NewNodeConnectionInfoGetter(nodes NodeGetter, config KubeletClientConfig) (ConnectionInfoGetter, error) {
	scheme := "https"

	transport, err := MakeTransport(&config)
	if err != nil {
		return nil, err
	}

	types := []v1.NodeAddressType{}
	for _, t := range config.PreferredAddressTypes {
		types = append(types, v1.NodeAddressType(t))
	}

	return &NodeConnectionInfoGetter{
		nodes:       nodes,
		scheme:      scheme,
		defaultPort: int(config.Port),
		transport:   transport,

		preferredAddressTypes: types,
	}, nil
}

// NodeConnectionInfoGetter obtains connection info from the status of a Node API object
type NodeConnectionInfoGetter struct {
	// nodes is used to look up Node objects
	nodes NodeGetter
	// scheme is the scheme to use to connect to all kubelets
	scheme string
	// defaultPort is the port to use if no Kubelet endpoint port is recorded in the node status
	defaultPort int
	// transport is the transport to use to send a request to all kubelets
	transport http.RoundTripper
	// preferredAddressTypes specifies the preferred order to use to find a node address
	preferredAddressTypes []v1.NodeAddressType
}

// GetConnectionInfo retrieves connection info from the status of a Node API object.
func (k *NodeConnectionInfoGetter) GetConnectionInfo(ctx context.Context, nodeName types.NodeName) (*ConnectionInfo, error) {
	node, err := k.nodes(string(nodeName))
	if err != nil {
		return nil, err
	}

	// Find a kubelet-reported address, using preferred address type
	host, err := GetPreferredNodeAddress(node, k.preferredAddressTypes)
	if err != nil {
		return nil, err
	}

	// Use the kubelet-reported port, if present
	port := int(node.Status.DaemonEndpoints.KubeletEndpoint.Port)
	if port <= 0 {
		port = k.defaultPort
	}

	return &ConnectionInfo{
		Scheme:    k.scheme,
		Hostname:  host,
		Port:      strconv.Itoa(port),
		Transport: k.transport,
	}, nil
}

// MakeTransport creates a RoundTripper for HTTP Transport.
func MakeTransport(config *KubeletClientConfig) (http.RoundTripper, error) {
	tlsConfig, err := transport.TLSConfigFor(config.transportConfig())
	if err != nil {
		return nil, err
	}

	rt := http.DefaultTransport
	dialer := config.Dial

	if dialer != nil || tlsConfig != nil {
		// If SSH Tunnel is turned on
		rt = utilnet.SetOldTransportDefaults(&http.Transport{
			DialContext:     dialer,
			TLSClientConfig: tlsConfig,
		})
	}

	return transport.HTTPWrappersForConfig(config.transportConfig(), rt)
}

// GetPreferredNodeAddress returns the address of the provided node, using the provided preference order.
// If none of the preferred address types are found, an error is returned.
func GetPreferredNodeAddress(node *v1.Node, preferredAddressTypes []v1.NodeAddressType) (string, error) {
	for _, addressType := range preferredAddressTypes {
		for _, address := range node.Status.Addresses {
			if address.Type == addressType {
				return address.Address, nil
			}
		}
	}
	return "", fmt.Errorf("no match")
}

// ConnectionInfo provides the information needed to connect to a kubelet
type ConnectionInfo struct {
	Scheme    string
	Hostname  string
	Port      string
	Transport http.RoundTripper
}

// ConnectionInfoGetter provides ConnectionInfo for the kubelet running on a named node
type ConnectionInfoGetter interface {
	GetConnectionInfo(ctx context.Context, nodeName types.NodeName) (*ConnectionInfo, error)
}
