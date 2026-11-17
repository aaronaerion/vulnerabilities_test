import hashlib
import os
import pickle
import sqlite3
import subprocess

import requests
import urllib3
import yaml
from django.conf import settings
from flask import Flask, request
from jinja2 import Template

if not settings.configured:
    settings.configure(DEBUG=True, SECRET_KEY="not-a-real-password")

app = Flask(__name__)
app.config["DEBUG"] = True
app.config["SECRET_KEY"] = "not-a-real-password"

AWS_ACCESS_KEY_ID = "AKIAIOSFODNN7EXAMPLE"
DB_PASSWORD = "not-a-real-password"


@app.route("/search")
def search():
    q = request.args.get("q", "")
    conn = sqlite3.connect(":memory:")
    cur = conn.cursor()
    cur.execute("SELECT * FROM items WHERE name = '" + q + "'")
    return str(cur.fetchall())


@app.route("/run")
def run():
    name = request.args.get("name", "world")
    return subprocess.getoutput("echo " + name)


@app.route("/file")
def file():
    path = request.args.get("path", "readme.txt")
    with open(path, "r") as handle:
        return handle.read()


@app.route("/preview")
def preview():
    msg = request.args.get("msg", "")
    return "<html><body>" + msg + "</body></html>"


@app.route("/fetch")
def fetch():
    url = request.args.get("url", "http://127.0.0.1/")
    return requests.get(url).text


@app.route("/hash")
def hash_password():
    password = request.args.get("password", "")
    return hashlib.md5(password.encode()).hexdigest()


@app.route("/load", methods=["POST"])
def load():
    raw = request.get_data()
    if request.args.get("fmt") == "yaml":
        return str(yaml.load(raw))
    return str(pickle.loads(raw))


@app.route("/render")
def render():
    body = request.args.get("body", "")
    return Template(body).render()


@app.route("/pool")
def pool():
    http = urllib3.PoolManager()
    target = request.args.get("url", "http://127.0.0.1/")
    return http.request("GET", target).data


if __name__ == "__main__":
    app.run(debug=True)
