# csvtail

学習用のcsvファイルをtailするコマンド用プログラム。

## Usage

```bash
csvtail file.csv  ## tail -f file.csv相当の挙動
csvtail file.csv -c 1,2,3  ## 1,2,3カラムのみ表示
csvtail file.csv -c 1,2,3 -d ","  ## 区切り文字を指定
csvtail file.csv -c 1,2,3 -d "," -s 1  ## 1秒ごとに監視
