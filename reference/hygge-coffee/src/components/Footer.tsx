import { MapPin, Clock, Phone, Mail, Globe, Share2, Heart } from "lucide-react";

export default function Footer() {
  return (
    <footer className="px-4 py-16 bg-foreground text-background sm:px-6 lg:px-8">
      <div className="max-w-6xl mx-auto">
        <div className="grid gap-12 mb-12 md:grid-cols-4">
          {/* Brand */}
          <div>
            <h3 className="mb-4 font-serif text-2xl font-bold">
              ☕ hygge-coffee
            </h3>
            <p className="leading-relaxed text-background/80">
              Creating moments of warmth and connection through exceptional
              coffee and a welcoming atmosphere.
            </p>
          </div>

          {/* Hours */}
          <div>
            <h4 className="flex items-center gap-2 mb-4 text-lg font-bold">
              <Clock size={20} /> Hours
            </h4>
            <ul className="space-y-2 text-background/80">
              <li>Mon - Fri: 6:00 AM - 8:00 PM</li>
              <li>Sat: 7:00 AM - 9:00 PM</li>
              <li>Sun: 8:00 AM - 7:00 PM</li>
              <li className="pt-2 font-semibold text-background">
                Always open for you! ☕
              </li>
            </ul>
          </div>

          {/* Location */}
          <div>
            <h4 className="flex items-center gap-2 mb-4 text-lg font-bold">
              <MapPin size={20} /> Location
            </h4>
            <address className="not-italic leading-relaxed text-background/80">
              <p className="mb-1 font-semibold">Main Store</p>
              <p>Jl. BSD Grand Boulevard</p>
              <p> Tangerang, Banten 15345</p>
            </address>
          </div>

          {/* Contact */}
          <div>
            <h4 className="mb-4 text-lg font-bold">Contact</h4>
            <div className="space-y-3 text-background/80">
              <a
                href="tel:+18005551234"
                className="flex items-center gap-2 transition-colors hover:text-background"
              >
                <Phone size={18} /> (+62) 818 710 117
              </a>
              <a
                href="mailto:hello@hygge.coffee"
                className="flex items-center gap-2 transition-colors hover:text-background"
              >
                <Mail size={18} /> hello@hygge.coffee
              </a>

              {/* Social Media */}
              <div className="flex gap-4 pt-4">
                <a href="#" className="transition-colors hover:text-background">
                  <Globe size={20} />
                </a>
                <a href="#" className="transition-colors hover:text-background">
                  <Heart size={20} />
                </a>
                <a href="#" className="transition-colors hover:text-background">
                  <Share2 size={20} />
                </a>
              </div>
            </div>
          </div>
        </div>

        {/* Bottom */}
        <div className="flex flex-col items-center justify-between gap-4 pt-8 text-sm border-t border-background/20 sm:flex-row text-background/70">
          <p>&copy; 2024 hygge-coffee. All rights reserved.</p>
          <div className="flex gap-6">
            <a href="#" className="transition-colors hover:text-background">
              Privacy Policy
            </a>
            <a href="#" className="transition-colors hover:text-background">
              Terms of Service
            </a>
            <a href="#" className="transition-colors hover:text-background">
              Sustainability
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
