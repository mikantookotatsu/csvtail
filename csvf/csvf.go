package csvf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// CsvfInf はCSVファイルの情報を保持する構造体
type CsvfInf struct {
	FileName  string   // 監視ファイル名
	Columns   []int    // 表示するカラム番号
	Delimiter string   // 区切り文字
	LineBreak string   // 改行コード
	Seconds   int      // 監視間隔(秒)
	FilePtr   *os.File // 対象ファイル
}

const (
	Delimiter = iota
	Lf
	Unfix
)

// ファイル存在チェック
func (c *CsvfInf) FileExists() (bool, error) {
	_, err := os.Stat(c.FileName)
	if err == nil {
		return true, nil // ファイル有り
	}
	// ファイルが存在しないエラーの場合
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // ファイル無し
	}

	// その他のエラー
	return false, err
}

// ファイルオープン
func (c *CsvfInf) FileOpen() error {
	var err error
	c.FilePtr, err = os.Open(c.FileName)
	return err
}

// 末尾にSeek
func (c *CsvfInf) SeekEnd() error {
	_, err := c.FilePtr.Seek(0, io.SeekEnd) // 末尾に移動
	return err
}

// ファイル全読み込み
func readAllBytes(reader *bufio.Reader) ([]byte, error) {
	var bytes []byte
	for {
		line, err := reader.ReadBytes('\n')
		bytes = append(bytes, line...)

		if err == io.EOF {
			return bytes, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

// ファイル監視
func (c *CsvfInf) FileWatch() error {
	// 監視関数の決定
	var watcher func(line []byte)
	if len(c.Columns) == 0 {
		// 全カラムを監視
		watcher = c.processAllColumnWatch()
	} else {
		// 指定カラムを監視
		watcher = c.processColumnWatch()
	}

	// 監視間隔を決定(0以下は1msとする)
	var d time.Duration
	if c.Seconds <= 0 {
		d = time.Duration(1) * time.Millisecond
	} else {
		d = time.Duration(c.Seconds) * time.Second
	}

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	// 対象ファイル読み込み
	reader := bufio.NewReader(c.FilePtr)

	// 無限ループ (Ctrl+Cで終了)
	for {
		select {
		case <-ticker.C:
			bytes, err := readAllBytes(reader)
			if err != nil {
				return err
			}
			// 新しい出力をチェック
			var n int = len(bytes)
			if n == 0 {
				continue
			}

			// 監視処理
			watcher(bytes)
		}
	}
}

// 監視処理：全カラム出力
func (c *CsvfInf) processAllColumnWatch() func(line []byte) {
	return func(line []byte) {
		fmt.Print(string(line))
	}
}

// 監視処理：指定カラムのみ出力制御
func (c *CsvfInf) processColumnWatch() func(line []byte) {
	var nowColumnNo int = 0        // 現在の列番号
	var columnCnt int = 0          // 列番号
	var unfixColumnStr string = "" // 未確定文字列
	return func(line []byte) {
		var totalByte int = 0
		var sb strings.Builder
		// カラムの抽出
		for {
			// 読み出しデータ解析済
			if totalByte >= len(line) {
				return
			}

			// 1カラム抽出する
			// 区切り文字または改行およびEOFまで抽出した文字列を返却
			// ex) "aaa,bbb\n" -> "aaa", 4, Delimiter -> "bbb", 4, Lf
			//     "aaa,bbb"   -> "aaa", 4, Delimiter -> "bbb", 4, Unfix
			extractStr, readByte, lastCharType := c.getColumn(line[totalByte:])
			totalByte += readByte

			switch {
			case lastCharType == Delimiter: // 区切り文字
				if columnCnt < len(c.Columns) && nowColumnNo == c.Columns[columnCnt] {
					sb.WriteString(unfixColumnStr)
					sb.WriteString(extractStr)
					sb.WriteString(c.Delimiter)
					fmt.Print(sb.String())
					sb.Reset()
					columnCnt++
				}
				unfixColumnStr = ""
				nowColumnNo++

			case lastCharType == Lf: // 改行コード
				if columnCnt < len(c.Columns) && nowColumnNo == c.Columns[columnCnt] {
					sb.WriteString(unfixColumnStr)
					sb.WriteString(extractStr)
				}
				// 改行コードの場合は次の行に移るため初期化
				sb.WriteString(c.LineBreak)
				fmt.Print(sb.String())
				sb.Reset()
				unfixColumnStr = ""
				nowColumnNo = 0
				columnCnt = 0

			default: // 未確定
				unfixColumnStr += extractStr
			}
		}
	}
}

// 1カラムを取り出し
func (c *CsvfInf) getColumn(bytes []byte) (extractStr string, readByte int, lastCharType int) {
	var b byte = 0
	lastCharType = Unfix
	for readByte, b = range bytes {
		// 区切り文字
		if string(b) == c.Delimiter {
			lastCharType = Delimiter
			extractStr = string(bytes[:readByte])
			readByte++ // 区切り文字も含める
			break
		}
		// 改行コード
		if string(b) == "\n" {
			lastCharType = Lf
			extractStr = string(bytes[:readByte])
			extractStr = strings.ReplaceAll(extractStr, "\r", "") // \rが残る場合削除
			readByte++                                            // 改行も含める
			break
		}
	}
	// もし未確定の場合は、EOFまでの文字列
	if lastCharType == Unfix {
		extractStr = string(bytes)
		readByte = len(bytes)
	}

	return
}

// ファイルクローズ
func (c *CsvfInf) FileClose() error {
	// オープンしているファイルがある場合のみクローズ
	if c.FilePtr != nil {
		err := c.FilePtr.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
