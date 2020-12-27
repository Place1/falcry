package api

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/place1/falcry/internal/webhook"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/place1/falcry/internal/storage"
	"github.com/place1/falcry/protos/protos"
	"google.golang.org/grpc"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
)

type ApiServer struct {
	protos.UnimplementedEventsServer
	db          *gorm.DB
	broadcaster *webhook.Broadcaster
}

func New(db *gorm.DB, broadcaster *webhook.Broadcaster) *ApiServer {
	return &ApiServer{
		db:          db,
		broadcaster: broadcaster,
	}
}

func (a *ApiServer) Handler() http.Handler {
	server := grpc.NewServer([]grpc.ServerOption{
		grpc.MaxRecvMsgSize(int(1 * math.Pow(2, 20))), // 1MB
		grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.StandardLogger()))(ctx, req, info, handler)
		}),
	}...)

	protos.RegisterEventsServer(server, a)

	return grpcweb.WrapServer(server,
		grpcweb.WithAllowNonRootResource(true),
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
	)
}

func (s *ApiServer) Listen(req *protos.ListenEventsReq, res protos.Events_ListenServer) error {
	query := s.db.
		WithContext(res.Context()).
		Order("time").
		Limit(100)

	if req.Since != nil {
		query = query.Where("time > ?", req.Since.AsTime())
	}

	if req.Until != nil {
		query = query.Where("time < ?", req.Until.AsTime())
	}

	var events []storage.Event
	if err := query.Find(&events).Error; err != nil {
		return err
	}

	for _, event := range events {
		res.Send(&protos.ListenEventsRes{
			Event: &protos.Event{
				Id:       fmt.Sprint(event.ID),
				Raw:      jsonstring(event.Raw),
				Output:   event.Output,
				Rule:     event.Rule,
				Priority: event.Priority,
				Time:     timestamp(event.Time),
			},
		})
	}

	if req.Until == nil {
		events := s.broadcaster.Channel(res.Context())
		for {
			select {
			case event := <-events:
				res.Send(&protos.ListenEventsRes{
					Event: &protos.Event{
						Id:       time.Now().String(),
						Raw:      event.Raw,
						Output:   event.Output,
						Rule:     event.Rule,
						Priority: event.Priority,
						Time:     timestamp(event.Time),
					},
				})
			case <-res.Context().Done():
				return nil
			}
		}
	}

	return nil
}

func jsonstring(j datatypes.JSON) string {
	b, err := j.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(b)
}
