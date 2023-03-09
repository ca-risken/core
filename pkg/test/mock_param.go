package test

import "github.com/stretchr/testify/mock"

// RepeatMockAnything は、mock.Anythingを指定した件数設定した配列を返します
// テストで使用するmockを自動生成したものに置き換える際の暫定対応として使う想定
func RepeatMockAnything(n int) []interface{} {
	var ret []interface{}
	for i := 0; i < n; i++ {
		ret = append(ret, mock.Anything)
	}
	return ret
}

func Int64(n int64) *int64 {
	return &n
}

func Int(n int) *int {
	return &n
}
