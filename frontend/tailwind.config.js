/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {},
  },
  plugins: [
    // Note: @tailwindcss/forms might need to be imported differently in ES modules
    // We'll add this after basic setup works
  ],
}
