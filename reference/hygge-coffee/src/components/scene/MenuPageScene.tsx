import { lazy, Suspense } from "react";

const FlipBookScene = lazy(() => import("./FlipBookScene"));

export default function MenuPageScene() {
  return (
    <div className="relative min-h-[70vh] overflow-hidden bg-gradient-to-b from-slate-900 via-slate-950 to-slate-980 text-white">
      {/* <div className="absolute inset-x-0 top-0 z-10 px-6 py-8 text-center pointer-events-none sm:px-12">
        <h1 className="font-serif text-4xl font-extrabold tracking-tight sm:text-5xl">
          Menu Flipbook
        </h1>
        <p className="mt-3 text-sm text-slate-300 sm:text-base">
          Swipe through the menu with a smooth, immersive flipbook layout that
          matches the landing page style.
        </p>
      </div> */}

      <button className="absolute z-20 px-4 py-3 text-sm text-white transition rounded-full prev bottom-8 left-6 bg-white/10 backdrop-blur-xl hover:bg-white/20">
        Previous
      </button>

      <button className="absolute z-20 px-4 py-3 text-sm text-white transition rounded-full next bottom-8 right-6 bg-white/10 backdrop-blur-xl hover:bg-white/20">
        Next
      </button>

      <div className="absolute inset-0">
        <Suspense
          fallback={
            <div className="flex items-center justify-center h-full text-sm bg-slate-950 text-slate-200">
              Loading menu preview...
            </div>
          }
        >
          <FlipBookScene />
        </Suspense>
      </div>
    </div>
  );
}
