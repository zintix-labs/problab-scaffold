// demo_cascade: marks the scatters on the trigger step; the side panel shows the
// trigger state, the scatter award, the cascade count and free-game progress.
(() => {
  const { kv, progress } = DEV.util;
  const T = (en, zh) => DEV.t({ en, zh });
  const scatters = (ctx) => (ctx.act.acttype === "trigger" && ctx.act.ext ? ctx.act.ext.scatter_hits || [] : []);

  // Cascades in this round up to the current step (each clear act is one cascade)
  function clears(ctx) {
    const acts = ctx.mode.acts;
    const at = acts.indexOf(ctx.act);
    let n = 0;
    for (let i = 0; i <= at; i++) if (acts[i].round === ctx.act.round && acts[i].acttype === "clear") n++;
    return n;
  }

  DEV.register(1, {
    title: { en: "demo_cascade info", zh: "demo_cascade 資訊" },
    modes: { 0: "BaseGame", 1: "FreeGame" },
    cellClass: (idx, ctx) => (scatters(ctx).includes(idx) ? ["mark"] : []),
    badge: (idx, ctx) => (scatters(ctx).includes(idx) ? "S" : ""),
    panel(ctx) {
      const fixed = ctx.game.fixed;
      const ext = ctx.ext;
      const rows = [
        [T("Triggered", "觸發"), ext ? T("yes", "是") : T("not yet", "尚未")],
        ["Scatter", ext ? ext.scatters : "—"],
      ];
      if (ext) rows.push([T("Scatter pay", "Scatter 加給"), `${fixed.scatter_pay} × ${T("multiplier", "押注倍數")} ${ctx.segment.betmult} = ${fixed.scatter_pay * ctx.segment.betmult}`]);
      rows.push([T("Clears this round", "本 round 消除"), T(`${clears(ctx)} (max ${fixed.max_step})`, `${clears(ctx)} 次（上限 ${fixed.max_step}）`)]);
      let html = kv(rows);
      if (ctx.mode.modeid === 1) html += progress(ctx.act.round + 1, fixed.free_rounds, T("FG spins", "FG 轉數"));
      return html;
    },
  });
})();
