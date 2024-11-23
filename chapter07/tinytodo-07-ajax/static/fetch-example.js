function doFetch() {
  fetch("/todo")
    .then((response) => response.text())   // <1>
    .then((text) => {                      // <2>
      console.log("Receive response");
      console.log(text);
    });
  console.log("Send request");
}
