// demo_choice: after a BG trigger the player picks one of two entries. The side
// panel explains each entry, derives the free-game rounds from the checkpoint's
// scatter count and shows free-game progress; the trigger step marks the scatters.
(() => {
  const { kv, progress, esc } = DEV.util;
  const T = (en, zh) => DEV.t({ en, zh });
  const scatters = (ctx) => (ctx.act.acttype === "trigger" && ctx.act.ext ? ctx.act.ext.scatter_hits || [] : []);

  // Entry value -> free game; rounds = base rounds + (scatters - 3)
  function entries(fixed) {
    return {
      1: { name: "FG-A", base: fixed.free_round_a, mult: 1 },
      2: { name: "FG-B", base: fixed.free_round_b, mult: fixed.multiplier_b },
    };
  }
  const rounds = (e, scatterCount) => e.base + scatterCount - 3;

  DEV.register(2, {
    title: { en: "demo_choice info", zh: "demo_choice 資訊" },
    modes: { 0: "BaseGame", 1: "FG-A", 2: "FG-B" },
    cellClass: (idx, ctx) => (scatters(ctx).includes(idx) ? ["mark"] : []),
    badge: (idx, ctx) => (scatters(ctx).includes(idx) ? "S" : ""),
    panel(ctx) {
      const table = entries(ctx.game.fixed);
      const seg = ctx.segment;
      let html = "";
      if (ctx.segIdx === 0) {
        const cp = seg.spin_state && seg.spin_state.cp;
        html += kv([[T("Triggered", "觸發"), ctx.ext ? T("yes", "是") : T("not yet", "尚未")], ["Scatter", ctx.ext ? ctx.ext.scatters : "—"]]);
        if (cp && seg.choose_next) {
          html += `<div class="muted">${T(`Options (${cp.scatters} Scatters)`, `選項（Scatter ${cp.scatters} 個）`)}</div>` + kv(seg.choose_next.map((entry, i) => {
            const e = table[entry];
            const n = e && rounds(e, cp.scatters);
            return [
              T(`#${i} entry ${entry}`, `#${i} 入口 ${entry}`),
              e ? T(`${e.name}: ${n} spins${e.mult > 1 ? `, win ×${e.mult}` : ""}`, `${e.name}：${n} 轉${e.mult > 1 ? `，贏分 ×${e.mult}` : ""}`) : T("unknown entry", "未知入口"),
            ];
          }));
        }
      } else {
        const prev = ctx.bet.segments[ctx.segIdx - 1];
        const entry = prev.choose_next[seg.choice];
        const e = table[entry];
        const scatterCount = prev.spin_state.cp ? prev.spin_state.cp.scatters : 3;
        const total = e ? rounds(e, scatterCount) : 0;
        html += kv([
          [T("Choice", "選擇"), T(`#${seg.choice} entry ${esc(entry)} (${e ? e.name : "unknown"})`, `#${seg.choice} 入口 ${esc(entry)}（${e ? e.name : "未知"}）`)],
          [T("Spins", "轉數推算"), e ? T(`${e.base} + (Scatter ${scatterCount} − 3) = ${total}`, `${e.base} ＋（Scatter ${scatterCount} − 3）＝ ${total}`) : "—"],
          [T("Win multiplier", "贏分倍數"), e ? `×${e.mult}` : "—"],
        ]);
        if (e) html += progress(ctx.act.round + 1, total, T(`${e.name} spins`, `${e.name} 轉數`));
      }
      return html;
    },
  });
})();
