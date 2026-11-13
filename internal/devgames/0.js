// demo_normal: marks the scatters on the trigger step; the side panel shows the
// trigger state and free-game progress.
(() => {
  const { kv, progress } = DEV.util;
  // ext appears only on the trigger step: scatter_hits are the scatter cells
  const scatters = (ctx) => (ctx.act.acttype === "trigger" && ctx.act.ext ? ctx.act.ext.scatter_hits || [] : []);

  DEV.register(0, {
    title: "demo_normal 資訊",
    cellClass: (idx, ctx) => (scatters(ctx).includes(idx) ? ["mark"] : []),
    badge: (idx, ctx) => (scatters(ctx).includes(idx) ? "S" : ""),
    panel(ctx) {
      const ext = ctx.ext;
      let html = kv([
        ["觸發", ext ? "是" : "尚未"],
        ["Scatter", ext ? ext.scatters : "—"],
      ]);
      if (ctx.mode.modeid === 1) html += progress(ctx.act.round + 1, ctx.game.fixed.free_round, "FG 轉數");
      return html;
    },
  });
})();
