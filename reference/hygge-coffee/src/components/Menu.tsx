import { Coffee } from 'lucide-react'

const menuCategories = [
  {
    name: 'Espresso Based',
    items: [
      { name: 'Espresso', description: 'Rich and bold shot of pure coffee', icon: '☕' },
      { name: 'Americano', description: 'Espresso diluted with hot water', icon: '💧' },
      { name: 'Cappuccino', description: 'Espresso with steamed milk and foam', icon: '🥛' },
      { name: 'Latte', description: 'Smooth espresso with velvety milk', icon: '🍶' },
      { name: 'Macchiato', description: 'Espresso marked with a touch of foam', icon: '☕' },
      { name: 'Flat White', description: 'Espresso with microfoam milk', icon: '🥛' },
    ]
  },
  {
    name: 'Cold Coffee',
    items: [
      { name: 'Iced Americano', description: 'Crisp and refreshing cold coffee', icon: '🧊' },
      { name: 'Iced Latte', description: 'Cold espresso with iced milk', icon: '🧊' },
      { name: 'Cold Brew', description: '24-hour slow-steeped perfection', icon: '⏱️' },
      { name: 'Iced Mocha', description: 'Coffee with chocolate and cold milk', icon: '🍫' },
      { name: 'Affogato', description: 'Espresso poured over vanilla ice cream', icon: '🍦' },
      { name: 'Nitro Cold Brew', description: 'Silky smooth nitrogen-infused brew', icon: '🌬️' },
    ]
  },
  {
    name: 'Specialty Drinks',
    items: [
      { name: 'Caramel Macchiato', description: 'Sweet caramel meets espresso elegance', icon: '✨' },
      { name: 'Vanilla Latte', description: 'Creamy vanilla-infused coffee', icon: '🌸' },
      { name: 'Hazelnut Cappuccino', description: 'Nutty and aromatic blend', icon: '🌰' },
      { name: 'Mocha', description: 'Perfect marriage of coffee and chocolate', icon: '🍫' },
      { name: 'Seasonal Spice', description: 'Limited edition seasonal flavors', icon: '✨' },
      { name: 'Cortado', description: 'Perfect balance of espresso and milk', icon: '⚖️' },
    ]
  }
]

export default function Menu() {
  return (
    <section id="menu" className="py-20 px-4 sm:px-6 lg:px-8 bg-accent/5">
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl sm:text-5xl font-serif font-bold mb-4 text-foreground">
            Our Menu
          </h2>
          <p className="text-lg text-secondary">Discover our carefully curated selection of premium coffee drinks</p>
        </div>

        <div className="grid md:grid-cols-3 gap-8">
          {menuCategories.map((category) => (
            <div key={category.name} className="bg-background rounded-2xl p-8 shadow-lg hover:shadow-xl transition-shadow">
              <h3 className="text-2xl font-serif font-bold mb-6 flex items-center gap-2 text-primary">
                <Coffee size={24} />
                {category.name}
              </h3>
              <div className="space-y-4">
                {category.items.map((item) => (
                  <div key={item.name} className="pb-4 border-b border-accent/20 last:border-b-0">
                    <div className="flex items-start gap-3">
                      <span className="text-2xl">{item.icon}</span>
                      <div>
                        <h4 className="font-semibold text-foreground">{item.name}</h4>
                        <p className="text-sm text-secondary">{item.description}</p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
