/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,jsx}",
  ],
  theme: {
    extend: {
      colors: {
        'sidebar-dark': '#1a1a1a',
        'content-light': '#ffffff',
        'primary': '#4F46E5',
        'primary-hover': '#4338CA',
      },
    },
  },
  plugins: [],
}
