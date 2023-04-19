package beacon

import (
	"context"
	"fmt"
	"strings"
	"time"

	api "github.com/troydai/grpcbeacon/gen/api/protos"
	"github.com/troydai/grpcbeacon/internal/settings"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	api.UnimplementedBeaconServer

	details map[string]string
	logger  *zap.Logger
}

var _ api.BeaconServer = (*Server)(nil)

const (
	_headerForwardedClientCert = "x-forwarded-client-cert"
)

func NewServer(env settings.Environment, logger *zap.Logger) *Server {
	s := &Server{
		details: make(map[string]string),
		logger:  logger,
	}

	s.details["Hostname"] = env.HostName
	s.details["Flockname"] = env.FlockName

	return s
}

func (s *Server) Signal(ctx context.Context, req *api.SignalReqeust) (*api.SignalResponse, error) {
	logger := s.logger
	if client, found := extractClient(ctx); found {
		logger = s.logger.With(zap.String("client", client))
	} else {
		logger = s.logger.With(zap.String("client", "unknown"))
	}

	logger.Info("Signal received")

	resp := &api.SignalResponse{
		Reply: fmt.Sprintf("Beacon signal at %s", time.Now().Format(time.RFC1123)),
	}

	if len(s.details) > 0 {
		resp.Details = s.details
	}

	return resp, nil
}

func extractClient(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	var clients []string
	for _, cert := range md[_headerForwardedClientCert] {
		for _, part := range strings.Split(cert, ";") {
			if strings.HasPrefix(part, "URI=") {
				clients = append(clients, part[4:])
			}
		}
	}

	if len(clients) <= 0 {
		return "", false
	}

	return strings.Join(clients, ","), true
}
