# semaphore

ゴルーチンの同時実行数を制御するパッケージ。

読み方はセマフォ。

- https://godoc.org/golang.org/x/sync/semaphore

## そもそも semaphore（セマフォ）とは？

- 同時に実行されているプログラム間でのリソースの排他制御や同期を行う仕組みのこと。
- 現在利用可能なリソースの数のこと。

以下の記事でなんとなく仕組みは理解できる。とりあえず「同時実行を制御するための仕組み・手法」ぐらいな認識をしておけば問題ないと思う。

- [セマフォ (semaphore)とは](https://wa3.i-3-i.info/word13357.html)
- [【問題 2】セマフォと排他制御](https://monoist.atmarkit.co.jp/mn/articles/1009/24/news100.html)

仕組みをざっくり理解できたら、コードを書いて動きを確認した方が理解が深まると思う。

## 利用例

```shell
go get golang.org/x/sync/semaphore
```

### ゴルーチンの同時実行数を１つにする

以下は３つのゴルーチンが並列で実行されるコード。

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func longProcess(ctx context.Context) {
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
```

semaphore を利用してゴルーチンを１つずつ実行させる場合、以下のようになる。

```go
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
```

### セマフォが取得できなかった場合（待機状態になる場合）、待機せずにゴルーチンを終了する

```go
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
	// s.TryAcquire(1) でセマフォを取得する。
	// 成功すれば true を返し、失敗（他のゴルーチンがセマフォを取得しており、セマフォを取得できないなどが原因で失敗）すれば false を返す。
	isAcquire := s.TryAcquire(1)
	if !isAcquire {
		fmt.Println("Could not get lock")
		return // ゴルチーンを終了
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
```
