// demo_entry: after a BG trigger the player picks one of two entries; a retrigger
// in FG-A can only enter FG-B and vice versa. The side panel explains each entry,
// shows the entry path of the bet and free-game progress; the trigger step marks
// the scatters.
(() => {
  const { kv, progress, esc } = DEV.util;
  const T = (en, zh) => DEV.t({ en, zh });
  const scatters = (ctx) => (ctx.act.acttype === "trigger" && ctx.act.ext ? ctx.act.ext.scatter_hits || [] : []);

  // Entry value -> free game; rounds are fixed, independent of the scatter count
  function entries(fixed) {
    return {
      3: { name: "FG-A", rounds: fixed.free_round_a, mult: 1 },
      7: { name: "FG-B", rounds: fixed.free_round_b, mult: fixed.multiplier_b },
    };
  }
  // The entry actually taken by segment i (i > 0)
  const entryOf = (bet, i) => bet.segments[i - 1].choose_next[bet.segments[i].choice];

  DEV.register(3, {
    title: { en: "demo_entry info", zh: "demo_entry 資訊" },
    modes: { 0: "BaseGame", 1: "FG-A", 2: "FG-B" },
    cellClass: (idx, ctx) => (scatters(ctx).includes(idx) ? ["mark"] : []),
    badge: (idx, ctx) => (scatters(ctx).includes(idx) ? "S" : ""),
    panel(ctx) {
      const table = entries(ctx.game.fixed);
      const name = (entry) => (table[entry] ? table[entry].name : T(`entry ${entry}`, `入口 ${entry}`));
      // Entry path: BG -> FG-A -> FG-B ..., the current segment in bold
      const path = ctx.bet.segments.map((_, i) => {
        const label = i === 0 ? "BG" : name(entryOf(ctx.bet, i));
        return i === ctx.segIdx ? `<b>${esc(label)}</b>` : esc(label);
      });
      let html = kv([
        [T("Entry path", "入口歷程"), path.join(" → ")],
        [T("Triggered", "觸發"), ctx.ext ? T("yes", "是") : T("not yet", "尚未")],
        ["Scatter", ctx.ext ? ctx.ext.scatters : "—"],
      ]);
      const seg = ctx.segment;
      if (seg.choose_next && seg.choose_next.length) {
        html += `<div class="muted">${T("Options", "選項")}</div>` + kv(seg.choose_next.map((entry, i) => {
          const e = table[entry];
          return [
            T(`#${i} entry ${entry}`, `#${i} 入口 ${entry}`),
            e ? T(`${e.name}: ${e.rounds} spins${e.mult > 1 ? `, win ×${e.mult}` : ""}`, `${e.name}：${e.rounds} 轉${e.mult > 1 ? `，贏分 ×${e.mult}` : ""}`) : T("unknown entry", "未知入口"),
          ];
        }));
      }
      if (ctx.segIdx > 0) {
        const e = table[entryOf(ctx.bet, ctx.segIdx)];
        if (e) html += progress(ctx.act.round + 1, e.rounds, T(`${e.name} spins${e.mult > 1 ? ` (win ×${e.mult})` : ""}`, `${e.name} 轉數${e.mult > 1 ? `（贏分 ×${e.mult}）` : ""}`));
      }
      html += `<div class="muted">${T("Retrigger: FG-A can only lead to FG-B, and FG-B only to FG-A.", "再觸發：FG-A 之後只能進 FG-B，FG-B 之後只能進 FG-A。")}</div>`;
      return html;
    },
  });
})();
