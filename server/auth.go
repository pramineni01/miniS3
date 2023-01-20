package main

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/metadata"
)

type AliasInfo struct {
	AccessKey string   `json:"access_key"`
	SecretKey string   `json:"secret_key"`
	Buckets   []string `json:"buckets"`
}

type Authenticator interface {
	Authenticate(ctx context.Context) error
	UpdateAuth(aliasName string, accessKey string, secretKey string)
	GetBuckets(aliasName string) ([]string, error)
	UpdateBuckets(aliasName string, buckets []string) error
	GetAliasName(ctx context.Context) string
}

type authenticator struct {
	AuthInfo map[string]AliasInfo
}

func NewAuthenticator() Authenticator {
	return &authenticator{
		AuthInfo: map[string]AliasInfo{},
	}
}

func (a *authenticator) Authenticate(ctx context.Context) error {
	var aliasName, accessKey, secretKey string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		aliasName = md["alias-name"][0]
		accessKey = md["access-key"][0]
		secretKey = md["secret-key"][0]
	}
	fmt.Printf("aliasName: %s, accessKey: %s, secretKey: %s\n", aliasName, accessKey, secretKey)
	if v, found := a.AuthInfo[aliasName]; !found {
		return errors.New("Alias not found")
	} else if v.AccessKey != accessKey || v.SecretKey != secretKey {
		return errors.New("Authentication failed.")
	}
	return nil
}

func (a *authenticator) UpdateAuth(aliasName string, accessKey string, secretKey string) {
	buckets := []string{}
	if v, found := a.AuthInfo[aliasName]; found {
		buckets = v.Buckets
	}
	a.AuthInfo[aliasName] = AliasInfo{AccessKey: accessKey, SecretKey: secretKey, Buckets: buckets}
	fmt.Printf("Updated Auth data: %v", a.AuthInfo)
}

func (a *authenticator) GetBuckets(aliasName string) ([]string, error) {
	if v, found := a.AuthInfo[aliasName]; !found {
		return nil, errors.New("Alias not found")
	} else {
		return v.Buckets, nil
	}
}

func (a *authenticator) UpdateBuckets(aliasName string, buckets []string) error {
	if v, found := a.AuthInfo[aliasName]; !found {
		return errors.New("Alias not found")
	} else {
		a.AuthInfo[aliasName] = AliasInfo{AccessKey: v.AccessKey, SecretKey: v.SecretKey, Buckets: buckets}
		return nil
	}
}

func (a authenticator) GetAliasName(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return md["alias-name"][0]
	}
	return ""
}
