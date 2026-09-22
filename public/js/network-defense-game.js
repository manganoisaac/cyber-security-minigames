(() => {
  "use strict";

  /*
   * ------------------------------------------------------------------------
   * GAME CONFIGURATION
   * ------------------------------------------------------------------------
   * A grid-based tower defence. Hostile traffic spawns at the EDGE and
   * follows a fixed path to the CORE. Defenses are placed on buildable
   * nodes bordering the path and automatically fire at traffic within
   * range once a wave is launched.
   *
   * Each tower has its own `fireRateMs` cooldown (tracked per placed
   * instance in `state.towersByCell`) rather than firing on every tick -
   * wave/tower numbers are a first-pass balance for playtesting, not a
   * final tune. Range is measured in path-tiles from the tower's anchor
   * point rather than true grid distance, which keeps the targeting math
   * (and its screen-reader description) simple. `speed` on enemies is in
   * path-tiles per second (tick-resolution independent), so `tickMs` can
   * be re-tuned later without touching enemy balance.
   *
   * Each tower hard-counters exactly one enemy type (`counters`) and does
   * nothing at all against the other three - if a type has no matching
   * tower built, it walks through untouched. This is deliberate (mirrors
   * "different controls stop different attacks"), not a bug. `icon` is the
   * real sprite; `color` remains as a colored-disc backdrop behind it, so
   * distinction never rests on the (fairly similar-looking, once small)
   * sprite shape alone.
   */
  const GAME_CONFIG = {
    cols: 8,
    rows: 5,
    tickMs: 100,
    startingIntegrity: 100,
    towers: {
      firewall: {
        id: "firewall",
        name: "Firewall",
        cost: 80,
        damage: 18,
        fireRateMs: 400,
        range: 2,
        counters: "malwarePacket",
        description: "Filters common payloads at the perimeter.",
        icon: "/public/images/network-defense/towers/firewall.png",
      },
      ids: {
        id: "ids",
        name: "IDS Sensor",
        cost: 120,
        damage: 27,
        fireRateMs: 500,
        range: 3,
        counters: "bruteWorm",
        description: "Flags and intercepts suspicious signatures early.",
        icon: "/public/images/network-defense/towers/ids.png",
      },
      patched: {
        id: "patched",
        name: "Patched Endpoint",
        cost: 150,
        damage: 42,
        fireRateMs: 700,
        range: 2,
        counters: "phishPayload",
        description: "Hardened node that punishes malware at close range.",
        icon: "/public/images/network-defense/towers/patched.png",
      },
      mfa: {
        id: "mfa",
        name: "MFA Prompt",
        cost: 200,
        damage: 72,
        fireRateMs: 900,
        range: 2,
        counters: "credentialBot",
        description: "Second-factor challenge stops identity-based attacks cold.",
        icon: "/public/images/network-defense/towers/mfa.png",
      },
    },
    enemies: {
      credentialBot: {
        id: "credentialBot", name: "Credential-stuffing bot", hp: 70, speed: 1.43, coreDamage: 5, reward: 10,
        color: "#3a5f8a", icon: "/public/images/network-defense/enemies/credential-bot.png",
      },
      phishPayload: {
        id: "phishPayload", name: "Phishing payload", hp: 80, speed: 1.14, coreDamage: 8, reward: 16,
        color: "#6b3fa0", icon: "/public/images/network-defense/enemies/phish-payload.png",
      },
      malwarePacket: {
        id: "malwarePacket", name: "Malware packet", hp: 130, speed: 0.86, coreDamage: 12, reward: 24,
        color: "#a85a1e", icon: "/public/images/network-defense/enemies/malware-packet.png",
      },
      bruteWorm: {
        id: "bruteWorm", name: "Brute-force worm", hp: 400, speed: 0.71, coreDamage: 40, reward: 60,
        color: "#8a1f1f", icon: "/public/images/network-defense/enemies/brute-worm.png",
      },
    },
  };

  const WAVE_COUNT = 10;

  // Enemy types unlock progressively rather than all being available from
  // wave 1. Chosen in difficulty order (lowest hp first) rather than by
  // tower name, since introducing the 400 hp "boss" (bruteWorm, IDS
  // Sensor's counter) as early as wave 1 would be brutal before the player
  // can realistically afford IDS Sensor. Waves 7-10 are "everything".
  const TYPE_UNLOCKS = [
    { fromWave: 1, types: ["credentialBot", "phishPayload"] },
    { fromWave: 5, types: ["malwarePacket"] },
    { fromWave: 7, types: ["bruteWorm"] },
  ];

  function unlockedTypesForWave(waveNumber) {
    const unlocked = [];
    for (const tier of TYPE_UNLOCKS) {
      if (waveNumber >= tier.fromWave) unlocked.push(...tier.types);
    }
    return unlocked;
  }

  /*
   * ------------------------------------------------------------------------
   * SOLVABILITY GUARANTEES
   * ------------------------------------------------------------------------
   * Every tower hard-counters exactly one enemy type and does nothing to the
   * others, so a wave is only answerable if the player can own a counter for
   * every type in it. Playtest feedback was that this wasn't always true -
   * wave 1 alone needs both an MFA Prompt and a Patched Endpoint, which cost
   * more together than the opening balance used to be, making wave 1
   * impossible to fully defend no matter how well it was played.
   *
   * Rather than re-hardcoding a number that happens to work today, the helpers
   * below derive the requirement from the tower/unlock config itself, and it is
   * enforced at the three points where a player could otherwise be left without
   * an answer: the opening balance (STARTING_CREDITS), the between-wave payout
   * (completeWave), and the supply of buildable nodes (placeTower). Retuning a
   * tower cost or moving an unlock earlier can no longer silently reintroduce
   * an unwinnable wave.
   */

  // enemy type id -> the one tower that counters it. Derived from the towers'
  // own `counters` field so the two can never drift apart.
  const COUNTER_TOWER_BY_TYPE = {};
  for (const tower of Object.values(GAME_CONFIG.towers)) {
    COUNTER_TOWER_BY_TYPE[tower.counters] = tower.id;
  }

  // The counter towers a wave demands, de-duplicated. Pure - depends only on
  // config, so it is safe to call while building the opening balance below,
  // before any runtime state exists.
  function counterTowerIdsForWave(waveNumber) {
    const ids = [];
    for (const type of unlockedTypesForWave(waveNumber)) {
      const towerId = COUNTER_TOWER_BY_TYPE[type];
      if (towerId && !ids.includes(towerId)) ids.push(towerId);
    }
    return ids;
  }

  function totalCostOf(towerIds) {
    return towerIds.reduce((total, id) => total + GAME_CONFIG.towers[id].cost, 0);
  }

  // Wave 1 introduces two types at once, so the opening balance has to cover
  // a counter for both. Derived, never hardcoded - see the block comment.
  const STARTING_CREDITS = totalCostOf(counterTowerIdsForWave(1));

  function ownedTowerIds() {
    const owned = new Set();
    for (const tower of state.towersByCell.values()) owned.add(tower.id);
    return owned;
  }

  // Counters `waveNumber` demands that haven't been built yet.
  function unbuiltCounterIdsForWave(waveNumber) {
    const owned = ownedTowerIds();
    return counterTowerIdsForWave(waveNumber).filter((id) => !owned.has(id));
  }

  // Counters never built, across the whole run. Towers can't be sold, so each
  // of these needs a buildable node kept free for it (see placeTower).
  function missingCounterTowerIds() {
    return unbuiltCounterIdsForWave(WAVE_COUNT);
  }

  function waveLabelFor(waveNumber) {
    if (waveNumber >= 7) return "Wave " + waveNumber + " - Full assault";
    if (waveNumber >= 5) return "Wave " + waveNumber + " - Escalation";
    return "Wave " + waveNumber + " - Probing";
  }

  // Randomly composed from whatever's unlocked at this wave number, with
  // enemy count scaling up over the 10 waves. bruteWorm is capped at one
  // per wave - it's meant as an occasional gut-check, not a swarm.
  function generateWave(waveNumber) {
    const pool = unlockedTypesForWave(waveNumber);
    const totalEnemies = Math.min(14, 4 + waveNumber);
    const spawns = [];
    let elapsed = 0;
    for (let i = 0; i < totalEnemies; i++) {
      let type = pool[Math.floor(Math.random() * pool.length)];
      if (type === "bruteWorm" && spawns.some((s) => s.type === "bruteWorm")) {
        const rest = pool.filter((t) => t !== "bruteWorm");
        type = rest[Math.floor(Math.random() * rest.length)];
      }
      spawns.push({ type, delay: Math.round(elapsed) });
      elapsed += 700 + Math.random() * 700;
    }
    return { label: waveLabelFor(waveNumber), spawns, pool };
  }

  /*
   * ------------------------------------------------------------------------
   * PATH + GRID SETUP (pure data, no DOM)
   * ------------------------------------------------------------------------
   */
  function buildPath() {
    const pts = [];
    const push = (c, r) => pts.push([c, r]);
    for (let c = 0; c <= 7; c++) push(c, 0); // top row, left to right
    push(7, 1);
    push(7, 2); // down the right side
    for (let c = 6; c >= 0; c--) push(c, 2); // middle row, right to left
    push(0, 3);
    push(0, 4); // down the left side
    for (let c = 1; c <= 7; c++) push(c, 4); // bottom row, left to right, to CORE
    return pts;
  }

  const PATH = buildPath();
  const cellKey = (c, r) => c + "," + r;
  const EDGE_KEY = cellKey(PATH[0][0], PATH[0][1]);
  const CORE_KEY = cellKey(PATH[PATH.length - 1][0], PATH[PATH.length - 1][1]);

  const pathIndexByKey = new Map();
  PATH.forEach(([c, r], i) => pathIndexByKey.set(cellKey(c, r), i));

  // Buildable = any non-path cell orthogonally adjacent to at least one path
  // cell. Range is measured in path-tiles from the cell's "anchors" rather
  // than as true 2D distance, which keeps the targeting math (and its
  // screen-reader description) simple.
  //
  // Every adjacent path index is kept, not just the lowest one. The path
  // doubles back on itself, so a node between two stretches of it - say the
  // gap between the middle row and the bottom row - genuinely borders both,
  // and a defense there should reach both. Keeping only the lowest index
  // silently orphaned the entire bottom row (7 of 26 tiles, the whole final
  // approach to the CORE), because those cells always anchored to the
  // middle-row stretch instead. Nothing could be built that reached the last
  // leg of the path, so anything still alive by then was a guaranteed breach.
  function computeBuildable() {
    const buildable = new Map();
    for (let r = 0; r < GAME_CONFIG.rows; r++) {
      for (let c = 0; c < GAME_CONFIG.cols; c++) {
        const key = cellKey(c, r);
        if (pathIndexByKey.has(key)) continue;
        const neighbours = [[c - 1, r], [c + 1, r], [c, r - 1], [c, r + 1]];
        const anchors = [];
        for (const [nc, nr] of neighbours) {
          const nIdx = pathIndexByKey.get(cellKey(nc, nr));
          if (nIdx !== undefined) anchors.push(nIdx);
        }
        if (anchors.length) buildable.set(key, anchors.sort((a, b) => a - b));
      }
    }
    return buildable;
  }

  const BUILDABLE = computeBuildable();

  function positionForProgress(pathPos) {
    const clamped = Math.max(0, Math.min(PATH.length - 1, pathPos));
    const i0 = Math.floor(clamped);
    const i1 = Math.min(PATH.length - 1, i0 + 1);
    const frac = clamped - i0;
    const [c0, r0] = PATH[i0];
    const [c1, r1] = PATH[i1];
    const c = c0 + (c1 - c0) * frac;
    const r = r0 + (r1 - r0) * frac;
    return {
      leftPct: ((c + 0.5) / GAME_CONFIG.cols) * 100,
      topPct: ((r + 0.5) / GAME_CONFIG.rows) * 100,
    };
  }

  /*
   * ------------------------------------------------------------------------
   * RUNTIME STATE
   * ------------------------------------------------------------------------
   */
  const state = {
    credits: STARTING_CREDITS,
    integrity: GAME_CONFIG.startingIntegrity,
    waveIndex: 0, // number of waves launched so far (0 = none yet)
    currentWaveLabel: "",
    status: "build", // build | wave-active | won | lost
    towersByCell: new Map(),
    enemies: [],
    selectedTower: null,
    pendingSpawns: [],
    waveElapsedMs: 0,
    enemySeq: 0,
    currentWavePool: [], // enemy type ids that can appear in the active/just-finished wave
    breachedTypes: new Set(), // enemy type ids that reached the core this wave
  };

  let tickTimer = null;

  /*
   * ------------------------------------------------------------------------
   * DOM REFERENCES (filled in on init)
   * ------------------------------------------------------------------------
   */
  let gridEl, waveLabelEl, creditsLabelEl, launchBtn, selectedInfoEl, upcomingPreviewEl;
  let integrityTextEl, integrityBarEl, threatTextEl, threatBarEl, eventLogEl;
  let introModal, introContinueBtn, introSkipEl, instructionsModal, beginBtn;
  let resultModal, resultTitleEl, resultMessageEl, resultBtn;
  let waveSummaryModal, waveSummaryListEl, waveSummaryContinueBtn;

  const towerButtons = new Map(); // towerId -> button element
  const cellEls = new Map(); // cellKey -> element (button for buildable, div otherwise)
  const enemyEls = new Map(); // enemy.uid -> element

  const MAX_LOG_ENTRIES = 40;

  function logEvent(text) {
    const li = document.createElement("li");
    li.textContent = text;
    eventLogEl.appendChild(li);
    while (eventLogEl.children.length > MAX_LOG_ENTRIES) {
      eventLogEl.removeChild(eventLogEl.firstChild);
    }
  }

  /*
   * ------------------------------------------------------------------------
   * SIMULATION
   * ------------------------------------------------------------------------
   */
  // Targets the enemy furthest along the path, so fire is focused on whatever
  // is closest to breaching rather than spread across the whole queue.
  function pickTarget(anchors, range, countersType) {
    let best = null;
    for (const enemy of state.enemies) {
      if (enemy.type !== countersType) continue;
      const tile = Math.floor(enemy.pathPos);
      if (!anchors.some((anchor) => Math.abs(tile - anchor) <= range)) continue;
      if (!best || enemy.pathPos > best.pathPos) best = enemy;
    }
    return best;
  }

  function spawnEnemy(type) {
    state.enemySeq += 1;
    const def = GAME_CONFIG.enemies[type];
    state.enemies.push({ uid: state.enemySeq, type, pathPos: 0, hp: def.hp });
    logEvent(def.name + " detected at the edge.");
  }

  function applyCoreDamage(def) {
    state.integrity = Math.max(0, state.integrity - def.coreDamage);
    state.breachedTypes.add(def.id);
    logEvent(def.name + " reached the core. Integrity -" + def.coreDamage + ".");
  }

  function tick() {
    state.waveElapsedMs += GAME_CONFIG.tickMs;

    for (const enemy of state.enemies) {
      enemy.pathPos += GAME_CONFIG.enemies[enemy.type].speed * (GAME_CONFIG.tickMs / 1000);
    }

    const arrived = state.enemies.filter((e) => e.pathPos >= PATH.length - 1);
    for (const enemy of arrived) applyCoreDamage(GAME_CONFIG.enemies[enemy.type]);
    state.enemies = state.enemies.filter((e) => e.pathPos < PATH.length - 1);

    for (const [key, tower] of state.towersByCell) {
      // Decrement before testing. Testing first burned a whole extra tick
      // between shots (the tick that took the cooldown to exactly 0 couldn't
      // also fire), so every tower fired at fireRateMs + tickMs - a Firewall
      // configured for 400ms was really firing every 500ms, 25% under its
      // advertised rate.
      if (tower.cooldown > 0) tower.cooldown -= GAME_CONFIG.tickMs;
      if (tower.cooldown > 0) continue;
      const def = GAME_CONFIG.towers[tower.id];
      const target = pickTarget(BUILDABLE.get(key), def.range, def.counters);
      if (target) {
        target.hp -= def.damage;
        spawnProjectile(key, target);
        tower.cooldown = def.fireRateMs;
      }
    }

    const dead = state.enemies.filter((e) => e.hp <= 0);
    for (const enemy of dead) {
      const def = GAME_CONFIG.enemies[enemy.type];
      state.credits += def.reward;
      logEvent(def.name + " neutralised (+" + def.reward + " credits).");
    }
    state.enemies = state.enemies.filter((e) => e.hp > 0);

    // Spawns happen last, after this tick's movement and combat have already
    // resolved - a newly spawned enemy always gets at least one render at
    // pathPos 0, untouched, before it can move or be targeted. Without this
    // ordering, a tower whose damage exceeds its target's max hp (MFA vs a
    // credential-stuffing bot) could spawn-move-and-die within a single
    // tick, before render() ever ran - which is exactly what looked like
    // MFA "nuking minions before they appear".
    while (state.pendingSpawns.length && state.pendingSpawns[0].delay <= state.waveElapsedMs) {
      spawnEnemy(state.pendingSpawns.shift().type);
    }

    if (state.integrity <= 0) {
      endGame(false);
      render();
      return;
    }

    if (state.status === "wave-active" && state.pendingSpawns.length === 0 && state.enemies.length === 0) {
      completeWave();
    }

    render();
  }

  function startTickTimer() {
    stopTickTimer();
    tickTimer = setInterval(tick, GAME_CONFIG.tickMs);
  }

  function stopTickTimer() {
    if (tickTimer) {
      clearInterval(tickTimer);
      tickTimer = null;
    }
  }

  function launchWave() {
    if (state.status !== "build") return;
    state.waveIndex += 1;
    const wave = generateWave(state.waveIndex);
    state.currentWaveLabel = wave.label;
    state.pendingSpawns = wave.spawns.slice();
    state.currentWavePool = wave.pool;
    state.breachedTypes = new Set();
    state.waveElapsedMs = 0;
    state.status = "wave-active";
    logEvent(wave.label + " incoming.");
    startTickTimer();
    render();
  }

  function completeWave() {
    stopTickTimer();
    logEvent(state.currentWaveLabel + " cleared.");

    // Guaranteed per-round credit floor, independent of kill count - a wave
    // "completing" only requires pendingSpawns and enemies to both be empty,
    // which happens whether the player killed everything or let it all
    // through. Without this, building the wrong counter-tower(s) for a
    // wave's composition can permanently softlock a player at 0 income.
    const roundBonus = 50 + 10 * state.waveIndex;
    state.credits += roundBonus;
    logEvent("Round bonus: +" + roundBonus + " credits.");

    if (state.waveIndex >= WAVE_COUNT) {
      endGame(true);
    } else {
      // Hard floor on the payout: whatever else has been spent, the player can
      // always afford a counter for every type the next wave can field. In
      // normal play the round bonus already covers this and nothing is added -
      // it only ever fires for someone who has spent themselves into a corner,
      // which is precisely the case that used to be unrecoverable.
      const needed = totalCostOf(unbuiltCounterIdsForWave(state.waveIndex + 1));
      const shortfall = needed - state.credits;
      if (shortfall > 0) {
        state.credits += shortfall;
        logEvent("Emergency budget released: +" + shortfall + " credits.");
      }
      state.status = "build";
      // Skipped on the final wave - the win modal from endGame() above is
      // the more important thing to show, and there's no "next wave" to
      // prepare for.
      showWaveSummary();
    }
  }

  // Recaps, between waves, what just came through and what stops it - the
  // "Threat reference" panel has this too, but it's easy to never look at it;
  // surfacing it right when a wave ends is when a player is most likely to
  // actually read it. Anything that reached the core gets called out by name
  // with the defense that should have stopped it, addressing the specific
  // "why did that get through" moment rather than just the general reference.
  function showWaveSummary() {
    const items = state.currentWavePool.map((typeId) => {
      const def = GAME_CONFIG.enemies[typeId];
      const counterTower = Object.values(GAME_CONFIG.towers).find((t) => t.counters === typeId);
      const counterName = counterTower ? counterTower.name : "nothing yet";
      const breached = state.breachedTypes.has(typeId);
      return (
        '<li class="network-defense-wave-summary__item">' +
        '<span class="network-defense-threat-legend__swatch" style="background-color:' +
        def.color +
        ";background-image:url(" +
        def.icon +
        ')"></span>' +
        "<span><strong>" +
        def.name +
        "</strong> - stopped by " +
        counterName +
        (breached
          ? '<br/><em class="network-defense-wave-summary__breach">Reached the core this wave - deploy ' +
            counterName +
            " to stop it next time.</em>"
          : "") +
        "</span></li>"
      );
    });
    waveSummaryListEl.innerHTML = items.join("");
    waveSummaryModal.hidden = false;
  }

  function endGame(won) {
    stopTickTimer();
    state.status = won ? "won" : "lost";
    resultTitleEl.textContent = won ? "Perimeter held" : "Core breached";
    resultMessageEl.textContent = won
      ? "All waves neutralised. Metacortex's network held the line."
      : "Network integrity reached zero. The intrusion broke through.";
    resultBtn.textContent = won ? "Continue" : "Try again";
    resultModal.hidden = false;
  }

  function selectTower(towerId) {
    state.selectedTower = state.selectedTower === towerId ? null : towerId;
    render();
  }

  function placeTower(key) {
    if (!state.selectedTower) return;
    if (state.towersByCell.has(key)) return;
    const def = GAME_CONFIG.towers[state.selectedTower];
    if (state.credits < def.cost) return;

    // Guaranteed credits are worth nothing if there is nowhere left to spend
    // them. Towers can't be sold, so filling every node with duplicates could
    // otherwise strand a player with the money for a counter they still need
    // and no node to put it on. Keep one node in reserve for each counter not
    // yet built; deploying a missing counter itself is never blocked.
    const stillMissingAfter = missingCounterTowerIds().filter((id) => id !== def.id);
    const freeCellsAfter = BUILDABLE.size - state.towersByCell.size - 1;
    if (freeCellsAfter < stillMissingAfter.length) {
      const names = stillMissingAfter.map((id) => GAME_CONFIG.towers[id].name).join(", ");
      logEvent("Node reserved - you still need room to deploy: " + names + ".");
      return;
    }

    state.credits -= def.cost;
    state.towersByCell.set(key, { id: state.selectedTower, cooldown: 0 });
    logEvent(def.name + " deployed.");
    state.selectedTower = null;
    render();
  }

  /*
   * ------------------------------------------------------------------------
   * RENDERING
   * ------------------------------------------------------------------------
   */
  function renderHUD() {
    waveLabelEl.textContent = "Wave " + state.waveIndex + " / " + WAVE_COUNT;
    creditsLabelEl.textContent = "Credits: " + state.credits;

    const pct = Math.max(0, Math.round(state.integrity));
    integrityTextEl.textContent = pct + "%";
    integrityBarEl.style.width = pct + "%";
    integrityBarEl.style.backgroundColor =
      pct <= 25 ? "var(--nd-danger)" : pct <= 50 ? "var(--nd-warning)" : "var(--nd-success)";

    const threatPct = Math.min(100, (state.waveIndex / WAVE_COUNT) * 100);
    threatBarEl.style.width = threatPct + "%";
    threatTextEl.textContent =
      state.waveIndex === 0
        ? "Standing by"
        : state.waveIndex >= WAVE_COUNT
        ? "Critical"
        : state.waveIndex >= Math.floor(WAVE_COUNT / 2)
        ? "Elevated"
        : "Active";

    launchBtn.disabled = state.status !== "build";
    launchBtn.textContent =
      state.status === "wave-active"
        ? "Wave in progress..."
        : state.waveIndex >= WAVE_COUNT
        ? "All waves cleared"
        : "Launch wave";

    if (state.status === "build" && state.waveIndex < WAVE_COUNT) {
      const nextPool = unlockedTypesForWave(state.waveIndex + 1);
      upcomingPreviewEl.textContent =
        "Next wave may include: " + nextPool.map((t) => GAME_CONFIG.enemies[t].name).join(", ");
    } else {
      upcomingPreviewEl.textContent = "";
    }

    if (!state.selectedTower) {
      selectedInfoEl.textContent = "Nothing selected. Choose a defense above.";
    } else {
      const def = GAME_CONFIG.towers[state.selectedTower];
      selectedInfoEl.textContent =
        "Selected: " + def.name + ". Click a highlighted node to deploy (" + def.cost + " credits).";
    }
  }

  function renderTowerList() {
    for (const [id, btn] of towerButtons) {
      const def = GAME_CONFIG.towers[id];
      btn.disabled = state.credits < def.cost;
      btn.setAttribute("aria-pressed", state.selectedTower === id ? "true" : "false");
    }
  }

  function renderGrid() {
    for (const [key, el] of cellEls) {
      const isPath = pathIndexByKey.has(key);
      const isEdge = key === EDGE_KEY;
      const isCore = key === CORE_KEY;
      const tower = state.towersByCell.get(key);
      const towerId = tower && tower.id;

      let className = "network-defense-grid__cell";
      if (isPath) className += " network-defense-grid__cell--path";
      else className += " network-defense-grid__cell--buildable";
      if (isEdge) className += " network-defense-grid__cell--edge";
      if (isCore) className += " network-defense-grid__cell--core";
      if (towerId) className += " network-defense-grid__cell--placed";
      el.className = className;

      if (isPath) {
        el.textContent = isEdge ? "EDGE" : isCore ? "CORE" : "";
      } else if (towerId) {
        el.textContent = "";
        el.style.backgroundImage = "url(" + GAME_CONFIG.towers[towerId].icon + ")";
        el.disabled = true;
        el.setAttribute(
          "aria-label",
          GAME_CONFIG.towers[towerId].name +
            " deployed. Damage " +
            GAME_CONFIG.towers[towerId].damage +
            ", range " +
            GAME_CONFIG.towers[towerId].range +
            "."
        );
      } else {
        el.textContent = "";
        el.style.backgroundImage = "";
        el.disabled = false;
        el.setAttribute("aria-label", "Buildable node. Empty. Select a defense, then activate to deploy it here.");
      }
    }
  }

  function spawnProjectile(fromKey, targetEnemy) {
    const [fc, fr] = fromKey.split(",").map(Number);
    const fromLeft = ((fc + 0.5) / GAME_CONFIG.cols) * 100;
    const fromTop = ((fr + 0.5) / GAME_CONFIG.rows) * 100;
    const toPos = positionForProgress(targetEnemy.pathPos);
    const enemyDef = GAME_CONFIG.enemies[targetEnemy.type];

    const dot = document.createElement("div");
    dot.className = "network-defense-projectile";
    dot.style.backgroundColor = enemyDef.color;
    dot.style.left = fromLeft + "%";
    dot.style.top = fromTop + "%";
    dot.setAttribute("aria-hidden", "true");
    gridEl.appendChild(dot);

    requestAnimationFrame(() => {
      dot.style.left = toPos.leftPct + "%";
      dot.style.top = toPos.topPct + "%";
    });
    setTimeout(() => dot.remove(), 320);
  }

  function renderEnemies() {
    const seen = new Set();
    for (const enemy of state.enemies) {
      seen.add(enemy.uid);
      const def = GAME_CONFIG.enemies[enemy.type];
      let el = enemyEls.get(enemy.uid);
      if (!el) {
        el = document.createElement("div");
        el.className = "network-defense-enemy";
        el.setAttribute("aria-hidden", "true");
        el.innerHTML =
          '<span class="network-defense-enemy__hpbar"><span class="network-defense-enemy__hpfill"></span></span>' +
          '<span class="network-defense-enemy__body"></span>';
        const bodyEl = el.querySelector(".network-defense-enemy__body");
        bodyEl.style.backgroundColor = def.color;
        bodyEl.style.backgroundImage = "url(" + def.icon + ")";
        gridEl.appendChild(el);
        enemyEls.set(enemy.uid, el);
      }
      const pos = positionForProgress(enemy.pathPos);
      el.style.left = pos.leftPct + "%";
      el.style.top = pos.topPct + "%";
      const hpPct = Math.max(0, Math.min(100, (enemy.hp / def.hp) * 100));
      el.querySelector(".network-defense-enemy__hpfill").style.width = hpPct + "%";
      el.title = def.name + " (" + Math.max(0, Math.round(enemy.hp)) + "/" + def.hp + " hp)";
    }
    for (const [uid, el] of enemyEls) {
      if (!seen.has(uid)) {
        el.remove();
        enemyEls.delete(uid);
      }
    }
  }

  function render() {
    renderHUD();
    renderTowerList();
    renderGrid();
    renderEnemies();
  }

  /*
   * ------------------------------------------------------------------------
   * INTRO TYPEWRITER
   * ------------------------------------------------------------------------
   * Matches the shared modal typewriter (public/js/modal-typewriter.js) - same
   * 20ms cadence and blinking caret - but that script drives a single
   * `.modal__text` inside a `.modal`, and this modal is three paragraphs in
   * this game's own markup. Giving the overlay a `.modal` class to reuse it
   * would drag in that component's positioning and dialog styles on top of
   * this game's, so the cadence is shared by constant rather than by code.
   */
  const INTRO_TYPE_SPEED_MS = 20;
  const INTRO_TYPING_CLASS = "network-defense-intro__line--typing";

  let introTimer = null;
  let introFinished = false;
  let introSkipKeyHandler = null;

  // Releases the interval and the document-level key listener. Called both when
  // the text finishes and when a boosted swap takes the page away mid-type -
  // otherwise either would outlive the page, against a detached modal.
  function stopIntroTyping() {
    if (introTimer) {
      clearInterval(introTimer);
      introTimer = null;
    }
    if (introSkipKeyHandler) {
      document.removeEventListener("keydown", introSkipKeyHandler);
      introSkipKeyHandler = null;
    }
  }

  // Fills in everything at once and hands control back to the player. Safe to
  // call twice - skipping and finishing naturally both land here.
  function finishIntroTyping(lines) {
    if (introFinished) return;
    introFinished = true;
    stopIntroTyping();
    for (const line of lines) {
      line.textContent = line.dataset.fullText;
      line.classList.remove(INTRO_TYPING_CLASS);
    }
    introContinueBtn.disabled = false;
    if (introSkipEl) introSkipEl.hidden = true;
  }

  function startIntroTyping() {
    const lines = Array.from(introModal.querySelectorAll(".network-defense-intro__line"));
    if (!lines.length) return;

    // Collapse the source indentation in the template down to single spaces so
    // the typed-out text doesn't stutter through runs of whitespace.
    for (const line of lines) {
      line.dataset.fullText = line.textContent.replace(/\s+/g, " ").trim();
    }

    // Respect a reduced-motion preference, and don't trap a screen reader
    // behind three paragraphs being spelled out one character at a time.
    const reduceMotion =
      window.matchMedia && window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    if (reduceMotion) {
      finishIntroTyping(lines);
      return;
    }

    for (const line of lines) line.textContent = "";
    introContinueBtn.disabled = true;
    if (introSkipEl) introSkipEl.hidden = false;

    let lineIndex = 0;
    let charIndex = 0;
    lines[0].classList.add(INTRO_TYPING_CLASS);

    introTimer = setInterval(() => {
      const line = lines[lineIndex];
      const full = line.dataset.fullText;
      if (charIndex < full.length) {
        charIndex += 1;
        line.textContent = full.slice(0, charIndex);
        return;
      }
      line.classList.remove(INTRO_TYPING_CLASS);
      lineIndex += 1;
      charIndex = 0;
      if (lineIndex >= lines.length) {
        finishIntroTyping(lines);
        return;
      }
      lines[lineIndex].classList.add(INTRO_TYPING_CLASS);
    }, INTRO_TYPE_SPEED_MS);

    // Anywhere on the modal, or any of the obvious "get on with it" keys.
    introModal.addEventListener("click", () => finishIntroTyping(lines));
    introSkipKeyHandler = (e) => {
      if (e.key === "Enter" || e.key === " " || e.key === "Escape") {
        e.preventDefault();
        finishIntroTyping(lines);
      }
    };
    document.addEventListener("keydown", introSkipKeyHandler);
  }

  /*
   * ------------------------------------------------------------------------
   * INITIALISATION
   * ------------------------------------------------------------------------
   */
  function buildTowerList() {
    Object.values(GAME_CONFIG.towers).forEach((def) => {
      const btn = document.createElement("button");
      btn.type = "button";
      btn.className = "network-defense-tower-list__item";
      btn.setAttribute("aria-pressed", "false");
      const countered = GAME_CONFIG.enemies[def.counters];
      btn.innerHTML =
        '<span class="network-defense-tower-list__top-row"><strong>' +
        def.name +
        '</strong><span class="network-defense-tower-list__cost">' +
        def.cost +
        ' cr</span></span>' +
        '<span class="network-defense-tower-list__stats">Damage ' +
        def.damage +
        " &middot; Range " +
        def.range +
        " &middot; Fires every " +
        (def.fireRateMs / 1000).toFixed(1) +
        "s</span>" +
        '<span class="network-defense-tower-list__counters">Effective against: ' +
        countered.name +
        " only</span>" +
        '<span class="network-defense-tower-list__desc">' +
        def.description +
        "</span>";
      btn.addEventListener("click", () => selectTower(def.id));
      document.getElementById("nd-tower-list").appendChild(btn);
      towerButtons.set(def.id, btn);
    });
  }

  function buildThreatLegend() {
    const list = document.getElementById("nd-threat-legend");
    Object.values(GAME_CONFIG.enemies).forEach((def) => {
      const counterTower = Object.values(GAME_CONFIG.towers).find((t) => t.counters === def.id);
      const li = document.createElement("li");
      li.className = "network-defense-threat-legend__item";
      li.innerHTML =
        '<span class="network-defense-threat-legend__swatch" style="background-color:' +
        def.color +
        ";background-image:url(" +
        def.icon +
        ')"></span>' +
        "<span><strong>" +
        def.name +
        "</strong> - stopped by " +
        (counterTower ? counterTower.name : "nothing yet") +
        "</span>";
      list.appendChild(li);
    });
  }

  function buildGrid() {
    for (let r = 0; r < GAME_CONFIG.rows; r++) {
      for (let c = 0; c < GAME_CONFIG.cols; c++) {
        const key = cellKey(c, r);
        const isPath = pathIndexByKey.has(key);
        const el = document.createElement(isPath ? "div" : "button");
        if (!isPath) {
          el.type = "button";
          el.addEventListener("click", () => placeTower(key));
        } else if (key === EDGE_KEY || key === CORE_KEY) {
          el.tabIndex = 0;
          el.setAttribute(
            "aria-label",
            key === EDGE_KEY
              ? "Edge. Intrusion traffic enters the network here."
              : "Core. Protect this asset - network integrity depends on it."
          );
        } else {
          el.setAttribute("aria-hidden", "true");
        }
        el.style.gridColumn = String(c + 1);
        el.style.gridRow = String(r + 1);
        gridEl.appendChild(el);
        cellEls.set(key, el);
      }
    }
  }

  function init() {
    gridEl = document.getElementById("nd-grid");
    waveLabelEl = document.getElementById("nd-wave-label");
    creditsLabelEl = document.getElementById("nd-credits-label");
    launchBtn = document.getElementById("nd-launch-wave");
    selectedInfoEl = document.getElementById("nd-selected-info");
    upcomingPreviewEl = document.getElementById("nd-upcoming-preview");
    integrityTextEl = document.getElementById("nd-integrity-text");
    integrityBarEl = document.getElementById("nd-integrity-bar");
    threatTextEl = document.getElementById("nd-threat-text");
    threatBarEl = document.getElementById("nd-threat-bar");
    eventLogEl = document.getElementById("nd-event-log");
    introModal = document.getElementById("nd-intro-modal");
    introContinueBtn = document.getElementById("nd-intro-continue");
    introSkipEl = document.getElementById("nd-intro-skip");
    instructionsModal = document.getElementById("nd-instructions-modal");
    beginBtn = document.getElementById("nd-begin-button");
    resultModal = document.getElementById("nd-result-modal");
    resultTitleEl = document.getElementById("nd-result-title");
    resultMessageEl = document.getElementById("nd-result-message");
    resultBtn = document.getElementById("nd-result-button");
    waveSummaryModal = document.getElementById("nd-wave-summary-modal");
    waveSummaryListEl = document.getElementById("nd-wave-summary-list");
    waveSummaryContinueBtn = document.getElementById("nd-wave-summary-continue");

    buildTowerList();
    buildThreatLegend();
    buildGrid();
    startIntroTyping();

    launchBtn.addEventListener("click", launchWave);
    introContinueBtn.addEventListener("click", () => {
      introModal.hidden = true;
      instructionsModal.hidden = false;
    });
    beginBtn.addEventListener("click", () => {
      instructionsModal.hidden = true;
      logEvent("Operation started. Deploy defenses, then launch the first wave.");
    });
    resultBtn.addEventListener("click", () => {
      // A win continues the story into this game's "What have we learned" page
      // (the route is /whwl/network-defense - American spelling, matching this
      // game's own route, even though the templ and component say "defence").
      // Only a loss's "Try again" replays the game.
      window.location.href =
        state.status === "won" ? "/whwl/network-defense" : "/network-defense";
    });
    waveSummaryContinueBtn.addEventListener("click", () => {
      waveSummaryModal.hidden = true;
    });

    render();
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init, { once: true });
  } else {
    // The layout boosts internal navigation (AGENTS.md's hx-boost
    // convention), so arriving here from chapter-select swaps this page's
    // markup into the existing document instead of doing a real navigation.
    // htmx re-clones and re-executes this <script src> tag on that swap
    // (rather than skip it, the way a plain innerHTML assignment would),
    // but DOMContentLoaded itself already fired once, for the very first
    // real page load - it will never fire again for this second execution,
    // which is exactly what left the page inert until a manual refresh.
    // By the time a re-run gets here the document has long since finished
    // loading, so init directly instead of waiting on an event that's
    // already in the past.
    init();
  }

  // A boosted swap never unloads the document, so leaving this page mid-wave
  // would otherwise leave this execution's tick timer running in the
  // background forever, against a grid that's no longer in the document.
  // Stop it right before the next swap replaces the page.
  //
  // `once` matters: htmx re-executes this script on every boosted swap back
  // into the page, and without it each visit would leave another listener
  // bound to a dead closure on a document that never unloads. One swap away
  // is the end of this execution's life, so the listener retires with it.
  document.addEventListener(
    "htmx:beforeSwap",
    () => {
      stopTickTimer();
      stopIntroTyping();
    },
    { once: true }
  );
})();
