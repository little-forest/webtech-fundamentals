document.addEventListener('DOMContentLoaded', setup);

function setup() {
  const reloadButton = document.getElementById("btn-reload");
  reloadButton.addEventListener("click", fetchTodo);

  const preflightButton = document.getElementById("btn-preflight");
  preflightButton.addEventListener("click", preflight);
  fetchTodo();
}

function fetchTodo() {
  showError("");
  clearTodo();

  const corsEnable = document.getElementById("cors-enable");

  let url;
  let fetchOpt;

  if (corsEnable.checked) {
    url = "https://tinytodo-10-cors.webtech.littleforest.jp/todos/";
    fetchOpt = { mode: "cors", credentials: "include" };
  } else {
    url = "https://tinytodo-09-webapi.webtech.littleforest.jp/todos/";
    fetchOpt = { credentials: "include" };
  }

  fetch(url, fetchOpt)
    .then(response => {
      return response.json();
    })
    .then(data => {
      data.items.forEach(todoItem => {
        addTodoItem(todoItem);
      });
    })
    .catch(() => {
      showError(`Failed to fetch URL : ${url}`);
    });
}

function clearTodo() {
  const todoList = document.querySelector("ul#todo-list");
  while (todoList.firstChild) {
    todoList.removeChild(todoList.firstChild);
  }
}

function addTodoItem(todoItem) {
  const listElement = document.createElement("li");

  todoItem = `
    <span>${todoItem.todo}</span>
  `;
  listElement.insertAdjacentHTML("afterbegin", todoItem);

  const todoListElement = document.querySelector("ul#todo-list");
  todoListElement.insertAdjacentElement("beforeend", listElement);
}

function showError(message) {
  const errorElement = document.getElementById("error-msg");
  errorElement.innerHTML = message;
}

function preflight() {
  fetch("https://tinytodo-10-cors.webtech.littleforest.jp/todos/", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    mode: "cors",
    credentials: "include",
    body: JSON.stringify({
      todo: "CORS test",
    })
  })
  .then(response => {
    return response.text();
  })
  .catch((err) => {
    console.error("Faile to request : ", err);
  });
}
