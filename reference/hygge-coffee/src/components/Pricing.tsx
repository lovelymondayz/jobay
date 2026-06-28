import { Check } from 'lucide-react'

const pricingPlans = [
  {
    name: 'Small',
    size: '8 oz',
    prices: {
      espresso: 2.50,
      cappuccino: 4.00,
      latte: 4.50,
      specialty: 5.00,
    }
  },
  {
    name: 'Medium',
    size: '12 oz',
    prices: {
      espresso: 3.00,
      cappuccino: 4.50,
      latte: 5.00,
      specialty: 5.75,
    },
    featured: true
  },
  {
    name: 'Large',
    size: '16 oz',
    prices: {
      espresso: 3.50,
      cappuccino: 5.00,
      latte: 5.50,
      specialty: 6.50,
    }
  }
]

const addOns = [
  { name: 'Extra Shot', price: 0.75 },
  { name: 'Flavored Syrup', price: 0.50 },
  { name: 'Extra Cream', price: 0.50 },
  { name: 'Whipped Cream', price: 0.75 },
  { name: 'Almond Milk', price: 0.75 },
  { name: 'Oat Milk', price: 0.75 },
]

export default function Pricing() {
  return (
    <section id="pricing" className="py-20 px-4 sm:px-6 lg:px-8 bg-background">
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl sm:text-5xl font-serif font-bold mb-4 text-foreground">
            Pricing
          </h2>
          <p className="text-lg text-secondary">Simple, transparent pricing for all our offerings</p>
        </div>

        {/* Size Pricing */}
        <div className="mb-16">
          <h3 className="text-2xl font-serif font-bold mb-8 text-center text-foreground">Coffee Sizes</h3>
          <div className="grid md:grid-cols-3 gap-8 mb-12">
            {pricingPlans.map((plan) => (
              <div
                key={plan.name}
                className={`rounded-2xl p-8 transition-all transform hover:scale-105 ${
                  plan.featured
                    ? 'bg-primary text-white shadow-2xl -translate-y-4'
                    : 'bg-accent/10 text-foreground border-2 border-accent/20'
                }`}
              >
                <h4 className={`text-2xl font-serif font-bold mb-2 ${plan.featured ? 'text-white' : 'text-primary'}`}>
                  {plan.name}
                </h4>
                <p className={`text-sm mb-6 ${plan.featured ? 'text-white/80' : 'text-secondary'}`}>
                  {plan.size}
                </p>
                
                <div className="space-y-3 mb-8">
                  <div className={`flex justify-between py-2 border-b ${plan.featured ? 'border-white/20' : 'border-accent/20'}`}>
                    <span>Espresso</span>
                    <span className="font-bold">${plan.prices.espresso.toFixed(2)}</span>
                  </div>
                  <div className={`flex justify-between py-2 border-b ${plan.featured ? 'border-white/20' : 'border-accent/20'}`}>
                    <span>Cappuccino</span>
                    <span className="font-bold">${plan.prices.cappuccino.toFixed(2)}</span>
                  </div>
                  <div className={`flex justify-between py-2 border-b ${plan.featured ? 'border-white/20' : 'border-accent/20'}`}>
                    <span>Latte</span>
                    <span className="font-bold">${plan.prices.latte.toFixed(2)}</span>
                  </div>
                  <div className={`flex justify-between py-2 ${plan.featured ? 'border-white/20' : 'border-accent/20'}`}>
                    <span>Specialty</span>
                    <span className="font-bold">${plan.prices.specialty.toFixed(2)}</span>
                  </div>
                </div>

                {plan.featured && (
                  <div className="bg-white/20 text-white text-sm py-2 px-4 rounded-lg text-center font-semibold">
                    Most Popular
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* Add-ons */}
        <div className="bg-accent/10 rounded-2xl p-12">
          <h3 className="text-2xl font-serif font-bold mb-8 text-center text-foreground">Customizations & Add-ons</h3>
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {addOns.map((addOn) => (
              <div key={addOn.name} className="flex items-center justify-between bg-background rounded-xl p-4">
                <div className="flex items-center gap-3">
                  <Check className="text-primary" size={20} />
                  <span className="font-medium text-foreground">{addOn.name}</span>
                </div>
                <span className="font-bold text-primary">${addOn.price.toFixed(2)}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Special Offers */}
        <div className="mt-16 text-center">
          <h3 className="text-2xl font-serif font-bold mb-8 text-foreground">Special Offers</h3>
          <div className="grid md:grid-cols-2 gap-8">
            <div className="bg-primary text-white rounded-2xl p-8">
              <h4 className="text-xl font-bold mb-2">Loyalty Program</h4>
              <p className="mb-4">Get a free coffee every 10th purchase!</p>
              <button className="bg-white text-primary px-6 py-2 rounded-lg hover:bg-white/90 transition-colors font-semibold">
                Join Today
              </button>
            </div>
            <div className="bg-accent text-white rounded-2xl p-8">
              <h4 className="text-xl font-bold mb-2">Daily Specials</h4>
              <p className="mb-4">20% off all drinks before 8 AM on weekdays</p>
              <button className="bg-white text-accent px-6 py-2 rounded-lg hover:bg-white/90 transition-colors font-semibold">
                Learn More
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
