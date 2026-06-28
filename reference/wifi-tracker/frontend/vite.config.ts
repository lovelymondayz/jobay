import { defineConfig, loadEnv } from "vite";
import viteReact from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { TanStackRouterVite } from "@tanstack/router-plugin/vite";
import { resolve } from "node:path";

export default defineConfig(({ mode }) => {
  // Load env vars based on the current mode (development, production, etc.)
  const env = loadEnv(mode, process.cwd(), "");

  return {
    plugins: [
      TanStackRouterVite({ autoCodeSplitting: true }),
      viteReact(),
      tailwindcss(),
    ],
    resolve: {
      alias: {
        "@": resolve(__dirname, "./src"),
      },
    },

    // Optional: make env accessible in your config
    define: {
      __VITE_API_URL__: JSON.stringify(env.VITE_API_URL),
    },
    server: {
      host: true,
    },
  };
});
