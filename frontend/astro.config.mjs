import { defineConfig } from "astro/config";

export default defineConfig({
  // bind mount上の変更検知(Windows/WSL対策)
  vite: { server: { watch: { usePolling: true } } },
});
