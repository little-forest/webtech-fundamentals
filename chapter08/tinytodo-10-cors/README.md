# tinytodo-10-cors : CORSに対応したTiny ToDo

**書籍での説明箇所** :「8.6.6 CORSに対応したTiny ToDo」p370

## 概要

Tiny Calからクロスオリジン通信でWebAPIを呼び出せるよう、[`tinytodo-09-webapi`](./tinytodo-09-webapi/) を、CORSに対応させたものです。

実行時、`ALLOWED_ORIGINS` 環境変数によって、許可オリジンをコンマ区切りで指定できます。

本サンプルは、以下のURLで実働公開しています。

https://tinytodo-10-cors.webtech.littleforest.jp/

## 主要コード

- **クライアントサイド**
  - [`static/todo.js`](./static/todo.js) ブラウザ上で動作するJavaScriptコード (リスト8.10, p372 / リスト8.11, p373)
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

## ソースコードの補足説明

### CORS関連処理

p375のHINTでは、サーバサイドの実装解説を省略すると書きました。実際の処理を参照したい方は、 [`main.go内のcheckCors()関数`](./main.go) を見てください。

checkCors()関数では、HTTPリクエストの `Origin` ヘッダ有無をチェックし、`Origin` ヘッダがあればその内容をチェックします。

許可するオリジンは、実行環境によって異なり、ローカル実行であれば localhost になるし、インターネット上に公開する場合はそのドメイン名になります。
そのため、 `getAllowedOrigins()` 関数で環境変数 (`ALLOWED_ORIGINS`)から取得するようにしてあります。
