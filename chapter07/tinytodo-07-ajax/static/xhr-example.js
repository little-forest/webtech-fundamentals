function sendXhr() {
  const xhr = new XMLHttpRequest();    // <1>
  xhr.open("GET", "/todo");            // <2>
  xhr.onload = () => {                 // <3>
    console.log("Receive response");   // <6>
    console.log(xhr.responseText);
  };
  xhr.send();                          // <4>
  console.log("Send request");         // <5>
}
