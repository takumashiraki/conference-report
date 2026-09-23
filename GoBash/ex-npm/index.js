// ex-Go/main.go と同じことをする実装
const express = require("express");

const app = express();

app.get("/", (req, res) => {
  res.send("Hello, Node server!\n");
});

app.listen(8080, () => {
  console.log("server started at http://localhost:8080");
});
