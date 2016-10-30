package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"

	pb "github.com/imgcompress/lossycompress"

	"google.golang.org/grpc"
)

const port = ":3030"

func main() {
	Send()
}

func Send() {
	hostname := flag.String("hostname", "localhost", "Hostname for the connection")
	flag.Parse()

	img, err := ioutil.ReadFile("face.jpg")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(*hostname+port, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewImgClient(conn)
	req := &pb.Request{string(img), 50, "face.jpg"}

	res, err := client.Compress(context.Background(), req)
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}
