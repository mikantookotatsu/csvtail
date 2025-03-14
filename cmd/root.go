/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mikan-to-kotatsu/csvtail/csvf"
	"github.com/spf13/cobra"
)

var (
	seconds   int    // -s オプションで指定された秒数
	columns   []int  // -c オプションで指定されたカラム
	delimiter string // -d オプションで指定された区切り文字
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "csvtail",
	Short: "csvtail は CSVファイルを tail -f 的に監視するツールです.",
	Long: `csvtail は CSVファイルを tail -f 的に監視するツールです.
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
		fmt.Println("監視ファイルを指定してください. ex) csvtail file.csv")
		os.Exit(1)
	}

	// パラメータセット
	csvf := csvf.CsvfInf{
		FileName:  args[0],
		Columns:   columns,
		Delimiter: delimiter,
		Seconds:   seconds,
	}

	// ファイル存在チェック
	if !csvf.FileExists() {
		fmt.Printf("[%s]ファイルは存在しません.", csvf.FileName)
		os.Exit(1)
	}

	// ファイルオープン
	csvf.FileOpen()

	// 末尾にSeek
	csvf.SeekEnd()

	// ファイル監視
	csvf.FileWatch()

	// ファイルクローズ
	defer csvf.FilePtr.Close()

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
	// フラグ情報の設定
	rootCmd.Flags().IntSliceVarP(&columns, "columns", "c", []int{}, "表示するカラム番号")
	rootCmd.Flags().IntVarP(&seconds, "seconds", "s", 1, "監視間隔(秒)")
	rootCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "区切り文字")
}
