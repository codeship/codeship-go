package codeship

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_paginate(t *testing.T) {
	type args struct {
		path string
		opts []PaginationOption
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
				opts: []PaginationOption{
					Page(1),
				},
			},
			want: "/organizations/123/projects?page=1",
		},
		{
			name: "paginates per_page",
			args: args{
				path: "/organizations/123/projects",
				opts: []PaginationOption{
					PerPage(10),
				},
			},
			want: "/organizations/123/projects?per_page=10",
		},
		{
			name: "paginates both page and per_page",
			args: args{
				path: "/organizations/123/projects",
				opts: []PaginationOption{
					PerPage(15),
					Page(5),
				},
			},
			want: "/organizations/123/projects?page=5&per_page=15",
		},
		{
			name: "handles multiple calls to page",
			args: args{
				path: "/organizations/123/projects",
				opts: []PaginationOption{
					Page(1),
					Page(5),
				},
			},
			want: "/organizations/123/projects?page=5",
		},
		{
			name: "handles multiple calls to per_page",
			args: args{
				path: "/organizations/123/projects",
				opts: []PaginationOption{
					PerPage(10),
					PerPage(15),
				},
			},
			want: "/organizations/123/projects?per_page=15",
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
				opts: []PaginationOption{
					Page(-1),
				},
			},
			want: "/organizations/123/projects",
		},
		{
			name: "handles negative per_page",
			args: args{
				path: "/organizations/123/projects",
				opts: []PaginationOption{
					PerPage(-1),
				},
			},
			want: "/organizations/123/projects",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := paginate(tt.args.path, tt.args.opts...)

			assert := assert.New(t)
			assert.NoError(err)
			assert.Equal(tt.want, got)
		})
	}
}

func TestLinks_NextPage(t *testing.T) {
	type fields struct {
		Next     string
		Previous string
		Last     string
		First    string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "returns next page if next page",
			fields: fields{
				Next: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=2",
			},
			want: 2,
		},
		{
			name: "returns 0 if no next page",
			fields: fields{
				Previous: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=2",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Links{
				Next:     tt.fields.Next,
				Previous: tt.fields.Previous,
				Last:     tt.fields.Last,
				First:    tt.fields.First,
			}
			got, err := l.NextPage()

			assert := assert.New(t)
			assert.NoError(err)
			assert.Equal(tt.want, got)
		})
	}
}

func TestLinks_PreviousPage(t *testing.T) {
	type fields struct {
		Next     string
		Previous string
		Last     string
		First    string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "returns previous page if previous page",
			fields: fields{
				Previous: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=1",
			},
			want: 1,
		},
		{
			name: "returns 0 if no previous page",
			fields: fields{
				Next: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=2",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Links{
				Next:     tt.fields.Next,
				Previous: tt.fields.Previous,
				Last:     tt.fields.Last,
				First:    tt.fields.First,
			}
			got, err := l.PreviousPage()

			assert := assert.New(t)
			assert.NoError(err)
			assert.Equal(tt.want, got)
		})
	}
}

func TestLinks_CurrentPage(t *testing.T) {
	type fields struct {
		Next     string
		Previous string
		Last     string
		First    string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "returns 1 if on first page",
			fields: fields{
				Next: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=2",
			},
			want: 1,
		},
		{
			name: "returns current page if not on first page",
			fields: fields{
				Next:     "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=3",
				Previous: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=1",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Links{
				Next:     tt.fields.Next,
				Previous: tt.fields.Previous,
				Last:     tt.fields.Last,
				First:    tt.fields.First,
			}
			got, err := l.CurrentPage()

			assert := assert.New(t)
			assert.NoError(err)
			assert.Equal(tt.want, got)
		})
	}
}

func TestLinks_LastPage(t *testing.T) {
	type fields struct {
		Next     string
		Previous string
		Last     string
		First    string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "returns last page if last page",
			fields: fields{
				Last: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=10",
			},
			want: 10,
		},
		{
			name: "returns current page if on last page",
			fields: fields{
				Previous: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=10",
			},
			want: 11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Links{
				Next:     tt.fields.Next,
				Previous: tt.fields.Previous,
				Last:     tt.fields.Last,
				First:    tt.fields.First,
			}
			got, err := l.LastPage()

			assert := assert.New(t)
			assert.NoError(err)
			assert.Equal(tt.want, got)
		})
	}
}

func TestLinks_IsLastPage(t *testing.T) {
	type fields struct {
		Next     string
		Previous string
		Last     string
		First    string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "returns true if on last page",
			fields: fields{
				Previous: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=10",
			},
			want: true,
		},
		{
			name: "returns false if not on last page",
			fields: fields{
				Previous: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=10",
				Next:     "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=12",
				Last:     "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=12",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Links{
				Next:     tt.fields.Next,
				Previous: tt.fields.Previous,
				Last:     tt.fields.Last,
				First:    tt.fields.First,
			}
			got := l.IsLastPage()

			assert := assert.New(t)
			assert.Equal(tt.want, got)
		})
	}
}
