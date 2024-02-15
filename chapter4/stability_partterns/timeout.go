package stability_partterns

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// SlowFunc - context.Context 를 매개변수로 받지 않아 timeout 설정이 불가능
type SlowFunc func(string) (string, error)

// SlowFuncWithCtx - SlowFunc 시그니처에 context 를 포함한 타입
type SlowFuncWithCtx func(context.Context, string) (string, error)

// MySlowFunc - [0,5) 구간의 응답시간이 걸리는 함수
func MySlowFunc(arg string) (string, error) {
	respTime := time.Duration(rand.Intn(5)) * time.Second
	fmt.Println("my slow func takes " + respTime.String() + " second(s).")
	time.Sleep(respTime + time.Millisecond*50) // 타임아웃과 동일한 응답시간에 대해 실패시키기위해 보정

	return fmt.Sprintf("%s success", arg), nil
}

// withCtx - SlowFunc 를 SlowFuncWithCtx 반환
func withCtx(f SlowFunc) SlowFuncWithCtx {
	return func(ctx context.Context, arg string) (string, error) {
		chres := make(chan string)
		cherr := make(chan error)

		go func() {
			res, err := f(arg)
			chres <- res
			cherr <- err
		}()

		select {
		case res := <-chres:
			return res, <-cherr
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

func useCase(timeToWait time.Duration) (string, error) {
	ctxT, cancel := context.WithTimeout(context.Background(), timeToWait)
	defer cancel()

	slowFuncWithCtx := withCtx(MySlowFunc)
	res, err := slowFuncWithCtx(ctxT, "arg")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("res: %s\n", res), nil
}
