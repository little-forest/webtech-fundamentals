# tinytodo-08-spa : SPA化したTiny ToDo

**書籍での説明箇所** :「7.6.2 ToDo の追加を JSON でやりとりする」p282

## 概要

ToDo管理をすべてSPA化したTiny ToDo。
ToDoの作成・編集処理をクライアントサイドJavaScriptで実現し、FetchAPIを使ってサーバサイドとAjax通信をしています。

ログインやアカウント作成に関する処理は、従来型Webアプリケーションのままです。

本サンプルは、以下のURLで実働公開しています。

https://tinytodo-08-spa.webtech.littleforest.jp/

## 主要コード

- **クライアントサイド**
  - [`static/todo.js`](./static/todo.js) ブラウザ上で動作するJavaScriptコード (リスト7.31, p283 / リスト7.32, p285 / リスト7.38, p290 / リスト7.40, p296 / リスト7.41, p297 / リスト7.42, p297)
- **HTMLテンプレート**
  - [`templates/todo.html`](./templates/todo.html) ToDo画面のHTMLテンプレート (リスト7.30, p283)
  - [`templates/login.html`](./templates/login.html) ログイン画面のHTMLテンプレート
  - [`templates/create-user-account.html`](./templates/create-user-account.html) アカウント作成画面のHTMLテンプレート
  - [`templates/new-user-account.html`](./templates/new-user-account.html) 作成アカウント表示画面のHTMLテンプレート
- **サーバサイド**
  - [`main.go`](./main.go) サーバサイドのメイン処理
  - [`args.go`](./args.go) コマンドライン引数の処理
  - [`page_create_account.go`](./page_create_account.go) アカウント作成画面に関する処理
  - [`page_login.go`](./page_login.go) ログイン画面に関する処理
  - [`page_new_account.go`](./page_new_account.go) 作成アカウント表示画面に関する処理
  - [`page_todo.go`](./page_todo.go) ToDo画面に関する処理
  - [`service_todo.go`](./service_todo.go) ToDoに関する処理 (リスト7.26, p288 / リスト7.37, p289)
  - [`todo_item.go`](./todo_item.go) ToDo項目を表す構造体と関連処理 (リスト7.35, p287)
  - [`todo_list.go`](./todo_list.go) ToDoリストを表す構造体と関連処理
  - [`user_account.go`](./user_account.go) ユーザアカウントを表す構造体と関連処理
  - [`user_account_manager.go`](./user_account_manager.go) ユーザアカウント管理にまつわる処理
  - [`session.go`](./session.go) セッション情報を保持する構造体と関連処理
  - [`session_manager.go`](./session_manager.go) セッション管理にまつわる処理

