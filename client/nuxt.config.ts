export default defineNuxtConfig({
  compatibilityDate: "2026-05-09",
  devtools: { enabled: false },
  css: ["~/assets/css/main.css"],
  runtimeConfig: {
    backendBaseUrl: process.env.BACKEND_BASE_URL || "http://127.0.0.1:8080",
    public: {
      apiProxyBase: "/api/backend"
    }
  }
})
