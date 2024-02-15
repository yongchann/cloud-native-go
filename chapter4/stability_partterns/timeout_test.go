package stability_partterns

import (
	"testing"
	"time"
)

func Test_useCase(t *testing.T) {
	tests := []struct {
		name       string
		timeToWait time.Duration
	}{
		{
			name:       "context_timeout=1초",
			timeToWait: time.Second,
		},
		{
			name:       "context_timeout=2초",
			timeToWait: 2 * time.Second,
		},
		{
			name:       "context_timeout=3초",
			timeToWait: 3 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := useCase(tt.timeToWait)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

//=== RUN   Test_useCase/context_timeout=1초
//this func takes 1s second(s).
//--- PASS: Test_useCase/context_timeout=1초 (1.00s)

//=== RUN   Test_useCase/context_timeout=2초
//this func takes 2s second(s).
//timeout_test.go:30: context deadline exceeded
//--- FAIL: Test_useCase/context_timeout=2초 (2.00s)

//=== RUN   Test_useCase/context_timeout=3초
//this func takes 2s second(s).
//--- PASS: Test_useCase/context_timeout=3초 (2.00s)
