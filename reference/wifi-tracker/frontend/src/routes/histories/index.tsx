import { createFileRoute } from "@tanstack/react-router";
import HistoriesPage from "@/pages/HistoriesPage";

export const Route = createFileRoute("/histories/")({
  component: HistoriesPage,
});
