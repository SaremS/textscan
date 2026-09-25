const express = require("express");

const app = express();
app.use(express.json());

function combine(a, b) {
  for (const k in b) {
    const v = b[k];

    if (v && typeof v === "object" && !Array.isArray(v)) {
      if (!a[k] || typeof a[k] !== "object") {
        a[k] = {};
      }

      combine(a[k], v);
    } else {
      a[k] = v;
    }
  }

  return a;
}

const state = {
  applicationName: "Configuration Service",
  debug: false
};

app.get("/", (req, res) => {
  res.json({
    name: "Configuration Service",
    endpoints: [
      "POST /api/config",
      "GET /api/config",
      "GET /api/status"
    ]
  });
});

app.post("/api/config", (req, res) => {
  combine(state, req.body);

  res.json({
    message: "Configuration updated",
    data: state
  });
});

app.get("/api/config", (req, res) => {
  res.json(state);
});

app.get("/api/status", (req, res) => {
  const item = {};

  res.json({
    local: Object.prototype.hasOwnProperty.call(item, "isAdmin"),
    value: item.isAdmin ?? null,
    active: item.isAdmin !== undefined
  });
});

app.listen(3000, "127.0.0.1", () => {
  console.log("Server listening on http://127.0.0.1:3000");
});
