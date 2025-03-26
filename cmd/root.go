/*
Copyright (c) 2025 mikantookotatsu
*/
package cmd

import (
	"fmt"
	"os"
	"runtime"
	"sort"

	"github.com/mikantookotatsu/csvtail/csvf"
	"github.com/spf13/cobra"
)

var (
	seconds   int    // -s オプションで指定された秒数
	columns   []int  // -c オプションで指定されたカラム
	delimiter string // -d オプションで指定された区切り文字
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "csvtail [ファイル名]",
	Short: "csvtail は CSVファイルを監視するツールです.",
	Long: `csvtail は CSVファイルを監視するツールです(tail -fのように).所定のカラムのみを出力できることが特徴です.
ex) csvtail file.csv  ## tail -f file.csv相当の挙動
    csvtail file.csv -c 1,2,3  ## 1,2,3カラムのみ表示
    csvtail file.csv -c 1,2,3 -d ","  ## 区切り文字を指定
    csvtail file.csv -c 1,2,3 -d "," -s 1  ## 1秒ごとに監視
`,
	Run: runCsvTail,
}

// csvtail のエントリーポイント
// ex) csvtail file.csv  ## tail -f file.csv相当の挙動
//
//	csvtail file.csv -c 1,2,3  ## 1,2,3カラムのみ表示
//	csvtail file.csv -c 1,2,3 -d ","  ## 区切り文字を指定
//	csvtail file.csv -c 1,2,3 -d "," -s 1  ## 1秒ごとに監視
func runCsvTail(cmd *cobra.Command, args []string) {
	// ファイル指定がない場合はエラー
	if len(args) == 0 {
		cmd.Help()
		os.Exit(1)
	}

	// Columnsが存在する場合、重複を除き昇順にソート
	fixColumns := uniqueSorted(columns)

	// OSから改行コードの設定
	var lineBreak string
	if runtime.GOOS == "windows" {
		lineBreak = "\r\n"
	} else {
		lineBreak = "\n"
	}

	// パラメータセット
	csvf := csvf.CsvfInf{
		FileName:  args[0],
		Columns:   fixColumns,
		Delimiter: delimiter,
		LineBreak: lineBreak,
		Seconds:   seconds,
	}

	// ファイル存在チェック
	if b, _ := csvf.FileExists(); !b {
		fmt.Printf("[%s]ファイルは存在しません.", csvf.FileName)
		os.Exit(1)
	}

	// ファイルオープン
	csvf.FileOpen()
	defer csvf.FilePtr.Close()

	// 末尾にSeek
	csvf.SeekEnd()

	// ファイル監視
	csvf.FileWatch()

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// helpの順番を変更
	rootCmd.SetHelpTemplate(`Usage:
  {{.UseLine}}

Flags:
{{.Flags.FlagUsages | trimTrailingWhitespaces}}

{{.Long}}`)

	// フラグ情報の設定
	rootCmd.Flags().IntSliceVarP(&columns, "columns", "c", []int{}, "表示するカラム番号")
	rootCmd.Flags().IntVarP(&seconds, "seconds", "s", 0, "監視間隔(秒)")
	rootCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "区切り文字")
}

// スライスの重複除去＆昇順ソート
func uniqueSorted(input []int) []int {
	// 引数チェック
	if len(input) == 0 {
		return []int{} // 空スライスを返しておく
	}

	// 重複を除去  struct{}はキーのみを持つmap メモリ効率よい
	uniqueMap := make(map[int]struct{})
	for _, val := range input {
		uniqueMap[val] = struct{}{} // mapは重複を許さないので、重複する値は上書きされている
	}

	// マップのキーをスライスに変換する
	uniqueSlice := make([]int, len(uniqueMap))
	i := 0
	for key := range uniqueMap {
		uniqueSlice[i] = key
		i++
	}

	// 昇順ソート
	sort.Ints(uniqueSlice)

	return uniqueSlice
}
