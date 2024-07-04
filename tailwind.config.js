/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["internal/resources/views/*.html"],
  safelist: [
    'bg-lime-600',
    'bg-red-500',
    'bg-orange-400',
    'text-white',
    'h-full',
  ],
  theme: {
    extend: {
      colors: {
        'avocado': {
          '50': '#f3f5f0',
          '100': '#e5e9de',
          '200': '#ced5c1',
          '300': '#b0ba9c',
          '400': '#93a07b',
          '500': '#849469',
          '600': '#5c6848',
          '700': '#48513a',
          '800': '#3b4232',
          '900': '#343a2d',
          '950': '#1a1e15',
        }
      }
    },
  },
  plugins: [],
  darkMode: ['variant', [
    '@media (prefers-color-scheme: dark) { &:not(.light *) }',
    '&:is(.dark *)',
  ]],
}
