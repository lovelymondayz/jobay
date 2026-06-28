/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: '#faf8f5',
        foreground: '#2d2d2d',
        primary: '#8b6f47',
        accent: '#d4a574',
        secondary: '#4a4a4a',
      },
      fontFamily: {
        sans: ['Geist', 'sans-serif'],
        serif: ['Lora', 'serif'],
      },
    },
  },
  plugins: [],
}
