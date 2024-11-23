# tinytodo-05-user-final : ユーザー管理機能を導入したTiny ToDo(完成版)

**書籍での説明箇所** :「6.8.1 Webアプリケーションの複雑さ」p229

## 概要

[`tinytodo-05-user`](../tinytodo-05-user/) に対して、不特定多数に対して公開できるように機能追加をしたものです。

本サンプルは、以下のURLで実働公開しています。

https://tinytodo-05-user.webtech.littleforest.jp/

### 追加機能

追加機能は以下の通りです。
以降のサンプルコードは、本アプリケーションをベースに拡張しています。

- セッション関連機能
  - セッションの延長機能
  - 期限切れセッションの自動削除
  - セッション作成時に、既存のセッションと重複していないことを確認する機能
  - セッションIDの検証機能
- ユーザアカウント管理機能
  - 期限切れユーザアカウントの自動削除機能
- その他
  - 複数ユーザによる同時アクセスで、管理する情報が崩れないようにするための排他処理
  - リクエストのロギング機能
  - ログイン中に/loginにアクセスしたらToDo画面へ遷移する機能
  - ToDoListを構造体ではなくUserAccount内のstringスライスに変更
  - 一度ログインしたアカウントIDをCookieへ保存する機能
  - Cookieにsecure属性を追加する機能

## 主要コード

- **クライアントサイド**
  - [`static/todo.js`](./static/todo.js) ブラウザ上で動作するJavaScriptコード
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
  - [`service_todo.go`](./service_todo.go) ToDoに関する処理
  - [`todo_item.go`](./todo_item.go) ToDo項目を表す構造体と関連処理
  - [`todo_list.go`](./todo_list.go) ToDoリストを表す構造体と関連処理
  - [`user_account.go`](./user_account.go) ユーザアカウントを表す構造体と関連処理
  - [`user_account_manager.go`](./user_account_manager.go) ユーザアカウント管理にまつわる処理
  - [`session.go`](./session.go) セッション情報を保持する構造体と関連処理
  - [`session_manager.go`](./session_manager.go) セッション管理にまつわる処理

