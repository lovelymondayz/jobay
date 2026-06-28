import { createFileRoute } from "@tanstack/react-router";

import RegistUserPage from "@/pages/RegistUserPage";

export const Route = createFileRoute("/regist-user/")({
  component: RegistUserPage,
});
