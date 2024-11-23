// defer属性によって、DOMの構築がすべて完了してから setup 関数を呼び出す
setup();

// 編集しているToDo項目のinput要素
let currentEditingTodo;

/**
 * 各要素にイベントハンドラを設定する。
 */
function setup() {
  addEventListenerByQuery('input[type="text"].todo', "click", onClickTodoInput);
  addEventListenerByQuery('button.btn-cancel', "click", onClickCancelButton);
  addEventListenerByQuery('button.btn-save', "click", onClickSaveButton);
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
 * ToDo項目がクリックされた時のイベントハンドラ。
 * @param {Event} イベントオブジェクト
 */
function onClickTodoInput(event) {
  if (isNotNull(currentEditingTodo) && currentEditingTodo !== event.target) {
    cancelTodoEdit(currentEditingTodo);
  }

  enableTodoInput(event.target);
  currentEditingTodo = event.target;
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
function onClickCancelButton(event) {
  const cancelBtn = event.target;
  const todoInput = cancelBtn.parentNode.previousElementSibling;
  cancelTodoEdit(todoInput);
}

/**
 * ToDo項目を編集をキャンセルする。
 * @param {HTMLElement} ToDo項目のinput要素
 */
function cancelTodoEdit(todoInput) {
  disableTodoInput(todoInput);
  // 編集内容を元に戻す
  todoInput.value = todoInput.dataset.originalValue;
  currentEditingTodo = null;
}

/**
 * Saveボタンがクリックされた時のイベントハンドラ。
 * @param {Event} イベントオブジェクト
 */
function onClickSaveButton(event) {
  const saveBtn = event.target;
  const todoInput = saveBtn.parentNode.previousElementSibling;
  disableTodoInput(todoInput);
  currentEditingTodo = null;
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

