import { defineConfig } from "astro/config";

export default defineConfig({
  devToolbar: { enabled: false },
  // bind mount上の変更検知(Windows/WSL対策)
  vite: { server: { watch: { usePolling: true } } },
});
