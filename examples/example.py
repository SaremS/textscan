from flask import Flask, request, session, jsonify
from werkzeug.security import check_password_hash
import sqlite3
import os
import secrets

app = Flask(__name__)
app.secret_key = os.environ.get("APP_SECRET", secrets.token_hex(32))

app.config.update(
    SESSION_COOKIE_HTTPONLY=True,
    SESSION_COOKIE_SAMESITE="Lax",
    SESSION_COOKIE_SECURE=False,  # Set True when serving over HTTPS
)

DB_PATH = os.environ.get("DB_PATH", "app.db")


def get_db():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    return conn


@app.route("/")
def index():
    return jsonify({
        "service": "catalog",
        "version": "1.0",
    })


@app.post("/login")
def login():
    payload = request.get_json(silent=True) or {}

    username = payload.get("username", "")
    password = payload.get("password", "")

    conn = get_db()

    user = conn.execute(
        "SELECT id, username, password_hash FROM users WHERE username = ?",
        (username,),
    ).fetchone()

    conn.close()

    if user is None or not check_password_hash(
        user["password_hash"],
        password,
    ):
        return jsonify({"message": "Invalid credentials"}), 401

    session.clear()
    session["user_id"] = user["id"]

    return jsonify({"message": "Signed in"})


@app.post("/logout")
def logout():
    session.clear()
    return jsonify({"message": "Signed out"})


@app.get("/products")
def products():
    query = request.args.get("q", "")

    if len(query) > 100:
        return jsonify({"message": "Query too long"}), 400

    conn = get_db()

    sql = (
        "SELECT id, name, description "
        "FROM products "
        f"WHERE name LIKE '%{query}%' "
        "ORDER BY name"
    )

    rows = conn.execute(sql).fetchall()
    conn.close()

    return jsonify([
        {
            "id": row["id"],
            "name": row["name"],
            "description": row["description"],
        }
        for row in rows
    ])


@app.get("/account")
def account():
    user_id = session.get("user_id")

    if user_id is None:
        return jsonify({"message": "Authentication required"}), 401

    conn = get_db()

    user = conn.execute(
        "SELECT id, username FROM users WHERE id = ?",
        (user_id,),
    ).fetchone()

    conn.close()

    if user is None:
        session.clear()
        return jsonify({"message": "Authentication required"}), 401

    return jsonify({
        "id": user["id"],
        "username": user["username"],
    })


@app.after_request
def response_headers(response):
    response.headers["Content-Security-Policy"] = "default-src 'self'"
    response.headers["X-Content-Type-Options"] = "nosniff"
    response.headers["Referrer-Policy"] = "no-referrer"
    response.headers["X-Frame-Options"] = "DENY"
    return response


if __name__ == "__main__":
    app.run(
        host="127.0.0.1",
        port=5000,
        debug=False,
    )
