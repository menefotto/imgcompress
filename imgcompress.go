package imgresizer

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"runtime"
	"time"

	"golang.org/x/net/context"

	"golang.org/x/oauth2/google"

	storage "google.golang.org/api/storage/v1"
)

// Limter stuff
var (
	GoroutineNum = runtime.NumCPU() * 2
	limiterChan  = make(chan struct{}, GoroutineNum)
)

const (
	projectId = "imgcompressor"
	port      = ":3030"
	bucktName = "imgcompressor_bucket"
	version   = 0.1
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

func (s *Server) ImgCompress(img string, quality int, name string) (string, error) {
	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(img))
	in, _, err := image.Decode(reader)
	if err != nil {
		return "", err
	}

	out := &bytes.Buffer{}
	opt := &jpeg.Options{Quality: quality}
	err = jpeg.Encode(out, in, opt)
	if err != nil {
		return "", err
	}
	// decide where to put the image
	return s.upload(out, name)
}

func (s *Server) upload(data *bytes.Buffer, name string) (string, error) {
	var (
		object = &storage.Object{Name: name}
	)

	res, err := s.Store.Objects.Insert(s.StoreName, object).Media(data).Do()
	if err != nil {
		log.Println("Failed to upload %v\n", name)
	} else {
		log.Println("Uploaded %v\n", name)
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
			log.Fatalf("Failed creating bucket %s: %v", s.StoreName, err)
		}
	}

	rpc.Register(s)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	http.Serve(listener, nil)
}
