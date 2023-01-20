package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/pramimeni01/minis3/gen/minis3"
)

type BucketServiceImpl struct {
	rootFolder string
	auth       Authenticator
	pb.UnimplementedBucketServiceServer
}

func NewBucketService(storage string) pb.BucketServiceServer {
	fmt.Println("NewBucketService(): storage: ", storage)
	return &BucketServiceImpl{
		rootFolder: storage,
		auth:       NewAuthenticator(),
	}
}

func main() {

	if len(os.Args) == 1 {
		log.Fatal("Missing storage path.\n Usage:\n\t mini_s3_server <storage_location>")
	}

	s := grpc.NewServer()
	pb.RegisterBucketServiceServer(s, NewBucketService(os.Args[1]))

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Unable to open tcp connection. Error: ", err)
	}

	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			aliasName := request.Header.Get("alias-name")
			accessKey := request.Header.Get("access-key")
			secretKey := request.Header.Get("secret-key")
			// send all the headers received from the client
			md := metadata.Pairs("alias-name", aliasName, "access-key", accessKey, "secret-key", secretKey)
			return md
		}),
	)
	err = pb.RegisterBucketServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())

}

// BucketService Implementation
func (bs *BucketServiceImpl) SetAlias(ctx context.Context, req *pb.SetAliasRequest) (*pb.Response, error) {
	fmt.Println("SetAlias(): begin")
	if (strings.TrimSpace(req.GetAliasName()) == "") ||
		(strings.TrimSpace(req.GetAccessKey()) == "") ||
		(strings.TrimSpace(req.GetSecretKey()) == "") {
		fmt.Printf("SetAlias(): Input: %v", req)
		return &pb.Response{Status: "Failed due to invalid credenials"}, errors.New("Invalid credentials")
	}
	fmt.Printf("SetAlias(): Input outside: %v\n", req)
	bs.auth.UpdateAuth(req.AliasName, req.AccessKey, req.SecretKey)
	return &pb.Response{Status: "Success"}, nil
}

// BucketService Implementation
func (bs *BucketServiceImpl) CreateBucket(ctx context.Context, req *pb.Request) (*pb.Response, error) {

	if err := bs.auth.Authenticate(ctx); err != nil {
		return &pb.Response{Status: "Authentication failed."}, errors.New("Authentication failed.")
	}

	err := createHome(bs.rootFolder, bs.auth.GetAliasName(ctx), req.GetBucketName())
	if err != nil {
		msg := fmt.Sprintln("Failed to create bucket: ", err)
		return &pb.Response{Status: msg}, err
	}
	return &pb.Response{Status: "Success"}, nil
}

func (bs *BucketServiceImpl) UploadToBucket(ctx context.Context, req *pb.UploadToBucketRequest) (*pb.Response, error) {

	if err := bs.auth.Authenticate(ctx); err != nil {
		return &pb.Response{Status: "Authentication failed."}, errors.New("Authentication failed.")
	}

	path := getPath(bs.rootFolder, bs.auth.GetAliasName(ctx), req.GetBucketName())

	// copy file into folder by bucket name
	err := os.WriteFile(fmt.Sprintf("%s/%s", path, req.GetFileName()), []byte(req.GetData()), 0660)
	if err != nil {
		return &pb.Response{Status: "Internal server error"}, err
	}

	return &pb.Response{Status: "success"}, nil
}

func (bs *BucketServiceImpl) DeleteBucket(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	if err := bs.auth.Authenticate(ctx); err != nil {
		return &pb.Response{Status: "Authentication failed."}, errors.New("Authentication failed.")
	}

	path := getPath(bs.rootFolder, bs.auth.GetAliasName(ctx), req.GetBucketName())
	if !pathExists(path) {
		err := fmt.Sprintf("%s does not exist", path)
		return &pb.Response{Status: err}, errors.New(err)
	}
	if err := os.RemoveAll(path); err != nil {
		return &pb.Response{Status: err.Error()}, err
	}
	return &pb.Response{Status: fmt.Sprintf("%s delete successful.", req.GetBucketName())}, nil
}

func getPath(rootFolder string, aliasName string, bucket string) string {
	fmt.Printf("Root: %s", rootFolder)
	return fmt.Sprintf("%s/home/%s/%s", strings.TrimSpace(rootFolder), strings.TrimSpace(aliasName), strings.TrimSpace(bucket))
}

func createHome(rootFolder string, aliasName string, bucketName string) error {
	if len(aliasName) <= 0 {
		return errors.New("Invalid user name")
	}

	err := os.MkdirAll(getPath(rootFolder, aliasName, bucketName), 0750)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
