import { createFileRoute } from "@tanstack/react-router";
import { MiningmapPage } from "@/pages/MiningmapPage";

export const Route = createFileRoute("/miningmap/")({
  component: MiningmapPage,
});
