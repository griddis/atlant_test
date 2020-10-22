package limiting

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLimiter(t *testing.T) {
	testCases := []struct {
		in  float64
		out int
	}{
		{1, 5},
		{2, 10},
	}
	ts := httptest.NewServer(Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})))
	defer ts.Close()

	for _, c := range testCases {
		SetLimiter(NewLimiter(c.in, 1))
		var counter int
		for i := 0; i < 10; i++ {
			resp, err := http.Get(ts.URL)
			require.NoError(t, err)
			if resp.StatusCode == 200 {
				counter++
			}
			time.Sleep(time.Millisecond * 500)
		}

		assert.Equal(t, c.out, counter)
	}

}
