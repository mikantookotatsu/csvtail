package csvf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// CsvfInf はCSVファイルの情報を保持する構造体
type CsvfInf struct {
	FileName  string   // 監視ファイル名
	Columns   []int    // 表示するカラム番号
	Delimiter string   // 区切り文字
	Seconds   int      // 監視間隔(秒)
	FilePtr   *os.File // 対象ファイル
}

// ファイル操作のエラーチェック
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// ファイル存在チェック
func (c *CsvfInf) FileExists() bool {
	_, err := os.Stat(c.FileName)
	if err == nil {
		return true // ファイル有り
	}
	// ファイルが存在しないエラー(ErrNotExist)でない場合
	if !errors.Is(err, os.ErrNotExist) {
		return true // ファイル有り
	}

	// ファイルなし
	return false
}

// ファイルオープン
func (c *CsvfInf) FileOpen() {
	var err error
	c.FilePtr, err = os.Open(c.FileName)
	check(err)
}

// 末尾にSeek
func (c *CsvfInf) SeekEnd() {
	_, err := c.FilePtr.Seek(0, io.SeekEnd) // 末尾に移動
	check(err)
}

// ファイル監視
func (c *CsvfInf) FileWatch() {

	reader := bufio.NewReader(c.FilePtr)

	ticker := time.NewTicker(time.Duration(c.Seconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			line, err := reader.ReadString('\n')
			if err != nil && err.Error() != "EOF" {
				log.Fatal(err)
			}

			// 新しい行があれば表示
			if line != "" {
				fmt.Println(strings.TrimSpace(line))
			}
		}
	}
}

// ファイルクローズ
func (c *CsvfInf) FileClose() {
	// オープンしているファイルがある場合のみクローズ
	if c.FilePtr != nil {
		err := c.FilePtr.Close()
		check(err)
	}
}
