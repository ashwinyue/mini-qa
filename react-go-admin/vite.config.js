import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src'),
        },
    },
    server: {
        port: 5173,
        proxy: {
            // Python 后端 API 代理
            '/api': {
                target: 'http://localhost:8000',
                changeOrigin: true,
            },
            '/chat': {
                target: 'http://localhost:8000',
                changeOrigin: true,
            },
            '/greet': {
                target: 'http://localhost:8000',
                changeOrigin: true,
            },
            '/suggest': {
                target: 'http://localhost:8000',
                changeOrigin: true,
            },
            '/health': {
                target: 'http://localhost:8000',
                changeOrigin: true,
            },
            '/models': {
                target: 'http://localhost:8000',
                changeOrigin: true,
            },
        },
    },
})
