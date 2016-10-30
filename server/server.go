package main

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"log"
	"net"
	"runtime"
	"time"

	pb "github.com/imgcompress/lossycompress"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
	"google.golang.org/grpc"
)

// Limter stuff
var (
	GoroutineNum = runtime.NumCPU() * 2
	limiterChan  = make(chan struct{}, GoroutineNum)
)

const (
	projectId  = "imgresizer-service"
	port       = ":3030"
	bucketName = "imgcompressor_bucket"
	version    = 0.1
)

func Limiter(fn func() error, wait time.Duration, errs chan error) {
	limiterChan <- struct{}{}
	go func() {
		defer func() {
			<-limiterChan
		}()

		err := fn()
		if err != nil {
			errs <- err
		}

		time.Sleep(wait)
	}()
}

type Server struct {
	Store     *storage.Service
	StoreName string
}

func (s *Server) Compress(ctx context.Context, req *pb.Request) (*pb.Result, error) {
	reader := bytes.NewBufferString(req.Data)
	in, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	opt := &jpeg.Options{Quality: int(req.Quality)}
	err = jpeg.Encode(out, in, opt)
	if err != nil {
		return nil, err
	}
	// decide where to put the image
	link, err := s.upload(out, req.Name)
	if err == nil {
		err = errors.New("")
	}

	return &pb.Result{link, req.Name}, nil
}

func (s *Server) upload(data *bytes.Buffer, name string) (string, error) {
	var (
		object = &storage.Object{Name: name}
	)

	res, err := s.Store.Objects.Insert(s.StoreName, object).Media(data).Do()
	if err != nil {
		log.Printf("Failed to upload %v\n", name)
	} else {
		log.Printf("Uploaded %v\n", name)
	}

	return res.SelfLink, err
}

var s *Server

func init() {
	const scope = storage.DevstorageFullControlScope
	client, err := google.DefaultClient(context.Background(), scope)
	if err != nil {
		log.Fatalf("Failed to initialize the default client %v\n", err)
	}

	service, err := storage.New(client)
	if err != nil {
		log.Fatalf("Failed to initialize the storage: %v\n", err)
	}

	s = &Server{service, bucketName}

	// bucket initialization stuff
	if _, err := s.Store.Buckets.Get(s.StoreName).Do(); err == nil {
		log.Println("Store bucket already there")
	} else {
		bucket := &storage.Bucket{Name: s.StoreName}
		res, err := s.Store.Buckets.Insert(projectId, bucket).Do()
		if err == nil {
			log.Printf("Created bucket %v\n", res.Name)
		} else {
			log.Fatalf("Failed creating bucket %s: %v\n", s.StoreName, err)
		}
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterImgServer(grpcServer, s)
	grpcServer.Serve(listener)
}
