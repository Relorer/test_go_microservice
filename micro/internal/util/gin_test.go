package util

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGinHandleError(t *testing.T) {

	type args struct {
		err    error
		status int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "With error",
			args: args{
				err:    errors.New("some error"),
				status: http.StatusNotFound,
			},
			want: true,
		},
		{
			name: "No error",
			args: args{
				err:    nil,
				status: http.StatusOK,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			got := GinHandleError(c, tt.args.err, tt.args.status)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.args.status, w.Code)
		})
	}
}
