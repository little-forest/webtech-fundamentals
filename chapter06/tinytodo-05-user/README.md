# tinytodo-05-user : ユーザー管理機能を導入したTiny ToDo

**書籍での説明箇所** :「6.7 ユーザ認証の実装」p214

## 概要

ユーザーアカウントの管理機能を導入したTiny ToDoです。
[`tinytodo-04-session`](../tinytodo-04-session/) に対して、以下のような機能を追加しています。

- ユーザーアカウント発行機能
- ログイン/ログアウト機能

本サンプルは、理解しやすくするために最低限の機能に絞っており、
インターネットを介して不特定大数の利用者へ公開するには不適切です。

実働するサンプルとして公開しているTiny ToDoは [`tinytodo-05-user-final`](../tinytodo-05-user-final/) です。

## 主要コード

- **HTMLテンプレート**
  - [`templates/todo.html`](./templates/todo.html) ToDo画面のHTMLテンプレート
  - [`templates/login.html`](./templates/login.html) ログイン画面のHTMLテンプレート (リスト6.16, p218)
  - [`templates/create-user-account.html`](./templates/create-user-account.html) アカウント作成画面のHTMLテンプレート
  - [`templates/new-user-account.html`](./templates/new-user-account.html) 作成アカウント表示画面のHTMLテンプレート
- **サーバサイド**
  - [`main.go`](./main.go) サーバサイドのメイン処理
  - [`args.go`](./args.go) コマンドライン引数の処理
  - [`page_create_account.go`](./page_create_account.go) アカウント作成画面に関する処理
  - [`page_login.go`](./page_login.go) ログイン画面に関する処理 (リスト6.17〜18, p219, p221)
  - [`page_new_account.go`](./page_new_account.go) 作成アカウント表示画面に関する処理
  - [`page_todo.go`](./page_todo.go) ToDo画面に関する処理 (リスト6.21, p225)
  - [`session.go`](./session.go) セッション情報を保持する構造体と関連処理 (リスト6.12〜13, p214, p216 / リスト6.19, 222)
  - [`session_manager.go`](./session_manager.go) セッション管理にまつわる処理 (リスト6.14, p216 / リスト6.22, p226)
  - [`user_account.go`](./user_account.go) ユーザアカウントを表す構造体と関連処理 (リスト6.15, p217)
  - [`user_account_manager.go`](./user_account_manager.go) ユーザアカウント管理にまつわる処理 (リスト6.20, p223)

