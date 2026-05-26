import { defineConfig, presetUno, presetAttributify } from 'unocss'

export default defineConfig({
  presets: [presetUno(), presetAttributify()],
  shortcuts: {
    'page-card': 'bg-white rounded-lg shadow-sm p-4',
    'flex-center': 'flex items-center justify-center',
    'text-muted': 'text-gray-500 text-sm',
  },
})
