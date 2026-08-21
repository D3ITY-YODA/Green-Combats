import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./lib/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        forest: { DEFAULT: "rgb(var(--color-forest) / <alpha-value>)", deep: "rgb(var(--color-forest-deep) / <alpha-value>)" },
        sage: { DEFAULT: "rgb(var(--color-sage) / <alpha-value>)", soft: "rgb(var(--color-sage-soft) / <alpha-value>)" },
        sky: { DEFAULT: "rgb(var(--color-sky) / <alpha-value>)", soft: "rgb(var(--color-sky-soft) / <alpha-value>)" },
        background: { DEFAULT: "rgb(var(--color-bg) / <alpha-value>)", mist: "rgb(var(--color-bg-mist) / <alpha-value>)", stone: "rgb(var(--color-bg-stone) / <alpha-value>)" },
        text: { charcoal: "rgb(var(--color-text-charcoal) / <alpha-value>)", muted: "rgb(var(--color-text-muted) / <alpha-value>)" },
        status: {
          normal: "rgb(var(--color-status-normal) / <alpha-value>)",
          watch: "rgb(var(--color-status-watch) / <alpha-value>)",
          important: "rgb(var(--color-status-important) / <alpha-value>)",
          emergency: "rgb(var(--color-status-emergency) / <alpha-value>)",
          delayed: "rgb(var(--color-status-delayed) / <alpha-value>)",
          unavailable: "rgb(var(--color-status-unavailable) / <alpha-value>)",
        },
      },
      fontSize: {
        display: ["48px", { lineHeight: "1.1", letterSpacing: "-0.02em" }],
        page: ["32px", { lineHeight: "1.2", letterSpacing: "-0.01em" }],
        section: ["22px", { lineHeight: "1.3" }],
        card: ["17px", { lineHeight: "1.4", fontWeight: "600" }],
        body: ["16px", { lineHeight: "1.5" }],
        supporting: ["13px", { lineHeight: "1.5" }],
        metadata: ["12px", { lineHeight: "1.4", letterSpacing: "0.02em" }],
      },
    },
  },
  plugins: [],
};
export default config;
