import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        'deep-forest': '#0D3B2E',
        'forest': '#176B4D',
        'sage': '#8FAF9D',
        'soft-sage': '#E6EFE9',
        'sky-blue': '#5D93B8',
        'soft-blue': '#E8F1F7',
        'warm-white': '#FCFDFC',
        'mist': '#F3F6F4',
        'stone': '#E7ECE8',
        'charcoal': '#17211D',
        'muted': '#64716B',
        'important': '#C76B35',
        'emergency': '#A93226',
        'delayed': '#6C7880',
      },
    },
  },
  plugins: [],
};
export default config;
