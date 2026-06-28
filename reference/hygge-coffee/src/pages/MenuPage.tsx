import { Link } from "react-router-dom";
import MenuPageScene from "../components/scene/MenuPageScene";
import Menu from "../components/Menu";

export default function MenuPage() {
  return (
    <main className="pb-20 bg-background">
      <div className="max-w-6xl px-4 pt-8 mx-auto sm:px-6 lg:px-8">
        <Link
          to="/"
          className="inline-flex items-center text-sm font-medium text-primary hover:text-primary/80"
        >
          ← Back to home
        </Link>

        <section className="mt-8 rounded-[32px] border border-white/10 bg-slate-950/95 shadow-2xl overflow-hidden">
          <div className="px-6 py-8 sm:px-12">
            <div className="max-w-3xl">
              <h1 className="font-serif text-4xl font-bold tracking-tight text-white sm:text-5xl">
                Menu Experience
              </h1>
              <p className="mt-4 text-base leading-7 text-slate-300">
                Explore the full hygge coffee menu in a rich flipbook experience
                designed to feel like part of the main landing page.
              </p>
            </div>
          </div>

          <MenuPageScene />
        </section>

        {/* <section className="mt-12 rounded-[32px] border border-accent/20 bg-background/80 p-8 shadow-lg">
          <h2 className="mb-6 font-serif text-3xl font-bold text-foreground">
            Full Menu
          </h2>
          <Menu />
        </section> */}
      </div>
    </main>
  );
}
