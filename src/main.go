package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	// 管理者権限のチェック
	isAdmin := checkAdmin()
	if !isAdmin {
		fmt.Println("This program requires administrator privileges. Restarting with admin rights...")
		// 管理者権限で再起動
		err := runAsAdmin()
		if err != nil {
			fmt.Printf("Failed to restart with admin privileges: %v\n", err)
		}
		return
	}

	// 管理者権限がある場合、続行
	fmt.Println("Running with administrator privileges!")

	// 探索を開始するディレクトリ
	startDir := "C:\\" // 必要に応じて変更！

	// 探索したいファイル名
	targetFile := "sprotect.sys"

	fmt.Printf("Searching for '%s' in '%s'...\n", targetFile, startDir)

	// 時間計測スタート
	startTime := time.Now()

	// ワーカー用のWaitGroup
	var wg sync.WaitGroup

	// 結果を格納するチャンネル
	results := make(chan string, 100)
	errors := make(chan error, 10)

	// ディレクトリを探索してゴルーチンで処理
	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// エラー（アクセス拒否など）をチャンネルに送信
			errors <- err
			return nil // 次のファイルへ
		}

		// ディレクトリならゴルーチンで処理
		if info.IsDir() {
			wg.Add(1)
			go func(dirPath string) {
				defer wg.Done()
				searchInDirectory(dirPath, targetFile, results, errors)
			}(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking initial directory: %v\n", err)
		return
	}

	// ゴルーチンの終了を待つ
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// 結果を集計
	var foundFiles []string
	for res := range results {
		foundFiles = append(foundFiles, res)
	}

	// エラーの確認
	for err := range errors {
		fmt.Printf("Error: %v\n", err)
	}

	// 結果を表示
	if len(foundFiles) > 0 {
		fmt.Println("\nFound the following matches:")
		for _, file := range foundFiles {
			fmt.Println(file)
		}
	} else {
		fmt.Println("\nNo matches found.")
	}

	// 検索時間を表示
	elapsedTime := time.Since(startTime)
	fmt.Printf("\nSearch completed in %s.\n", elapsedTime)

	// プログラムの終了を待つ
	fmt.Print(`
   /\_/\
  ( o.o )  < Thank you for using the program!
   > ^ <     Press Enter to exit...
`)
	fmt.Scanln() // Enterキーを待つ
}

// searchInDirectory は1つのディレクトリ内を探索します
func searchInDirectory(dir string, target string, results chan<- string, errors chan<- error) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return nil // アクセス拒否を無視
			}
			errors <- err
			return nil
		}
		// ファイル名が一致した場合
		if strings.EqualFold(info.Name(), target) {
			results <- path
		}
		return nil
	})
	if err != nil {
		errors <- err
	}
}

// checkAdmin は現在のプロセスが管理者権限で実行されているかを確認します
func checkAdmin() bool {
	// 管理者権限を確認するため、物理ディスクにアクセスできるか試みる
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

// runAsAdmin は管理者権限でプログラムを再起動します
func runAsAdmin() error {
	// 現在の実行ファイルのパスを取得
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	// 管理者権限で実行するコマンドを準備
	cmd := exec.Command("powershell", "-Command", "Start-Process", executable, "-Verb", "runAs")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return cmd.Run()
}
