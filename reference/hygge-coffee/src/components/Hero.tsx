export default function Hero() {
  return (
    <section className="pt-32 pb-20 px-4 sm:px-6 lg:px-8 bg-gradient-to-br from-background via-background to-accent/10">
      <div className="max-w-6xl mx-auto">
        <div className="grid md:grid-cols-2 gap-12 items-center min-h-[500px]">
          <div>
            <h1 className="text-5xl sm:text-6xl font-serif font-bold mb-6 text-foreground leading-tight">
              Welcome to <span className="gradient-text">hygge-coffee</span>
            </h1>
            <p className="text-lg text-secondary mb-8 leading-relaxed">
              Discover the perfect blend of cozy warmth and exceptional coffee.
              We create moments of connection through carefully crafted
              beverages and a welcoming atmosphere.
            </p>
            <div className="flex flex-col sm:flex-row gap-4">
              <button className="bg-primary text-white px-8 py-3 rounded-lg hover:bg-primary/90 transition-all transform hover:scale-105 font-semibold">
                Explore Menu
              </button>
              <button className="border-2 border-primary text-primary px-8 py-3 rounded-lg hover:bg-primary/10 transition-colors font-semibold">
                View Hours
              </button>
            </div>
          </div>
          <div className="relative h-96 sm:h-[500px] rounded-2xl overflow-hidden shadow-2xl">
            <div className="absolute inset-0 bg-gradient-to-br from-primary/20 via-accent/20 to-primary/30 flex items-center justify-center">
              <div className="text-6xl">☕</div>
            </div>
            <div className="absolute inset-0 opacity-10">
              <svg viewBox="0 0 100 100" className="w-full h-full">
                <circle
                  cx="50"
                  cy="50"
                  r="40"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="0.5"
                  opacity="0.3"
                />
                <circle
                  cx="50"
                  cy="50"
                  r="30"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="0.5"
                  opacity="0.3"
                />
                <circle
                  cx="50"
                  cy="50"
                  r="20"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="0.5"
                  opacity="0.3"
                />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
