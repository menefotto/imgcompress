package main

import (
	"context"
	"io/ioutil"
	"log"

	pb "github.com/imgcompress/lossycompress"

	"google.golang.org/grpc"
)

func main() {
	Send()
}

func Send() {
	img, err := ioutil.ReadFile("face.jpg")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial("localhost:3030", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewImgClient(conn)
	req := &pb.Request{string(img), 50, "image-name"}

	res, err := client.Compress(context.Background(), req)
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}
