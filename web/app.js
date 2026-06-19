"use strict";

// Set window.API_BASE before loading this script to point at a custom backend.
// Otherwise: same-origin when served from the Go server, else the local backend.
const API_BASE =
  window.API_BASE ||
  (location.protocol.startsWith("http") && location.host === "localhost:8080"
    ? location.origin
    : "http://localhost:8080");

const TOKENS = {
  get access() { return localStorage.getItem("access_token"); },
  get refresh() { return localStorage.getItem("refresh_token"); },
  set({ access_token, refresh_token }) {
    if (access_token) localStorage.setItem("access_token", access_token);
    if (refresh_token) localStorage.setItem("refresh_token", refresh_token);
  },
  clear() {
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
  },
};

let ME = null;          // {id, login, games_played, wins}
let pollTimer = null;   // multiplayer polling

// ---------------------------------------------------------------- HTTP layer

async function api(path, { method = "GET", body, auth = true, retry = true } = {}) {
  const headers = {};
  if (body !== undefined) headers["Content-Type"] = "application/json";
  if (auth && TOKENS.access) headers["Authorization"] = `Bearer ${TOKENS.access}`;

  const res = await fetch(API_BASE + path, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  // Try to transparently refresh an expired access token once.
  if (res.status === 401 && auth && retry && TOKENS.refresh) {
    const refreshed = await tryRefresh();
    if (refreshed) return api(path, { method, body, auth, retry: false });
  }

  const text = await res.text();
  const data = text ? safeJSON(text) : null;

  if (!res.ok) {
    const msg = (data && (data.message || data.error)) || text || `HTTP ${res.status}`;
    const err = new Error(msg);
    err.status = res.status;
    throw err;
  }
  return data;
}

function safeJSON(t) { try { return JSON.parse(t); } catch { return t; } }

async function tryRefresh() {
  try {
    const data = await api("/api/auth/refresh/access", {
      method: "POST",
      auth: false,
      retry: false,
      body: { refresh_token: TOKENS.refresh },
    });
    if (data && data.access_token) {
      TOKENS.set(data);
      return true;
    }
  } catch { /* fall through */ }
  return false;
}

// ---------------------------------------------------------------- UI helpers

const $ = (sel, root = document) => root.querySelector(sel);
const app = $("#app");

function tpl(id) {
  return document.importNode($("#" + id).content, true);
}

function render(node) {
  stopPolling();
  app.innerHTML = "";
  app.appendChild(node);
}

let toastTimer;
function toast(msg, kind = "") {
  const el = $("#toast");
  el.textContent = msg;
  el.className = "toast" + (kind ? " " + kind : "");
  el.hidden = false;
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => (el.hidden = true), 3200);
}

function spinner() {
  const d = document.createElement("div");
  d.className = "spinner";
  d.textContent = "Загрузка…";
  return d;
}

// ---------------------------------------------------------------- Router

const routes = {
  home: viewHome,
  auth: viewAuth,
  "play-ai": startAIGame,
  "create-mp": startMultiplayerGame,
  available: viewAvailable,
  join: viewJoin,
  leaderboard: viewLeaderboard,
  profile: viewProfile,
};

function navigate(name, arg) {
  if (name !== "auth" && !TOKENS.access) return viewAuth();
  (routes[name] || viewHome)(arg);
}

// Global click delegation for [data-nav] elements.
document.addEventListener("click", (e) => {
  const el = e.target.closest("[data-nav]");
  if (el) navigate(el.dataset.nav);
});

$("#btn-logout").addEventListener("click", () => {
  TOKENS.clear();
  ME = null;
  syncNav();
  viewAuth();
});

function syncNav() {
  $("#nav-auth").hidden = !TOKENS.access;
}

// ---------------------------------------------------------------- Auth view

function viewAuth() {
  const node = tpl("tpl-auth");
  render(node);
  let mode = "signin";

  const form = $("#auth-form");
  const submit = $("#auth-submit");
  app.querySelectorAll(".tab").forEach((t) =>
    t.addEventListener("click", () => {
      mode = t.dataset.tab;
      app.querySelectorAll(".tab").forEach((x) => x.classList.toggle("active", x === t));
      submit.textContent = mode === "signin" ? "Войти" : "Создать аккаунт";
    })
  );

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const fd = new FormData(form);
    const body = { login: fd.get("login").trim(), password: fd.get("password") };
    submit.disabled = true;
    try {
      // Signup only registers the account (no tokens returned), so sign in
      // afterwards with the same credentials to obtain access/refresh tokens.
      if (mode === "signup") {
        await api("/api/auth/signup", { method: "POST", auth: false, body });
      }
      const data = await api("/api/auth/signin", { method: "POST", auth: false, body });
      TOKENS.set(data);
      if (!TOKENS.access) throw new Error("Сервер не вернул токен доступа");
      syncNav();
      toast(mode === "signin" ? "С возвращением!" : "Аккаунт создан!", "ok");
      await loadMe();
      viewHome();
    } catch (err) {
      toast(err.message || "Ошибка авторизации", "error");
    } finally {
      submit.disabled = false;
    }
  });
}

async function loadMe() {
  try { ME = await api("/api/user/me"); } catch { /* ignore */ }
}

// ---------------------------------------------------------------- Home

function viewHome() {
  render(tpl("tpl-home"));
}

// ---------------------------------------------------------------- Games

async function startAIGame() {
  try {
    const r = await api("/api/game/ai", { method: "POST" });
    openGame(r.game_id, "ai");
  } catch (err) {
    toast(err.message, "error");
  }
}

async function startMultiplayerGame() {
  try {
    const r = await api("/api/game/multiplayer", { method: "POST" });
    openGame(r.game_id, "mp");
  } catch (err) {
    toast(err.message, "error");
  }
}

function viewJoin() {
  render(tpl("tpl-join"));
  $("#join-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const id = new FormData(e.target).get("gameID").trim();
    if (!id) return;
    try {
      await api(`/api/game/multiplayer/${id}`, { method: "POST" });
      openGame(id, "mp");
    } catch (err) {
      toast(err.message, "error");
    }
  });
}

// Open a game by id; mode is "ai" or "mp" (multiplayer).
async function openGame(gameID, mode) {
  if (!ME) await loadMe();
  render(tpl("tpl-game"));
  $("#game-id").textContent = gameID;
  $("#game-mode").textContent = mode === "ai" ? "🤖 vs ИИ" : "🌐 Мультиплеер";
  $("#copy-id").addEventListener("click", () => {
    navigator.clipboard?.writeText(gameID);
    toast("ID скопирован", "ok");
  });

  const state = { gameID, mode, game: null, busy: false };
  buildBoard(state);
  await refreshGame(state);

  if (mode === "mp") startPolling(state);
}

function buildBoard(state) {
  const board = $("#board");
  board.innerHTML = "";
  for (let r = 0; r < 3; r++) {
    for (let c = 0; c < 3; c++) {
      const cell = document.createElement("div");
      cell.className = "cell";
      cell.dataset.r = r;
      cell.dataset.c = c;
      cell.addEventListener("click", () => onCellClick(state, r, c));
      board.appendChild(cell);
    }
  }
}

async function refreshGame(state) {
  try {
    const game = await api(`/api/game/${state.gameID}`);
    state.game = game;
    drawGame(state);
  } catch (err) {
    toast(err.message, "error");
  }
}

// My icon: 1 (X) if I'm first player, 2 (O) if second. null if spectating.
function myIcon(game) {
  if (!ME) return null;
  if (game.firstPlayer_ID === ME.id) return 1;
  if (game.secondPlayer_ID === ME.id) return 2;
  return null;
}

const ICON_CHAR = { 1: "X", 2: "O" };

function drawGame(state) {
  const game = state.game;
  const board = $("#board");
  const cells = board.children;
  for (let r = 0; r < 3; r++) {
    for (let c = 0; c < 3; c++) {
      const v = game.board[r][c];
      const cell = cells[r * 3 + c];
      cell.textContent = v ? ICON_CHAR[v] : "";
      cell.className = "cell " + (v === 1 ? "x" : v === 2 ? "o" : "empty");
    }
  }

  const mine = myIcon(game);
  const finished = game.status === "draw" || game.status === "game_over";
  const myTurn = !finished && mine !== null && game.current_turn === ICON_CHAR[mine];

  board.classList.toggle("locked", !myTurn);

  const statusEl = $("#game-status");
  statusEl.className = "status";

  if (game.status === "waiting") {
    statusEl.textContent = "⏳ Ждём соперника… поделись ID игры";
  } else if (game.status === "draw") {
    statusEl.textContent = "🤝 Ничья!";
    statusEl.classList.add("draw");
  } else if (game.status === "game_over") {
    if (mine !== null) {
      const iWon = game.winner === ME.id;
      statusEl.textContent = iWon ? "🏆 Ты победил!" : "😞 Ты проиграл";
      statusEl.classList.add(iWon ? "win" : "lose");
    } else {
      statusEl.textContent = "Игра окончена";
    }
  } else { // playing
    if (mine === null) {
      statusEl.textContent = `Ходит: ${game.current_turn}`;
    } else {
      statusEl.textContent = myTurn ? "Твой ход" : "Ход соперника…";
    }
  }

  const actions = $("#game-actions");
  actions.innerHTML = "";
  if (finished) {
    const rematch = document.createElement("button");
    rematch.className = "btn primary";
    rematch.textContent = state.mode === "ai" ? "Сыграть снова" : "Новая игра";
    rematch.addEventListener("click", () =>
      state.mode === "ai" ? startAIGame() : startMultiplayerGame()
    );
    actions.appendChild(rematch);
    stopPolling();
  }
}

async function onCellClick(state, r, c) {
  const game = state.game;
  if (!game || state.busy) return;
  if (game.status !== "playing") return;
  if (game.board[r][c] !== 0) return;

  const mine = myIcon(game);
  if (mine === null || game.current_turn !== ICON_CHAR[mine]) {
    if (state.mode === "mp") toast("Сейчас не твой ход", "");
    return;
  }

  state.busy = true;
  $("#board").classList.add("locked");
  try {
    const updated = await api(`/api/game/${state.gameID}`, {
      method: "POST",
      body: { row: r, col: c },
    });
    state.game = updated;
    drawGame(state);
  } catch (err) {
    toast(err.message, "error");
  } finally {
    state.busy = false;
  }
}

function startPolling(state) {
  stopPolling();
  pollTimer = setInterval(async () => {
    if (state.busy) return;
    const g = state.game;
    if (g && (g.status === "draw" || g.status === "game_over")) {
      stopPolling();
      return;
    }
    // Only refresh while waiting on the opponent (or for them to join).
    const mine = myIcon(g || {});
    const theirTurn = g && (g.status === "waiting" ||
      (g.status === "playing" && (mine === null || g.current_turn !== ICON_CHAR[mine])));
    if (theirTurn) await refreshGame(state);
  }, 2000);
}

function stopPolling() {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null; }
}

// ---------------------------------------------------------------- Available games

async function viewAvailable() {
  render(tpl("tpl-list"));
  $("#list-title").textContent = "Открытые игры";
  $("#list-refresh").addEventListener("click", loadAvailable);
  loadAvailable();
}

async function loadAvailable() {
  const bodyEl = $("#list-body");
  bodyEl.innerHTML = "";
  bodyEl.appendChild(spinner());
  try {
    const data = await api("/api/game/available", { auth: false });
    const games = (data && data.games) || [];
    bodyEl.innerHTML = "";
    if (!games.length) {
      bodyEl.innerHTML = `<p class="empty-note">Нет открытых игр. Создай свою!</p>`;
      return;
    }
    for (const g of games) {
      const row = document.createElement("div");
      row.className = "list-row";
      row.innerHTML = `
        <div class="grow">
          <div>Игра <strong>${g.status}</strong></div>
          <code>${g.game_id}</code>
        </div>`;
      const btn = document.createElement("button");
      btn.className = "btn primary small";
      btn.textContent = "Войти";
      btn.addEventListener("click", async () => {
        try {
          await api(`/api/game/multiplayer/${g.game_id}`, { method: "POST" });
          openGame(g.game_id, "mp");
        } catch (err) { toast(err.message, "error"); }
      });
      row.appendChild(btn);
      bodyEl.appendChild(row);
    }
  } catch (err) {
    bodyEl.innerHTML = `<p class="empty-note">${err.message}</p>`;
  }
}

// ---------------------------------------------------------------- Leaderboard

async function viewLeaderboard() {
  render(tpl("tpl-list"));
  $("#list-title").textContent = "🏆 Лидеры";
  $("#list-refresh").addEventListener("click", loadLeaderboard);
  loadLeaderboard();
}

async function loadLeaderboard() {
  const bodyEl = $("#list-body");
  bodyEl.innerHTML = "";
  bodyEl.appendChild(spinner());
  try {
    const data = await api("/api/user/leaderboard/10");
    const rows = Array.isArray(data) ? data : [];
    bodyEl.innerHTML = "";
    if (!rows.length) {
      bodyEl.innerHTML = `<p class="empty-note">Пока нет данных</p>`;
      return;
    }
    rows.forEach((u, i) => {
      const medal = ["gold", "silver", "bronze"][i] || "";
      const row = document.createElement("div");
      row.className = "list-row";
      row.innerHTML = `
        <div class="rank ${medal}">${i + 1}</div>
        <div class="grow"><code>${u.id}</code></div>
        <strong>${u.winrate}%</strong>`;
      bodyEl.appendChild(row);
    });
  } catch (err) {
    bodyEl.innerHTML = `<p class="empty-note">${err.message}</p>`;
  }
}

// ---------------------------------------------------------------- Profile

async function viewProfile() {
  render(tpl("tpl-profile"));
  const bodyEl = $("#profile-body");
  try {
    const me = await api("/api/user/me");
    ME = me;
    const played = me.games_played || 0;
    const wins = me.wins || 0;
    const winrate = played ? Math.round((wins / played) * 100) : 0;
    bodyEl.innerHTML = `
      <div class="profile-name">${me.login}</div>
      <div class="stat-grid">
        <div class="stat"><div class="num">${played}</div><div class="lbl">Игр сыграно</div></div>
        <div class="stat"><div class="num">${wins}</div><div class="lbl">Побед</div></div>
        <div class="stat"><div class="num">${winrate}%</div><div class="lbl">Винрейт</div></div>
      </div>
      <code style="color:var(--muted);font-size:.8rem">${me.id}</code>`;
  } catch (err) {
    bodyEl.innerHTML = `<p class="empty-note">${err.message}</p>`;
  }
}

// ---------------------------------------------------------------- Boot

(async function boot() {
  syncNav();
  if (TOKENS.access) {
    await loadMe();
    viewHome();
  } else {
    viewAuth();
  }
})();
