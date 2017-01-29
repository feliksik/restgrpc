package restgrpc

import (
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// grpcEndpoint defines a grpc endpoint, and a RegisterHandler function
// that registers an endpoint to be used as http gateway
type grpcEndpoint struct {
	Name                string
	Endpoint            string
	RegisterRestHandler func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
}

type gateway struct {
	grpcEndpoints []grpcEndpoint
	mux           *runtime.ServeMux
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewEndpoint creates a new GRPC endpoint, with it's associated function to Register the
// REST handler, as generated by the grpc-gateway package
func NewEndpoint(name string, endpoint string, registerHandler func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error) grpcEndpoint {
	return grpcEndpoint{
		Name:                name,
		Endpoint:            endpoint,
		RegisterRestHandler: registerHandler,
	}
}

func NewGateway() *gateway {
	p := &gateway{}

	ctx := context.Background()
	p.ctx, p.cancel = context.WithCancel(ctx)
	p.grpcEndpoints = make([]grpcEndpoint, 0)
	p.mux = runtime.NewServeMux()

	return p
}

func (p *gateway) AddService(endpoint grpcEndpoint) {
	p.grpcEndpoints = append(p.grpcEndpoints, endpoint)
}

// connectToEndpoint connects to the remote grpc endpoint and registers
// a local REST endpoint on the gateway mux
func (p *gateway) connectToEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint *grpcEndpoint) error {

	opts := []grpc.DialOption{grpc.WithInsecure()}

	conn, err := grpc.Dial(endpoint.Endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", endpoint.Endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", endpoint.Endpoint, cerr)
			}
		}()
	}()

	endpoint.RegisterRestHandler(ctx, mux, conn)

	// unfortunately we cannot easily print the REST endpoint name (/v1/compute/compute)
	// as it is private in the generated pb.gw.go file
	log.Printf("connected to grpc endpoint  %s for service %s\n", endpoint.Endpoint, endpoint.Name)
	return nil
}

// connectToEndpoints connects to all registered grpc endpoints
func (p *gateway) connectToEndpoints() error {

	for _, e := range p.grpcEndpoints {
		err := p.connectToEndpoint(p.ctx, p.mux, &e)
		if err != nil {
			return err
		}
	}

	return nil
}

// Bind connects on all remote andpoints and
// binds on the given socket, e.g. ":8080".
// It then blocks.
func (p *gateway) Bind(socket string) error {
	defer p.cancel()

	if err := p.connectToEndpoints(); err != nil {
		return err
	}
	// listen
	return http.ListenAndServe(socket, p.mux)

}
