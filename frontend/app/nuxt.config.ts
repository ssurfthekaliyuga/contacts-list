export default defineNuxtConfig({
    modules: ['@nuxtjs/tailwindcss'],
    runtimeConfig: {
        public: {
            apiBase: process.env.NUXT_PUBLIC_BACKEND_URL
        }
    },
    tailwindcss: {
        configPath: '~/tailwind.config.js'
    },
    compatibilityDate: '2025-03-02'
})