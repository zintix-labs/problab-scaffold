// demo_normal: marks the scatters on the trigger step; the side panel shows the
// trigger state and free-game progress.
// Panel text follows the UI language: { en, zh } objects or DEV.t({ en, zh }).
(() => {
  const { kv, progress } = DEV.util;
  const T = (en, zh) => DEV.t({ en, zh });
  // ext appears only on the trigger step: scatter_hits are the scatter cells
  const scatters = (ctx) => (ctx.act.acttype === "trigger" && ctx.act.ext ? ctx.act.ext.scatter_hits || [] : []);

  DEV.register(0, {
    title: { en: "demo_normal info", zh: "demo_normal 資訊" },
    modes: { 0: "BaseGame", 1: "FreeGame" },
    cellClass: (idx, ctx) => (scatters(ctx).includes(idx) ? ["mark"] : []),
    badge: (idx, ctx) => (scatters(ctx).includes(idx) ? "S" : ""),
    panel(ctx) {
      const ext = ctx.ext;
      let html = kv([
        [T("Triggered", "觸發"), ext ? T("yes", "是") : T("not yet", "尚未")],
        ["Scatter", ext ? ext.scatters : "—"],
      ]);
      if (ctx.mode.modeid === 1) html += progress(ctx.act.round + 1, ctx.game.fixed.free_round, T("FG spins", "FG 轉數"));
      return html;
    },
  });
})();
