export default defineNuxtConfig({
    modules: ['@nuxtjs/tailwindcss'],
    runtimeConfig: {
        public: {
            apiBase: process.env.NUXT_PUBLIC_API_BASE
        }
    },
    tailwindcss: {
        configPath: '~/tailwind.config.js'
    },
    compatibilityDate: '2025-03-02',
})