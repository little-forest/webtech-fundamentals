# tinycal : Tiny ToDoのAPIを呼び出すカレンダーアプリケーション

**書籍での説明箇所** :「8.6.3 クロスオリジン通信」p359

※ 書籍内ではソースコード自体は解説していません

## 概要

クロスオリジン通信の動作を確認するためのサンプルアプリケーションです。
Tiny ToDoのWebAPIを呼び出すことで、Tiny ToDoが管理するToDoリストを表示します。

本サンプルは、以下のURLで実働公開しています。

https://tinycal.webtech.littleforest.jp/

## 主要コード

- [`static/tinycal.js`](./static/tinycal.js) カレンダーを表示するJavaScriptコード
- [`static/todo.js`](./static/todo.js) Tiny ToDoのWebAPIを呼び出すJavaScriptコード
- [`static/index.html`](./static/index.html) Tiny Calの画面を表示するHTML

