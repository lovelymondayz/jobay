"use client";

import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import React from "react";
import { Outlet, createRootRoute } from "@tanstack/react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import Sidebar from "@/components/SideBar";
import NotFoundPage from "@/pages/NotFoundPage";
import { Toaster } from "sonner";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 5 * 60 * 1000, // 5 minutes
    },
  },
});

export const Route = createRootRoute({
  component: () => {
    const [isSidebarOpen, setIsSidebarOpen] = React.useState(true);

    return (
      <QueryClientProvider client={queryClient}>
        <div className="flex">
          <Sidebar isOpen={isSidebarOpen} setIsOpen={setIsSidebarOpen} />
          <main
            className={`flex-1 transition-all duration-300 ease-in-out ${
              isSidebarOpen ? "ml-60" : "ml-16"
            }`}
          >
            <Outlet />
          </main>
        </div>
        <Toaster richColors position="top-right" />
        <TanStackRouterDevtools />
        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    );
  },
  notFoundComponent: NotFoundPage,
});
