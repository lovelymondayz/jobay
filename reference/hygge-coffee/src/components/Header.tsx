import { useState } from "react";
import { Menu, X } from "lucide-react";
import { Link } from "react-router-dom";

export default function Header() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <header className="fixed top-0 w-full bg-background/95 backdrop-blur-sm z-50 border-b border-accent/20">
      <nav className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Link to="/" className="text-2xl font-serif font-bold text-primary">
            ☕ hygge-coffee
          </Link>
        </div>

        {/* Desktop Menu */}
        <div className="hidden md:flex items-center gap-8">
          <a href="/#about" className="hover:text-primary transition-colors">
            About
          </a>
          <Link to="/menu" className="hover:text-primary transition-colors">
            Menu
          </Link>
          <a href="/#pricing" className="hover:text-primary transition-colors">
            Pricing
          </a>
          <a
            href="/#contact"
            className="bg-primary text-white px-6 py-2 rounded-lg hover:bg-primary/90 transition-colors"
          >
            Contact Us
          </a>
        </div>

        {/* Mobile Menu Button */}
        <button onClick={() => setIsOpen(!isOpen)} className="md:hidden p-2">
          {isOpen ? <X size={24} /> : <Menu size={24} />}
        </button>
      </nav>

      {/* Mobile Menu */}
      {isOpen && (
        <div className="md:hidden bg-background border-t border-accent/20 px-4 py-4 flex flex-col gap-4">
          <a href="/#about" className="hover:text-primary transition-colors">
            About
          </a>
          <Link to="/menu" className="hover:text-primary transition-colors">
            Menu
          </Link>
          <a href="/#pricing" className="hover:text-primary transition-colors">
            Pricing
          </a>
          <a
            href="/#contact"
            className="bg-primary text-white px-6 py-2 rounded-lg hover:bg-primary/90 transition-colors text-center"
          >
            Contact Us
          </a>
        </div>
      )}
    </header>
  );
}
