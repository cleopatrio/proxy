package core

import (
	"net/url"
	"testing"
)

func Test_DownstreamURL(t *testing.T) {
	type args struct {
		proxyPath   ProxyPath
		requestHost string
		requestPath string
	}

	exactMatches := []struct {
		name        string
		args        args
		expectedURL string
	}{
		{
			name: "example.com/files",
			args: args{
				proxyPath:   ProxyPath{Path: "/files", PathType: ExactPathType, TLS: false},
				requestHost: "example.com",
				requestPath: "/files",
			},
			expectedURL: "http://example.com/files",
		},
		{
			name: "example.com/files - with tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/files", PathType: ExactPathType, TLS: true},
				requestHost: "example.com",
				requestPath: "/files",
			},
			expectedURL: "https://example.com/files",
		},
		{
			name: "example.com/people - with port number - 1",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: ExactPathType, PortNumber: 3000, TLS: false},
				requestHost: "example.com:4000",
				requestPath: "/people",
			},
			expectedURL: "http://example.com:3000/people",
		},
		{
			name: "example.com/people - with port number - 2",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: ExactPathType, PortNumber: 3000, TLS: false},
				requestHost: "example.com",
				requestPath: "/people",
			},
			expectedURL: "http://example.com:3000/people",
		},
		{
			name: "example.com/people - with port number and tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: ExactPathType, PortNumber: 3000, TLS: true},
				requestHost: "example.com",
				requestPath: "/people",
			},
			expectedURL: "https://example.com:3000/people",
		},
		{
			name: "example.com/animals",
			args: args{
				proxyPath:   ProxyPath{Path: "/animals", PathType: ExactPathType, TLS: false},
				requestHost: "https://example.com",
				requestPath: "/animals",
			},
			expectedURL: "http://example.com/animals",
		},
		{
			name: "example.com/animals - with tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/animals", PathType: ExactPathType, TLS: true},
				requestHost: "https://example.com",
				requestPath: "/animals",
			},
			expectedURL: "https://example.com/animals",
		},
		{
			name: "example.com/cars - mismatch",
			args: args{
				proxyPath:   ProxyPath{Path: "/drivers", PathType: ExactPathType, TLS: false},
				requestHost: "example.com",
				requestPath: "/cars",
			},
			expectedURL: "",
		},
		{
			name: "example.com/cars - mismatch - with tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/drivers", PathType: ExactPathType, TLS: true},
				requestHost: "example.com",
				requestPath: "/cars",
			},
			expectedURL: "",
		},
		{
			name: "example.com/posts - with host port",
			args: args{
				proxyPath:   ProxyPath{Path: "/posts", PathType: ExactPathType, TLS: false},
				requestHost: "example.com:4000",
				requestPath: "/posts",
			},
			expectedURL: "http://example.com/posts",
		},
		{
			name: "example.com/posts - with host port and tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/posts", PathType: ExactPathType, TLS: true},
				requestHost: "example.com:4000",
				requestPath: "/posts",
			},
			expectedURL: "https://example.com/posts",
		},
	}

	prefixMatches := []struct {
		name        string
		args        args
		expectedURL string
	}{
		{
			name: "example.com/files*",
			args: args{
				proxyPath:   ProxyPath{Path: "/files", PathType: PrefixPathType, TLS: false},
				requestHost: "example.com",
				requestPath: "/files/1",
			},
			expectedURL: "http://example.com/files/1",
		},
		{
			name: "example.com/files* - with tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/files", PathType: PrefixPathType, TLS: true},
				requestHost: "example.com",
				requestPath: "/files/1",
			},
			expectedURL: "https://example.com/files/1",
		},
		{
			name: "example.com/people* - with port number - 1",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: PrefixPathType, PortNumber: 3000, TLS: false},
				requestHost: "example.com:4000",
				requestPath: "/people/john",
			},
			expectedURL: "http://example.com:3000/people/john",
		},
		{
			name: "example.com/people* - with port number - 2",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: PrefixPathType, PortNumber: 3000, TLS: false},
				requestHost: "example.com",
				requestPath: "/people/john",
			},
			expectedURL: "http://example.com:3000/people/john",
		},
		{
			name: "example.com/people* - with port number and tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: PrefixPathType, PortNumber: 3000, TLS: true},
				requestHost: "example.com",
				requestPath: "/people/john",
			},
			expectedURL: "https://example.com:3000/people/john",
		},
		{
			name: "example.com/animals*",
			args: args{
				proxyPath:   ProxyPath{Path: "/animals", PathType: PrefixPathType, TLS: false},
				requestHost: "https://example.com",
				requestPath: "/animals/birds",
			},
			expectedURL: "http://example.com/animals/birds",
		},
		{
			name: "example.com/animals* - with tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/animals", PathType: PrefixPathType, TLS: true},
				requestHost: "https://example.com",
				requestPath: "/animals/birds",
			},
			expectedURL: "https://example.com/animals/birds",
		},
		{
			name: "example.com/cars - mismatch",
			args: args{
				proxyPath:   ProxyPath{Path: "/cars/v1", PathType: PrefixPathType, TLS: false},
				requestHost: "example.com",
				requestPath: "/cars",
			},
			expectedURL: "",
		},
		{
			name: "example.com/cars - mismatch - with tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/cars/v1", PathType: PrefixPathType, TLS: true},
				requestHost: "example.com",
				requestPath: "/cars",
			},
			expectedURL: "",
		},
		{
			name: "example.com/posts* - with host port",
			args: args{
				proxyPath:   ProxyPath{Path: "/posts", PathType: PrefixPathType, TLS: false},
				requestHost: "example.com:4000",
				requestPath: "/posts/5",
			},
			expectedURL: "http://example.com/posts/5",
		},
		{
			name: "example.com/posts* - with host port and tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/posts", PathType: PrefixPathType, TLS: true},
				requestHost: "example.com:4000",
				requestPath: "/posts/5",
			},
			expectedURL: "https://example.com/posts/5",
		},
	}

	for _, tt := range append(exactMatches, prefixMatches...) {
		t.Run(tt.name, func(t *testing.T) {
			url := tt.args.proxyPath.DownstreamURL(tt.args.requestHost, tt.args.requestPath)
			if url != tt.expectedURL {
				t.Errorf(`expected %s but got %s`, tt.expectedURL, url)
			}
		})
	}
}

func Test_RequestURL(t *testing.T) {
	type args struct {
		proxyPath   ProxyPath
		requestHost string
		requestPath string
	}

	exactMatches := []struct {
		name        string
		args        args
		expectedURL *url.URL
	}{
		{
			name: "example.com/files",
			args: args{
				proxyPath:   ProxyPath{Path: "/files", PathType: ExactPathType, TLS: false},
				requestHost: "example.com",
				requestPath: "/files",
			},
			expectedURL: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "files",
			},
		},
		{
			name: "example.com/files - with tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/files", PathType: ExactPathType, TLS: true},
				requestHost: "example.com",
				requestPath: "/files",
			},
			expectedURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "files",
			},
		},
		{
			name: "example.com/people - with port number - 1",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: ExactPathType, PortNumber: 3000, TLS: false},
				requestHost: "example.com:4000",
				requestPath: "/people",
			},
			expectedURL: &url.URL{
				Scheme: "http",
				Host:   "example.com:3000",
				Path:   "people",
			},
		},
		{
			name: "example.com/people - with port number - 2",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: ExactPathType, PortNumber: 3000, TLS: false},
				requestHost: "example.com",
				requestPath: "/people",
			},
			expectedURL: &url.URL{
				Scheme: "http",
				Host:   "example.com:3000",
				Path:   "people",
			},
		},
		{
			name: "example.com/people - with port number and tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: ExactPathType, PortNumber: 3000, TLS: true},
				requestHost: "example.com",
				requestPath: "/people",
			},
			expectedURL: &url.URL{
				Scheme: "https",
				Host:   "example.com:3000",
				Path:   "people",
			},
		},
		{
			name: "example.com/animals",
			args: args{
				proxyPath:   ProxyPath{Path: "/animals", PathType: ExactPathType, TLS: false},
				requestHost: "https://example.com",
				requestPath: "/animals",
			},
			expectedURL: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "animals",
			},
		},
		{
			name: "example.com/animals - with tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/animals", PathType: ExactPathType, TLS: true},
				requestHost: "https://example.com",
				requestPath: "/animals",
			},
			expectedURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "animals",
			},
		},
		{
			name: "example.com/cars - mismatch",
			args: args{
				proxyPath:   ProxyPath{Path: "/drivers", PathType: ExactPathType, TLS: false},
				requestHost: "example.com",
				requestPath: "/cars",
			},
			expectedURL: nil,
		},
		{
			name: "example.com/cars - mismatch - with tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/drivers", PathType: ExactPathType, TLS: true},
				requestHost: "example.com",
				requestPath: "/cars",
			},
			expectedURL: nil,
		},
		{
			name: "example.com/posts - with host port",
			args: args{
				proxyPath:   ProxyPath{Path: "/posts", PathType: ExactPathType, TLS: false},
				requestHost: "example.com:4000",
				requestPath: "/posts",
			},
			expectedURL: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "posts",
			},
		},
		{
			name: "example.com/posts - with host port and tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/posts", PathType: ExactPathType, TLS: true},
				requestHost: "example.com:4000",
				requestPath: "/posts",
			},
			expectedURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "posts",
			},
		},
	}

	prefixMatches := []struct {
		name        string
		args        args
		expectedURL *url.URL
	}{
		{
			name: "example.com/files*",
			args: args{
				proxyPath:   ProxyPath{Path: "/files", PathType: PrefixPathType, TLS: false},
				requestHost: "example.com",
				requestPath: "/files/1",
			},
			expectedURL: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/files/1",
			},
		},
		{
			name: "example.com/files* - with tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/files", PathType: PrefixPathType, TLS: true},
				requestHost: "example.com",
				requestPath: "/files/1",
			},
			expectedURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/files/1",
			},
		},
		{
			name: "example.com/people* - with port number",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: PrefixPathType, PortNumber: 3000, TLS: false},
				requestHost: "example.com",
				requestPath: "/people/john",
			},
			expectedURL: &url.URL{
				Scheme: "http",
				Host:   "example.com:3000",
				Path:   "/people/john",
			},
		},
		{
			name: "example.com/people* - with port number and tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: PrefixPathType, PortNumber: 3000, TLS: true},
				requestHost: "example.com",
				requestPath: "/people/john",
			},
			expectedURL: &url.URL{
				Scheme: "https",
				Host:   "example.com:3000",
				Path:   "/people/john",
			},
		},
		{
			name: "example.com/animals*",
			args: args{
				proxyPath:   ProxyPath{Path: "/animals", PathType: PrefixPathType, TLS: false},
				requestHost: "https://example.com",
				requestPath: "/animals/birds",
			},
			expectedURL: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/animals/birds",
			},
		},
		{
			name: "example.com/animals* - with tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/animals", PathType: PrefixPathType, TLS: true},
				requestHost: "https://example.com",
				requestPath: "/animals/birds",
			},
			expectedURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/animals/birds",
			},
		},
		{
			name: "example.com/cars - mismatch",
			args: args{
				proxyPath:   ProxyPath{Path: "/cars/v1", PathType: PrefixPathType, TLS: false},
				requestHost: "example.com",
				requestPath: "/cars",
			},
			expectedURL: nil,
		},
		{
			name: "example.com/cars - mismatch - with tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/cars/v1", PathType: PrefixPathType, TLS: true},
				requestHost: "example.com",
				requestPath: "/cars",
			},
			expectedURL: nil,
		},
		{
			name: "example.com/posts* - with host port",
			args: args{
				proxyPath:   ProxyPath{Path: "/posts", PathType: PrefixPathType, TLS: false},
				requestHost: "example.com:4000",
				requestPath: "/posts/5",
			},
			expectedURL: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/posts/5",
			},
		},
		{
			name: "example.com/posts* - with host port and tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/posts", PathType: PrefixPathType, TLS: true},
				requestHost: "example.com:4000",
				requestPath: "/posts/5",
			},
			expectedURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/posts/5",
			},
		},
	}

	for _, tt := range append(exactMatches, prefixMatches...) {
		t.Run(tt.name, func(t *testing.T) {
			url := tt.args.proxyPath.RequestURL(tt.args.requestHost, tt.args.requestPath)
			if url != nil && tt.expectedURL != nil {
				if url.String() != tt.expectedURL.String() {
					t.Errorf(`expected %s but got %s`, tt.expectedURL, url)
				}
			}

			// t.Errorf(`expected non nil url`)
		})
	}
}
