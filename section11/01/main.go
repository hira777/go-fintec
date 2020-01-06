package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"time"
)

// 引数に渡した数が最大カウント（同時アクセス数）として指定されたセマフォを作成する。
// 今回は NewWeighted(1) なので同時に実行できるゴルーチンは１つ。
var s *semaphore.Weighted = semaphore.NewWeighted(1)

func longProcess(ctx context.Context) {
	// s.Acquire(ctx, 1) でセマフォを取得する。セマフォを取得することで、セマフォのカウントが１つ減る。
	// 今回のセマフォの最大カウントは１なので、セマフォのカウントは０になり、他のゴルチーンはブロックされて待機状態になる。
	if err := s.Acquire(ctx, 1); err != nil {
		fmt.Println(err)
		return
	}
	// セマフォをリリースする。セマフォのカウントは１に戻り、待機状態のゴルーチンが実行される。
	defer s.Release(1)
	fmt.Println("Wait...")
	time.Sleep(1 * time.Second)
	fmt.Println("Done")
}

func main() {
	ctx := context.TODO()
	go longProcess(ctx)
	go longProcess(ctx)
	go longProcess(ctx)
	time.Sleep(5 * time.Second)
}
