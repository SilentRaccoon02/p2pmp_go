package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/SilentRaccoon02/p2pmp_go/common"
	"github.com/hashicorp/yamux"
	"github.com/pion/randutil"
	"github.com/pion/webrtc/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OfferClientServerImpl struct {
	ICEMutex    *sync.Mutex
	pendingICE  []*webrtc.ICECandidate
	peerConn    *webrtc.PeerConnection
	dataChannel *webrtc.DataChannel
	config      webrtc.Configuration

	client common.P2PManagerClient
	common.UnimplementedP2PClientServer
}

type AnswerClientServerImpl struct {
	ICEMutex   *sync.Mutex
	pendingICE []*webrtc.ICECandidate
	peerConn   *webrtc.PeerConnection
	config     webrtc.Configuration

	client common.P2PManagerClient
	common.UnimplementedP2PClientServer
}

func (server OfferClientServerImpl) CheckClientType(ctx context.Context, empty *common.Empty) (*common.ClientType, error) {
	return &common.ClientType{Type: "offer"}, nil
}

func (server AnswerClientServerImpl) CheckClientType(ctx context.Context, empty *common.Empty) (*common.ClientType, error) {
	return &common.ClientType{Type: "answer"}, nil
}

func (server *OfferClientServerImpl) StartOffer(ctx context.Context, empty *common.Empty) (*common.Empty, error) {
	log.Println("Offer connection started")

	server.ICEMutex = &sync.Mutex{}
	server.pendingICE = make([]*webrtc.ICECandidate, 0)

	server.config = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	var err error
	server.peerConn, err = webrtc.NewPeerConnection(server.config)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := server.peerConn.Close(); err != nil {
			log.Printf("Cannot close peer connection: %v\n", err)
		}
	}()

	server.peerConn.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		server.ICEMutex.Lock()
		defer server.ICEMutex.Unlock()

		desc := server.peerConn.RemoteDescription()

		if desc == nil {
			server.pendingICE = append(server.pendingICE, c)
		} else {
			server.client.AddICEOffer(context.Background(), &common.ICECandidate{Candidate: c.ToJSON().Candidate})
		}
	})

	server.dataChannel, err = server.peerConn.CreateDataChannel("data", nil)

	if err != nil {
		panic(err)
	}

	server.peerConn.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Printf("Peer connection state changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			log.Println("Peer connection has gone to failed exiting")
			os.Exit(0)
		}

		if s == webrtc.PeerConnectionStateClosed {
			log.Println("Peer connection has gone to closed exiting")
			os.Exit(0)
		}
	})

	server.dataChannel.OnOpen(func() {
		log.Printf("Data channel open: %d\n", server.dataChannel.ID())

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			message, err := randutil.GenerateCryptoRandomString(10, "012345689")

			if err != nil {
				panic(err)
			}

			log.Printf("Sending: %s\n", message)

			err = server.dataChannel.SendText(message)

			if err != nil {
				panic(err)
			}
		}
	})

	server.dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		log.Printf("Message: '%s'\n", string(msg.Data))
	})

	offer, err := server.peerConn.CreateOffer(nil)

	if err != nil {
		panic(err)
	}

	err = server.peerConn.SetLocalDescription(offer)

	if err != nil {
		panic(err)
	}

	payload, err := json.Marshal(offer)

	if err != nil {
		panic(err)
	}

	server.client.AddSessionDescOffer(context.Background(), &common.SessionDesc{Desc: string(payload)})

	select {}
}

func (server *AnswerClientServerImpl) StartAnswer(ctx context.Context, empty *common.Empty) (*common.Empty, error) {
	log.Println("Answer connection started")

	server.ICEMutex = &sync.Mutex{}
	server.pendingICE = make([]*webrtc.ICECandidate, 0)

	server.config = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	var err error
	server.peerConn, err = webrtc.NewPeerConnection(server.config)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := server.peerConn.Close(); err != nil {
			log.Printf("Cannot close peer connection: %v\n", err)
		}
	}()

	server.peerConn.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		server.ICEMutex.Lock()
		defer server.ICEMutex.Unlock()

		desc := server.peerConn.RemoteDescription()

		if desc == nil {
			server.pendingICE = append(server.pendingICE, c)
		} else {
			server.client.AddICEAnswer(context.Background(), &common.ICECandidate{Candidate: c.ToJSON().Candidate})
		}
	})

	server.peerConn.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Printf("Peer connection state changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			log.Println("Peer connection has gone to failed exiting")
			os.Exit(0)
		}

		if s == webrtc.PeerConnectionStateClosed {

			log.Println("Peer connection has gone to closed exiting")
			os.Exit(0)
		}
	})

	server.peerConn.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnOpen(func() {
			log.Printf("Data channel open: %d\n", d.ID())
		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			log.Printf("Message: %s\n", string(msg.Data))
		})
	})

	select {}
}

func (server *OfferClientServerImpl) AddICEOffer(ctx context.Context, candidate *common.ICECandidate) (*common.Empty, error) {
	log.Printf("Add ICE %s\n", candidate.Candidate)

	server.peerConn.AddICECandidate(webrtc.ICECandidateInit{Candidate: candidate.Candidate})

	return &common.Empty{}, nil
}

func (server *AnswerClientServerImpl) AddICEAnswer(ctx context.Context, candidate *common.ICECandidate) (*common.Empty, error) {
	log.Printf("Add ICE %s\n", candidate.Candidate)

	server.peerConn.AddICECandidate(webrtc.ICECandidateInit{Candidate: candidate.Candidate})

	return &common.Empty{}, nil
}

func (server *OfferClientServerImpl) AddSessionDescOffer(ctx context.Context, desc *common.SessionDesc) (*common.Empty, error) {
	log.Println("Add session description")

	sdp := webrtc.SessionDescription{}
	json.Unmarshal([]byte(desc.Desc), &sdp)

	err := server.peerConn.SetRemoteDescription(sdp)

	if err != nil {
		panic(err)
	}

	server.ICEMutex.Lock()
	defer server.ICEMutex.Unlock()

	for _, c := range server.pendingICE {
		server.client.AddICEOffer(context.Background(), &common.ICECandidate{Candidate: c.ToJSON().Candidate})
	}

	return &common.Empty{}, nil
}

func (server *AnswerClientServerImpl) AddSessionDescAnswer(ctx context.Context, desc *common.SessionDesc) (*common.Empty, error) {
	log.Println("Add session description")

	sdp := webrtc.SessionDescription{}
	json.Unmarshal([]byte(desc.Desc), &sdp)

	err := server.peerConn.SetRemoteDescription(sdp)

	if err != nil {
		panic(err)
	}

	answer, err := server.peerConn.CreateAnswer(nil)

	if err != nil {
		panic(err)
	}

	payload, err := json.Marshal(answer)

	if err != nil {
		panic(err)
	}

	server.client.AddSessionDescAnswer(context.Background(), &common.SessionDesc{Desc: string(payload)})

	err = server.peerConn.SetLocalDescription(answer)

	if err != nil {
		panic(err)
	}

	server.ICEMutex.Lock()
	defer server.ICEMutex.Unlock()

	for _, c := range server.pendingICE {
		server.client.AddICEAnswer(context.Background(), &common.ICECandidate{Candidate: c.ToJSON().Candidate})
	}

	return &common.Empty{}, nil
}

func offerClient() {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:3001", 10*time.Second)

	if err != nil {
		panic(err)
	}

	yamuxSession, err := yamux.Server(conn, yamux.DefaultConfig())

	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	clientServerImpl := &OfferClientServerImpl{}
	clientServerImpl.client = connectServer()
	common.RegisterP2PClientServer(grpcServer, clientServerImpl)

	log.Println("Launching gRPC over TCP")

	err = grpcServer.Serve(yamuxSession)

	if err != nil {
		panic(err)
	}
}

func answerClient() {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:3001", 10*time.Second)

	if err != nil {
		panic(err)
	}

	yamuxSession, err := yamux.Server(conn, yamux.DefaultConfig())

	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	clientServerImpl := &AnswerClientServerImpl{}
	clientServerImpl.client = connectServer()
	common.RegisterP2PClientServer(grpcServer, clientServerImpl)

	log.Println("Launching gRPC over TCP")

	err = grpcServer.Serve(yamuxSession)

	if err != nil {
		panic(err)
	}
}

func connectServer() common.P2PManagerClient {
	conn, err := grpc.NewClient("127.0.0.1:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	log.Println("Server connected")

	return common.NewP2PManagerClient(conn)
}

func main() {
	if len(os.Args) < 2 {
		log.Println("Invalid params")
		return
	}

	clientType := os.Args[1]

	if clientType == "off" {
		log.Println("Offer client")
		offerClient()
	} else if clientType == "ans" {
		log.Println("Answer client")
		answerClient()
	} else {
		log.Println("Invalid params")
		return
	}
}
