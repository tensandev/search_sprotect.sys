# search_sprotect.sys

このプログラムは、指定されたディレクトリ内で特定のファイル（デフォルトでは`sprotect.sys`）を検索するGo言語で書かれたツールです。管理者権限が必要な場合は、自動的に再起動して管理者権限で実行します。

## 機能

- 指定されたディレクトリ内でファイルを再帰的に検索
- 管理者権限の確認と自動再起動
- 検索結果とエラーの表示
- 検索にかかった時間の計測

## 使用方法

1. プロジェクトをクローンまたはダウンロードします。
2. Goがインストールされていることを確認します。
3. 以下のコマンドでプログラムを実行します。

```bash
go run src/main.go
```

## 設定

- 検索を開始するディレクトリは、`main.go`内の`startDir`変数で設定できます。デフォルトは`C:\`です。
- 検索対象のファイル名は、`main.go`内の`targetFile`変数で設定できます。デフォルトは`sprotect.sys`です。

## 出力

- 検索が成功した場合、見つかったファイルのパスが一覧表示されます。
- エラーが発生した場合、エラー内容が表示されます。
- 検索にかかった時間が表示されます。

## 注意事項

- このプログラムはWindows環境で動作することを前提としています。
- 管理者権限が必要な場合、PowerShellを使用して再起動します。

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。詳細は[LICENSE](LICENSE)ファイルをご確認ください。

## お礼

```
   /\_/\
  ( o.o )  < Thank you for using the program!
   > ^ <
```