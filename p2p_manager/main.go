package main

import (
	"context"
	"log"
	"net"

	"github.com/SilentRaccoon02/p2pmp_go/common"
	"github.com/hashicorp/yamux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type P2PManagerServerImpl struct {
	offerClient  common.P2PClientClient
	answerClient common.P2PClientClient
	common.UnimplementedP2PManagerServer
}

func (server *P2PManagerServerImpl) AddICEOffer(ctx context.Context, candidate *common.ICECandidate) (*common.Empty, error) {
	if server.answerClient == nil {
		panic("Answer client == nil")
	}

	log.Println("ICE offer -> answer")
	server.answerClient.AddICEAnswer(context.Background(), candidate)
	return &common.Empty{}, nil
}

func (server *P2PManagerServerImpl) AddICEAnswer(ctx context.Context, candidate *common.ICECandidate) (*common.Empty, error) {
	if server.offerClient == nil {
		panic("Offer client == nil")
	}

	log.Println("ICE answer -> offer")
	server.offerClient.AddICEOffer(context.Background(), candidate)
	return &common.Empty{}, nil
}

func (server *P2PManagerServerImpl) AddSessionDescOffer(ctx context.Context, desc *common.SessionDesc) (*common.Empty, error) {
	if server.answerClient == nil {
		panic("Answer client == nil")
	}

	log.Println("SDP offer -> answer")
	server.answerClient.AddSessionDescAnswer(context.Background(), desc)
	return &common.Empty{}, nil
}

func (server *P2PManagerServerImpl) AddSessionDescAnswer(ctx context.Context, desc *common.SessionDesc) (*common.Empty, error) {
	if server.offerClient == nil {
		panic("Offer client == nil")
	}

	log.Println("SDP answer -> offer")
	server.offerClient.AddSessionDescOffer(context.Background(), desc)
	return &common.Empty{}, nil
}

func startP2PManagerServer() {
	listener, err := net.Listen("tcp", ":3000")

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	grpcServer := grpc.NewServer()
	p2pManagerServer := &P2PManagerServerImpl{}
	common.RegisterP2PManagerServer(grpcServer, p2pManagerServer)

	go startYamuxServer(p2pManagerServer)

	log.Println("gRPC listening on port 3000")
	err = grpcServer.Serve(listener)

	if err != nil {
		panic(err)
	}
}

func startYamuxServer(p2pManagerServer *P2PManagerServerImpl) {
	listener, err := net.Listen("tcp", ":3001")

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	log.Println("Yamux listening on port 3001")

	for {
		conn, err := listener.Accept()

		if err != nil {
			panic(err)
		}

		yamuxSession, err := yamux.Client(conn, yamux.DefaultConfig())

		if err != nil {
			panic(err)
		}

		log.Println("Launching gRPC over TCP")

		clientConn, err := grpc.NewClient(":3000", grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return yamuxSession.Open() }))

		if err != nil {
			panic(err)
		}

		go handleConn(clientConn, p2pManagerServer)
	}
}

func handleConn(clientConn *grpc.ClientConn, p2pManagerServer *P2PManagerServerImpl) {
	client := common.NewP2PClientClient(clientConn)
	res, err := client.CheckClientType(context.Background(), &common.Empty{})

	if err != nil {
		panic(err)
	}

	if res.Type == "offer" {
		log.Println("Offer client connected")
		p2pManagerServer.offerClient = client
		p2pManagerServer.offerClient.StartOffer(context.Background(), &common.Empty{})
	} else if res.Type == "answer" {
		log.Println("Answer client connected")
		p2pManagerServer.answerClient = client
		p2pManagerServer.answerClient.StartAnswer(context.Background(), &common.Empty{})
	} else {
		panic("Incorrect client type")
	}
}

func main() {
	startP2PManagerServer()
}
