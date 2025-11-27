/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                primary: {
                    DEFAULT: '#1677ff',
                    hover: '#4096ff',
                    active: '#0958d9',
                },
                success: '#52c41a',
                warning: '#faad14',
                error: '#ff4d4f',
                info: '#1677ff',
            },
        },
    },
    plugins: [],
    corePlugins: {
        preflight: true,
    },
}
