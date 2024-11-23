setup();

function setup() {
  // a要素がクリックされた時の処理
  addEventListenerByQuery('nav a', "click", (e) => {
    // ブラウザ本来の遷移動作を禁止する
    e.preventDefault();                               // <1>

    // a要素が指すリンクのパス部分を取得し、ルーティングする
    route(e.target.pathname, true);                   // <2>

    console.log(`onClick : ${e.target.pathname}`);
  });

  // 現在のURLに応じてルーティングする
  route(location.pathname, false);                     // <1>
  console.log(`setup : ${location.pathname}`);

  // ページ履歴が辿られたとき、URLに応じて表示内容を切り替える
  window.addEventListener("popstate", () => {          // <2>
    console.log(`onPopState : ${location.pathname}`);
    route(location.pathname, false);
  });
}

/**
 * 簡易的なルーティング関数。
 * @param {path} URLのパス部分
 * @param {pushState} trueのとき、URLを変更する
 */
function route(path, pushState) {
  // パスに応じて表示を切り替える
  switch (path) {                            // <1>
    case "/page1":
      showContant("Page1 content");
      break;
    case "/page2":
      showContant("Page2 content");
      break;
    case "/page3":
      showContant("Page3 content");
      break;
  }

  // History APIでURLを変更する
  if (pushState) {
    history.pushState(null, "", path);       // <3>
  }
}

/**
 * ページ内容を表示する関数。
 * @param {text} 表示内容
 */
function showContant(text) {
  document.getElementById("main").innerHTML = text;  // <2>
}

function addEventListenerByQuery(query, eventName, func) {
  const elements = document.querySelectorAll(query);
  for (let i = 0; i < elements.length; i++) {
    elements[i].addEventListener(eventName, func);
  }
}
