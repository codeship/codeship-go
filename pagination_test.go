package codeship

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_paginate(t *testing.T) {
	type args struct {
		path string
		opts ListOptions
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "paginates page",
			args: args{
				path: "/organizations/123/projects",
				opts: ListOptions{
					Page: 1,
				},
			},
			want: "/organizations/123/projects?page=1",
		},
		{
			name: "paginates per_page",
			args: args{
				path: "/organizations/123/projects",
				opts: ListOptions{
					PerPage: 10,
				},
			},
			want: "/organizations/123/projects?per_page=10",
		},
		{
			name: "paginates both page and per_page",
			args: args{
				path: "/organizations/123/projects",
				opts: ListOptions{
					PerPage: 15,
					Page:    5,
				},
			},
			want: "/organizations/123/projects?page=5&per_page=15",
		},
		{
			name: "handles empty options",
			args: args{
				path: "/organizations/123/projects",
			},
			want: "/organizations/123/projects",
		},
		{
			name: "handles negative page",
			args: args{
				path: "/organizations/123/projects",
				opts: ListOptions{
					Page: -1,
				},
			},
			want: "/organizations/123/projects",
		},
		{
			name: "handles negative per_page",
			args: args{
				path: "/organizations/123/projects",
				opts: ListOptions{
					PerPage: -1,
				},
			},
			want: "/organizations/123/projects",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := paginate(tt.args.path, tt.args.opts)

			assert := assert.New(t)
			assert.NoError(err)
			assert.Equal(tt.want, got)
		})
	}
}
