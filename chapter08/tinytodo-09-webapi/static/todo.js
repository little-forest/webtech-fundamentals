// defer属性によって、DOMの構築がすべて完了してから setup 関数を呼び出す
setup();

/**
 * 各要素にイベントハンドラを設定する。
 */
function setup() {
  fetchTodoItems();

  addEventListenerByQuery("#btn-add", "click", onClickAddButton);
  addEventListenerByQuery("#new-todo", "click", onClickAddInput);
  addEventListenerByQuery("#new-todo", "keydown", onKeydownAddInput);

  window.addEventListener("hashchange", onHashChange);  // <1>
}

/**
 * 関数をイベントハンドラとして登録する。
 * @param {string} 登録先要素を指定するクエリ
 * @param {string} イベント名
 * @param {function} ハンドラとして登録する関数
 */
function addEventListenerByQuery(query, eventName, callback) {
  const elements = document.querySelectorAll(query);
  for (let i = 0; i < elements.length; i++) {
    elements[i].addEventListener(eventName, callback);
  }
}

/**
 * ハッシュフラグメントが変化したときのイベントハンドラ。
 * @param {event} hashchangeイベントオブジェクト
 */
function onHashChange(event) {
  // 編集中のToDo項目があれば、編集をキャンセルする
  const oldInputElement = getToDoInputFromUrl(event.oldURL); // <1>
  if (isNotNull(oldInputElement)) {
    cancelTodoEdit(oldInputElement);
  }

  // フラグメントで指定されたToDo項目を編集状態にする
  const curInputElement = getToDoInputFromUrl(event.newURL); // <2>
  if (isNotNull(curInputElement)) {
    enableTodoInput(curInputElement);
  }
}

/**
 * URLのハッシュフラグメントから、対応するToDo Input要素を取得する。
 * @param {string} URL
 */
function getToDoInputFromUrl(url) {
  const match = new URL(url).hash.match(/^#edit\/([0-9a-f]+)$/);
  if (isNull(match)) {
    return;
  }
  const todoItemId = match[1];
  return document.getElementById(todoItemId);
}

/**
 * ToDo項目がクリックされた時のイベントハンドラ。
 * @param {Event} イベントオブジェクト
 */
function onClickTodoInput(event) {
  location.hash = `edit/${event.target.id}`;
}

/**
 * ToDo項目を編集可能にする。
 * save/cancelボタンも表示する。
 * @param {HTMLElement} ToDo項目のinput要素
 */
function enableTodoInput(todoInput) {
  // キャンセル時に戻すため、編集前の内容を保存
  if (isNull(todoInput.dataset.originalValue)) {
    todoInput.dataset.originalValue = todoInput.value;
  }

  todoInput.readOnly = false;
  todoInput.focus();

  // save/cancelボタンを表示
  const todoEditorControl = getTodoEditorControl(todoInput);
  if (isNotNull(todoEditorControl)) {
    todoEditorControl.classList.remove("hidden");
  } else {
    console.error("TodoEditorControl not found.");
  }
}

/**
 * ToDo項目を編集不能にする。
 * save/cancelボタンも非表示にする。
 * @param {HTMLElement} ToDo項目のinput要素
 */
function disableTodoInput(todoInput) {
  todoInput.readOnly = true;
  todoInput.blur();

  // save/cancelボタンを非表示
  const todoEditorControl = getTodoEditorControl(todoInput);
  if (isNotNull(todoEditorControl)) {
    todoEditorControl.classList.add("hidden");
  } else {
    console.error("TodoEditorControl not found.");
  }
}

/**
 * ToDo項目に対応するsave/cancelボタンのコンテナ要素を取得する。
 * @param {HTMLElement} ToDo項目のinput要素
 * @return {HTMLElement} save/cancelボタンの親div要素、存在しない場合はnull
 */
function getTodoEditorControl(todoInput) {
  const todoEditorControl = todoInput.nextElementSibling;
  if (todoEditorControl.classList.contains("todo-item-control")) {
    return todoEditorControl;
  }
  return null;
}

/**
 * Cancelボタンがクリックされた時のイベントハンドラ。
 * @param {Event} イベントオブジェクト
 */
function onClickCancelButton() {
  // ハッシュフラグメントをクリアすることで、
  // 編集中のToDoをキャンセルさせる
  location.hash = "";
}

/**
 * ToDo項目を編集をキャンセルする。
 * @param {HTMLElement} ToDo項目のinput要素
 */
function cancelTodoEdit(todoInput) {
  disableTodoInput(todoInput);
  // 編集内容を元に戻す
  todoInput.value = todoInput.dataset.originalValue;
}

/**
 * ToDoリストを取得して表示。
 */
function fetchTodoItems() {
  fetch("/todos/")
    .then(response => {
      if (!response.ok) {
        // エラー時はログイン画面に戻る
        location.href = "/login";
      }
      return response.json();
    })
    .then(data => {
      // 取得したToDoを画面に追加する
      data.items.forEach(todoItem => {
        addTodoItem(todoItem);
      });

      // 画面の状態を復元する
      restoreState();
    });
}


/**
 * ハッシュフラグメントから画面の状態を復元する
 */
function restoreState() {
  const curInputElement = getToDoInputFromUrl(location.href);
  if (isNotNull(curInputElement)) {
    enableTodoInput(curInputElement);
  }
}

/**
 * ToDo入力ボックスがクリックされた時のイベントハンドラ。
 */
function onClickAddInput() {
  // ハッシュフラグメントをクリアすることで、
  // 編集中のToDoをキャンセルさせる
  location.hash = "";
}

/**
 * ToDo入力時のキー入力イベントハンドラ。
 * @param {Event} キーボードイベントオブジェクト
 */
function onKeydownAddInput(event) {
  if (event.isComposing || event.keyCode === 229) {
    return;
  }
  if (event.code === "Enter") {
    onClickAddButton();
  }
}

/**
 * Addボタンがクリックされた時のイベントハンドラ。
 */
function onClickAddButton() {
  // 新しいToDoを入力したinput要素の取得
  const addInput = document.getElementById("new-todo");

  if (addInput.value.trim() === "") {
    return;
  }

  // POSTリクエストの準備
  const request = {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      todo: addInput.value,
    })
  };

  // リクエストの送信
  fetch("/todos/", request)                               // <1>
    .then(response => {
      if (!response.ok) {
        // エラー時はログイン画面に戻る
        location.href = "/login";
      }
      // LocationヘッダからToDo Idを取得する
      const location = response.headers.get("location");  // <2>
      const todoId = location.replace(/^\/todos\//, '');

      // 画面にToDoを追加する
      addTodoItem({                                       // <3>
        id: todoId,
        todo: addInput.value,
      });

      // 入力したToDoをクリア
      addInput.value = "";
    })
    .catch((err) => {
      console.error("Failed to send request: ", err);
    });
}

/**
 * ToDo項目を画面に追加する。
 * @param {Object} 追加するToDo項目
 */
function addTodoItem(todoItem) {
  // ToDo項目を表示するためのli要素を生成してCSSクラスを追加  // <1>
  const listElement = document.createElement("li");
  listElement.classList.add("todo-item");

  // li要素内部を構築                                         // <2>
  todoItem = `
    <div class="todo-item-container">
      <input type="checkbox" />
      <input type="text" class="todo" id="${todoItem.id}" value="${todoItem.todo}" readonly />
      <div class="todo-item-control hidden">
        <button class="btn-save">save</button>
        <button class="btn-cancel">cancel</button>
      </div>
    </div>
  `;
  listElement.insertAdjacentHTML("afterbegin", todoItem);    // <3>

  // 編集操作用のイベントハンドラを登録                      // <4>
  listElement.querySelector('input[type="text"].todo').addEventListener("click", onClickTodoInput);
  listElement.querySelector('input[type="text"].todo').addEventListener("keydown", onKeydownTodoEdit);
  listElement.querySelector('button.btn-cancel').addEventListener("click", onClickCancelButton);
  listElement.querySelector('button.btn-save').addEventListener("click", onClickSaveButton);

  // li要素を親のul要素に追加                                // <5>
  const todoListElement = document.querySelector("ul#todo-list");
  todoListElement.insertAdjacentElement("afterbegin", listElement);
}

/**
 * ToDo項目を画面に追加する。
 * @param {Event} キーボードイベント
 */
function onKeydownTodoEdit(event) {
  if (event.isComposing || event.keyCode === 229) {
    return;
  }
  if (event.code === "Enter") {
    const saveButton = event.target.nextElementSibling.querySelector('button.btn-save');
    saveButton.click();
  }
}

/**
 * Saveボタンがクリックされた時のイベントハンドラ。
 * @param {Event} イベントオブジェクト
 */
function onClickSaveButton(event) {
  const saveBtn = event.target;
  const todoInput = saveBtn.parentNode.previousElementSibling;
  todoInput.blur();

  // リクエストの準備
  const request = {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      id: todoInput.id,
      todo: todoInput.value,
    }),
  };

  // リクエストの送信
  fetch(`/todos/${todoInput.id}`, request)
    .then(response => {
      if (!response.ok) {
        // エラー時はログイン画面に戻る
        location.href = "/login";
      }
      todoInput.dataset.originalValue = todoInput.value;
      location.hash = "";
    })
    .catch((err) => {
      console.error("Failed to send request: ", err);
    });
}

/**
 * 値が null または undefined かどうかを判定する。
 *
 * @param value 判定する値
 * @returns {boolean} null または undefined の場合は true
 */
function isNull(value) {
  return value === null || typeof value === "undefined";
}

/**
 * 値が null または undefined ではないかどうかを判定する。
 *
 * @param value 判定する値
 * @returns {boolean} null または undefined の場合は false
 */
function isNotNull(value) {
  return !isNull(value);
}

