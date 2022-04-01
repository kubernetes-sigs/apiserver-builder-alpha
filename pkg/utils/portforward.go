package utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/klog/v2"
)

type LocalProxyServer interface {
	Listen() (func(), error)
}

// PortForwardProtocolV1Name is the subprotocol used for port forwarding.
const PortForwardProtocolV1Name = "portforward.k8s.io"

func NewLocalProxy(
	restConfig *rest.Config,
	podNamespace,
	podName string,
	targetPort int32) LocalProxyServer {
	return &roundRobin{
		restConfig:   restConfig,
		podNamespace: podNamespace,
		podName:      podName,
		targetPort:   targetPort,
	}
}

var _ LocalProxyServer = &roundRobin{}

type roundRobin struct {
	podNamespace string
	podName      string
	targetPort   int32

	restConfig *rest.Config
}

func (r *roundRobin) Listen() (func(), error) {
	klog.V(4).Infof("Started local proxy server at port %d", r.targetPort)
	listener, err := net.Listen(
		"tcp",
		net.JoinHostPort("localhost", strconv.Itoa(int(r.targetPort))))
	if err != nil {
		return nil, fmt.Errorf("unable to create listener: Error %s", err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if !strings.Contains(strings.ToLower(err.Error()), "use of closed network connection") {
					runtime.HandleError(fmt.Errorf("error accepting connection on port %d: %v", r.targetPort, err))
				}
			}
			go func() {
				if err := r.handle(conn); err != nil {
					runtime.HandleError(fmt.Errorf("error handling connection: %v", err))
				}
			}()
		}
	}()
	return func() { listener.Close() }, nil

}

func (r *roundRobin) handle(conn net.Conn) error {
	klog.V(6).Info("Receiving connection")
	transport, upgrader, err := spdy.RoundTripperFor(r.restConfig)
	if err != nil {
		return err
	}
	nativeClient, err := kubernetes.NewForConfig(r.restConfig)
	if err != nil {
		return err
	}
	req := nativeClient.RESTClient().
		Post().
		Prefix("api", "v1").
		Resource("pods").
		Namespace(r.podNamespace).
		Name(r.podName).
		SubResource("portforward")
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", req.URL())
	streamConn, _, err := dialer.Dial(PortForwardProtocolV1Name)
	if err != nil {
		return err
	}
	defer streamConn.Close()

	// create error stream
	headers := http.Header{}
	headers.Set(v1.StreamType, v1.StreamTypeError)
	headers.Set(v1.PortHeader, fmt.Sprintf("%d", r.targetPort))
	headers.Set(v1.PortForwardRequestIDHeader, strconv.Itoa(0))

	errorStream, err := streamConn.CreateStream(headers)
	if err != nil {
		runtime.HandleError(fmt.Errorf("error creating error stream for port %d -> %d: %v", r.targetPort, r.targetPort, err))
		return err
	}
	// we're not writing to this stream
	errorStream.Close()

	errorChan := make(chan error)
	go func() {
		message, err := ioutil.ReadAll(errorStream)
		switch {
		case err != nil:
			errorChan <- fmt.Errorf("error reading from error stream for port %d -> %d: %v", r.targetPort, r.targetPort, err)
		case len(message) > 0:
			errorChan <- fmt.Errorf("an error occurred forwarding %d -> %d: %v", r.targetPort, r.targetPort, string(message))
		}
		close(errorChan)
	}()

	// create data stream
	headers.Set(v1.StreamType, v1.StreamTypeData)
	dataStream, err := streamConn.CreateStream(headers)
	if err != nil {
		runtime.HandleError(fmt.Errorf("error creating forwarding stream for port %d -> %d: %v", r.targetPort, r.targetPort, err))
		return err
	}

	localError := make(chan struct{})
	remoteDone := make(chan struct{})

	go func() {
		// Copy from the remote side to the local port.
		if _, err := io.Copy(conn, dataStream); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			runtime.HandleError(fmt.Errorf("error copying from remote stream to local connection: %v", err))
		}

		// inform the select below that the remote copy is done
		close(remoteDone)
	}()

	go func() {
		// inform server we're not sending any more data after copy unblocks
		defer dataStream.Close()

		// Copy from the local port to the remote side.
		if _, err := io.Copy(dataStream, conn); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			runtime.HandleError(fmt.Errorf("error copying from local connection to remote stream: %v", err))
			// break out of the select below without waiting for the other copy to finish
			close(localError)
		}
	}()

	// wait for either a local->remote error or for copying from remote->local to finish
	select {
	case <-remoteDone:
		klog.V(6).Info("Connection closed from remote")
	case <-localError:
		klog.V(6).Info("Connection closed due to local error")
	}

	return nil
}

var _ net.Conn = &conn{}

type conn struct {
	dataStream httpstream.Stream
}

func (c conn) Read(b []byte) (n int, err error) {
	return c.dataStream.Read(b)
}

func (c conn) Write(b []byte) (n int, err error) {
	return c.dataStream.Write(b)
}

func (c conn) Close() error {
	return c.dataStream.Close()
}

func (c conn) LocalAddr() net.Addr {
	return nil
}

func (c conn) RemoteAddr() net.Addr {
	return nil
}

func (c conn) SetDeadline(t time.Time) error {
	return errors.New("not implemented")
}

func (c conn) SetReadDeadline(t time.Time) error {
	return errors.New("not implemented")
}

func (c conn) SetWriteDeadline(t time.Time) error {
	return errors.New("not implemented")
}
