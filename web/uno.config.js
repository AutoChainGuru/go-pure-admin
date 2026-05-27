import { defineConfig, presetUno, presetAttributify } from 'unocss'

export default defineConfig({
  darkMode: 'class',
  presets: [presetUno(), presetAttributify()],
  shortcuts: {
    'page-card': 'rounded-lg p-4',
    'flex-center': 'flex items-center justify-center',
    'text-muted': 'text-gray-500 text-sm',
  },
})
