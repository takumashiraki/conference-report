const express = require("express");  // サーバーを起動するパッケージ

const app = express();

app.get("/", (req, res) => {
  res.send("Hello, Node server!\n");
});

app.listen(8080, () => {
  console.log("Go Bash vol.3");
});

// https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=vscode&wt=none&l=text%2Fx-go&width=680&ds=true&dsyoff=0px&dsblur=0px&wc=true&wa=true&pv=0px&ph=0px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=2x&wm=false&code=package%2520main%250A%250Aimport%2520%28%250A%2509%2522fmt%2522%250A%2509%2522log%2522%250A%2509%2522net%252Fhttp%2522%250A%250A%2509%2522github.com%252Fgo-chi%252Fchi%252Fv5%2522%2520%252F%252F%2520%25E3%2582%25B5%25E3%2583%25BC%25E3%2583%2590%25E3%2583%25BC%25E3%2582%2592%25E8%25B5%25B7%25E5%258B%2595%25E3%2581%2599%25E3%2582%258B%25E3%2583%2591%25E3%2583%2583%25E3%2582%25B1%25E3%2583%25BC%25E3%2582%25B8%250A%29%250A%250Afunc%2520main%28%29%2520%257B%250A%2509r%2520%253A%253D%2520chi.NewRouter%28%29%250A%2509r.Get%28%2522%252F%2522%252C%2520func%28w%2520http.ResponseWriter%252C%2520req%2520*http.Request%29%2520%257B%250A%2509%2509fmt.Fprintln%28w%252C%2520%2522Hello%252C%2520Go%2520server%21%2522%29%250A%2509%257D%29%250A%250A%2509log.Println%28%2522Go%2520Bash%2520vol.3%2522%29%250A%2509log.Fatal%28http.ListenAndServe%28%2522%253A8080%2522%252C%2520r%29%29%250A%257D
