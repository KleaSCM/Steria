/**
 * Vite Configuration.
 * 
 * Configures the build process for Svelte and Electron integration.
 * 
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */

import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()], // Use the imported plugin function
  base: './', // Important for Electron to find assets
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  }
})
