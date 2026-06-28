export default function About() {
  return (
    <section id="about" className="px-4 py-20 sm:px-6 lg:px-8 bg-background">
      <div className="max-w-6xl mx-auto">
        <div className="grid items-center gap-12 md:grid-cols-2">
          <div className="h-96 sm:h-[500px] rounded-2xl overflow-hidden shadow-xl">
            <div className="flex items-center justify-center w-full h-full bg-gradient-to-br from-primary/30 via-accent/30 to-primary/20">
              <div className="text-8xl">🏡</div>
            </div>
          </div>

          <div>
            <h2 className="mb-6 font-serif text-4xl font-bold sm:text-5xl text-foreground">
              About <span className="gradient-text">hygge-cafe</span>
            </h2>
            <p className="mb-6 text-lg leading-relaxed text-secondary">
              Established in 2020, hygge-cafe is dedicated to bringing the
              Scandinavian concept of hygge to every cup. We believe that coffee
              is more than just a beverage—it&apos;s a moment of warmth,
              connection, and comfort.
            </p>
            <p className="mb-8 text-lg leading-relaxed text-secondary">
              Our beans are carefully sourced from sustainable farms around the
              world, and every drink is crafted with passion and precision.
              Whether you&apos;re looking for your daily ritual or a cozy
              meeting spot, hygge-coffee is your sanctuary.
            </p>

            <div className="grid grid-cols-2 gap-6">
              <div className="p-6 bg-accent/10 rounded-xl">
                <div className="mb-2 text-3xl font-bold text-primary">100%</div>
                <p className="font-semibold text-secondary">Sustainable</p>
              </div>
              <div className="p-6 bg-accent/10 rounded-xl">
                <div className="mb-2 text-3xl font-bold text-primary">
                  1000+
                </div>
                <p className="font-semibold text-secondary">Daily Customers</p>
              </div>
              <div className="p-6 bg-accent/10 rounded-xl">
                <div className="mb-2 text-3xl font-bold text-primary">50+</div>
                <p className="font-semibold text-secondary">Menu Items</p>
              </div>
              <div className="p-6 bg-accent/10 rounded-xl">
                <div className="mb-2 text-3xl font-bold text-primary">4.9★</div>
                <p className="font-semibold text-secondary">Customer Rating</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
