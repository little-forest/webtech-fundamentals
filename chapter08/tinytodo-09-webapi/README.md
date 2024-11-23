# tinytodo-09-webapi : WebAPI化したTiny ToDo

**書籍での説明箇所** :「8.5 Tiny ToDo の Web API 化」p346

## 概要

[`tinytodo-08-spa`](../tinytodo-08-spa) をベースに、サーバとの通信をWebAPI化したものです。

利用者視点の動作は [`tinytodo-08-spa`](../tinytodo-08-spa) と同じですが、サーバ側のAPIをRESTを意識したものに変更しています。

本サンプルは、以下のURLで実働公開しています。

https://tinytodo-09-webapi.webtech.littleforest.jp/

## 主要コード

- **クライアントサイド**
  - [`static/todo.js`](./static/todo.js) ブラウザ上で動作するJavaScriptコード (リスト8.6, p347)
- **HTMLテンプレート**
  - [`templates/todo.html`](./templates/todo.html) ToDo画面のHTMLテンプレート
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
  - [`service_todo.go`](./service_todo.go) ToDoに関する処理 (リスト8.8, p350)
  - [`todo_item.go`](./todo_item.go) ToDo項目を表す構造体と関連処理
  - [`todo_list.go`](./todo_list.go) ToDoリストを表す構造体と関連処理
  - [`user_account.go`](./user_account.go) ユーザアカウントを表す構造体と関連処理
  - [`user_account_manager.go`](./user_account_manager.go) ユーザアカウント管理にまつわる処理
  - [`session.go`](./session.go) セッション情報を保持する構造体と関連処理
  - [`session_manager.go`](./session_manager.go) セッション管理にまつわる処理
