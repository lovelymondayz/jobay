import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
  },
  build: {
    chunkSizeWarningLimit: 1500,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes("node_modules")) {
            if (
              id.includes("three") ||
              id.includes("@react-three") ||
              id.includes("quick_flipbook")
            ) {
              return "vendor_3d";
            }
            if (id.includes("react-router-dom")) {
              return "vendor_router";
            }
            if (id.includes("lucide-react")) {
              return "vendor_icons";
            }
            return "vendor";
          }
        },
      },
    },
  },
});
