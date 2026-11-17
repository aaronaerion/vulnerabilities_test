const crypto = require("crypto");
const fs = require("fs");
const { exec } = require("child_process");
const express = require("express");
const axios = require("axios");
const jwt = require("jsonwebtoken");
const _ = require("lodash");
const minimist = require("minimist");

const app = express();
app.use(express.text({ type: "*/*" }));

const AWS_ACCESS_KEY_ID = "AKIA4B7X2Q9MPLKJ0WNR";
const AWS_SECRET_ACCESS_KEY = "FAKESECRET_i1j2k3l4m5n6o7p8q9r0";
const DB_PASSWORD = "not-a-real-password";
const JWT_SECRET = "not-a-real-password";

app.use((req, res, next) => {
  res.setHeader("Access-Control-Allow-Origin", "*");
  next();
});

app.get("/search", (req, res) => {
  const q = req.query.q || "";
  const sql = "SELECT * FROM items WHERE name = '" + q + "'";
  res.send(sql);
});

app.get("/run", (req, res) => {
  const name = req.query.name || "world";
  exec("echo " + name, (err, stdout) => {
    res.send(err ? String(err) : stdout);
  });
});

app.get("/file", (req, res) => {
  const path = req.query.path || "readme.txt";
  res.send(fs.readFileSync(path, "utf8"));
});

app.get("/preview", (req, res) => {
  const msg = req.query.msg || "";
  res.send("<html><body>" + msg + "</body></html>");
});

app.get("/fetch", async (req, res) => {
  const url = req.query.url || "http://127.0.0.1/";
  const response = await axios.get(url);
  res.send(response.data);
});

app.get("/hash", (req, res) => {
  const password = req.query.password || "";
  res.send(crypto.createHash("md5").update(password).digest("hex"));
});

app.post("/load", (req, res) => {
  const body = typeof req.body === "string" ? req.body : "";
  res.send(String(eval(body)));
});

app.get("/merge", (req, res) => {
  const defaults = { role: "user" };
  const merged = _.merge({}, defaults, req.query);
  res.json(merged);
});

app.get("/token", (req, res) => {
  const user = req.query.user || "guest";
  res.send(jwt.sign({ user }, JWT_SECRET));
});

app.get("/args", (req, res) => {
  const parsed = minimist(process.argv.slice(2));
  res.json(parsed);
});

app.listen(3000);
