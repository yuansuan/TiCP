package util

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/metadata"
)

// RemoteHost RemoteHost
const RemoteHost = "X-REMOTE-EDGEPROXY"

// AppendToIncomingContext append `k`, `v` to incoming context
func AppendToIncomingContext(ctx context.Context, k, v string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	md.Set(k, v)
	return metadata.NewIncomingContext(ctx, md)
}

// SetInMetadata SetInMetadata
func SetInMetadata(ctx context.Context, key, value string) context.Context {
	return AppendToIncomingContext(ctx, key, value)
}

// GetInMetadata returns metadata
func GetInMetadata(ctx context.Context, key string) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if len(md[key]) == 0 {
			return "", fmt.Errorf("metadata[%q] is empty", key)
		}
		return md[key][0], nil
	}
	return "", fmt.Errorf("metadata does not exists")
}

// GetOutMetadata return metadata
func GetOutMetadata(ctx context.Context, key string) (string, error) {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		if len(md[key]) == 0 {
			return "", fmt.Errorf("metadata[%q] is empty", key)
		}
		return md[key][0], nil
	}
	return "", fmt.Errorf("metadata does not exists")
}

// AppendToOutgoingContext append `k`, `v` to incoming context
func AppendToOutgoingContext(ctx context.Context, k, v string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	md.Set(k, v)
	return metadata.NewOutgoingContext(ctx, md)
}

// AppendToOutgoingContextByHeader AppendToOutgoingContextByHeader
func AppendToOutgoingContextByHeader(ctx context.Context, header http.Header) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	for k, v := range header {
		md.Append(k, v...)
	}
	return metadata.NewOutgoingContext(ctx, md)
}

// HeaderFromIncomingContext HeaderFromIncomingContext
func HeaderFromIncomingContext(ctx context.Context) http.Header {
	header := http.Header{}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for k, vs := range md {
			for _, v := range vs {
				header.Set(k, v)
			}
		}
	}
	return header
}
