package config

import "testing"

func Test_DownstreamURL_ExactMatch(t *testing.T) {
	type args struct {
		proxyPath   ProxyPath
		requestHost string
		requestPath string
	}

	tests := []struct {
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
			name: "example.com/people - with port number",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: ExactPathType, Port: 3000, TLS: false},
				requestHost: "example.com",
				requestPath: "/people",
			},
			expectedURL: "http://example.com:3000/people",
		},
		{
			name: "example.com/people - with port number and tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: ExactPathType, Port: 3000, TLS: true},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := tt.args.proxyPath.DownstreamURL(tt.args.requestHost, tt.args.requestPath)
			if url != tt.expectedURL {
				t.Errorf(`expected %s but got %s`, tt.expectedURL, url)
			}
		})
	}
}

func Test_DownstreamURL_PrefixMatch(t *testing.T) {
	type args struct {
		proxyPath   ProxyPath
		requestHost string
		requestPath string
	}

	tests := []struct {
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
			name: "example.com/people* - with port number",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: PrefixPathType, Port: 3000, TLS: false},
				requestHost: "example.com",
				requestPath: "/people/john",
			},
			expectedURL: "http://example.com:3000/people/john",
		},
		{
			name: "example.com/people* - with port number and tls",
			args: args{
				proxyPath:   ProxyPath{Path: "/people", PathType: PrefixPathType, Port: 3000, TLS: true},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := tt.args.proxyPath.DownstreamURL(tt.args.requestHost, tt.args.requestPath)
			if url != tt.expectedURL {
				t.Errorf(`expected %s but got %s`, tt.expectedURL, url)
			}
		})
	}
}
