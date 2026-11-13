// demo_choice: after a BG trigger the player picks one of two entries. The side
// panel explains each entry, derives the free-game rounds from the checkpoint's
// scatter count and shows free-game progress; the trigger step marks the scatters.
(() => {
  const { kv, progress, esc } = DEV.util;
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
    title: "demo_choice 資訊",
    cellClass: (idx, ctx) => (scatters(ctx).includes(idx) ? ["mark"] : []),
    badge: (idx, ctx) => (scatters(ctx).includes(idx) ? "S" : ""),
    panel(ctx) {
      const table = entries(ctx.game.fixed);
      const seg = ctx.segment;
      let html = "";
      if (ctx.segIdx === 0) {
        const cp = seg.spin_state && seg.spin_state.cp;
        html += kv([["觸發", ctx.ext ? "是" : "尚未"], ["Scatter", ctx.ext ? ctx.ext.scatters : "—"]]);
        if (cp && seg.choose_next) {
          html += `<div class="muted">選項（Scatter ${cp.scatters} 個）</div>` + kv(seg.choose_next.map((entry, i) => {
            const e = table[entry];
            return [`#${i} 入口 ${entry}`, e ? `${e.name}：${rounds(e, cp.scatters)} 轉${e.mult > 1 ? `，贏分 ×${e.mult}` : ""}` : "未知入口"];
          }));
        }
      } else {
        const prev = ctx.bet.segments[ctx.segIdx - 1];
        const entry = prev.choose_next[seg.choice];
        const e = table[entry];
        const scatterCount = prev.spin_state.cp ? prev.spin_state.cp.scatters : 3;
        const total = e ? rounds(e, scatterCount) : 0;
        html += kv([
          ["選擇", `#${seg.choice} 入口 ${esc(entry)}（${e ? e.name : "未知"}）`],
          ["轉數推算", e ? `${e.base} ＋（Scatter ${scatterCount} − 3）＝ ${total}` : "—"],
          ["贏分倍數", e ? `×${e.mult}` : "—"],
        ]);
        if (e) html += progress(ctx.act.round + 1, total, `${e.name} 轉數`);
      }
      return html;
    },
  });
})();
