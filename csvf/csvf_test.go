package csvf

import (
	"os"
	"testing"
)

// memo : Fatal系はテストを中断、Error系はテストは続行

// FileExists() のテスト
func TestFileExists_FileExists(t *testing.T) {
	// ファイル有り時のテスト
	// テスト用のファイルを作成
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("FileExists() テスト用ファイル作成失敗: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // テスト終了後にファイルを削除
	tmpFile.Close()

	c := CsvfInf{FileName: tmpFile.Name()}
	exists, err := c.FileExists()
	if err != nil {
		t.Errorf("FileExists() エラーが発生しました: %v", err)
	}
	if !exists {
		t.Errorf("FileExists() ファイル存在時に false が返されました")
	}
}

func TestFileExists_FileNotExists(t *testing.T) {
	// ファイル無し時のテスト
	c := CsvfInf{FileName: "nonexistentfile"}
	exists, err := c.FileExists()
	if err != nil {
		t.Errorf("FileExists() でエラーが発生しました: %v", err)
	}
	if exists {
		t.Errorf("FileExists() ファイル無し時に true が返されました")
	}
}

func TestFileExists_Directory(t *testing.T) {
	// ディレクトリ指定時のテスト
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("FileExists() テスト用ディレクトリ作成失敗: %v", err)
	}
	defer os.RemoveAll(tmpDir) // テスト終了後にディレクトリを削除

	c := CsvfInf{FileName: tmpDir}
	exists, err := c.FileExists()
	if err != nil {
		t.Errorf("FileExists() でエラーが発生しました: %v", err)
	}
	if exists {
		t.Errorf("FileExists() ディレクトリ指定時に true が返されました")
	}
}

// FileOpen()のテスト
func TestFileOpen_Success(t *testing.T) {
	// ファイルが正常にオープンできるかのテスト
	// テスト用のファイルを作成
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("FileOpen() テストファイル作成失敗: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // テスト終了後にファイルを削除
	tmpFile.Close()

	c := CsvfInf{FileName: tmpFile.Name()}
	err = c.FileOpen()
	if err != nil {
		t.Errorf("FileOpen() でエラーが発生しました: %v", err)
	}
	c.FileClose()
}

func TestFileOpen_FileNotExists(t *testing.T) {
	// ファイルが存在しない場合にエラーを返すかのテスト
	c := CsvfInf{FileName: "nonexistentfile"}
	err := c.FileOpen()
	if err == nil {
		t.Errorf("FileOpen() はエラーを返すはずですが、nil を返しました")
	}
}
