import { defineConfig } from "vitepress";

export default defineConfig({
  title: "Mahaam",
  description: "Mahaam app docs",
  cleanUrls: true,
  head: [
    ["link", { rel: "icon", href: "/logo.png" }],
    ["link", { rel: "preconnect", href: "https://fonts.googleapis.com" }],
    ["link", { rel: "preconnect", href: "https://fonts.gstatic.com", crossorigin: "" }],
    [
      "link",
      {
        href: "https://fonts.googleapis.com/css2?family=JetBrains+Mono:ital,wght@0,100..800;1,100..800&display=swap",
        rel: "stylesheet",
      },
    ],
  ],

  themeConfig: {
    // Logo
    logo: "/logo.png",

    // Navigation
    nav: [],

    // Sidebar
    sidebar: {
      "/": [
        {
          text: "",
          items: [{ text: "Home", link: "/" }],
        },
        {
          text: "Service Setup",
          items: [
            { text: "Create", link: "/setup/creation" },
            { text: "Dependencies", link: "/setup/dependencies" },
          ],
        },
        {
          text: "Service Design",
          items: [
            { text: "Intro", link: "/design/intro" },
            { text: "Functions", link: "/design/functions" },
            { text: "Model", link: "/design/model" },
            { text: "Design", link: "/design/design" },
            { text: "Structure", link: "/design/structure" },
            { text: "Maintainability", link: "/design/maintainability" },
          ],
        },
        {
          text: "Service Module",
          items: [
            { text: "Intro", link: "/module/intro" },
            { text: "Repo", link: "/module/repos" },
            { text: "Service", link: "/module/services" },
            { text: "Controller", link: "/module/controllers" },
            { text: "Model", link: "/module/models" },
            { text: "Transactions", link: "/module/transactions" },
          ],
        },
        {
          text: "Service Infra",
          items: [
            { text: "Intro", link: "/infra/intro" },
            { text: "Exceptions", link: "/infra/exceptions" },
            { text: "Req Context", link: "/infra/request-context" },
            { text: "Security", link: "/infra/security" },
            { text: "Logging", link: "/infra/logging" },
            { text: "Validation", link: "/infra/validations" },
            { text: "Swagger", link: "/infra/swagger" },
            { text: "Config", link: "/infra/configs" },
            { text: "Cache", link: "/infra/caching" },
            { text: "Monitor", link: "/infra/monitoring" },
          ],
        },
        {
          text: "Service Testing",
          items: [{ text: "Test", link: "/test/test" }],
        },
      ],
    },

    // Social links
    socialLinks: [
      {
        icon: {
          svg: '<svg height="32" aria-hidden="true" viewBox="0 0 24 24" version="1.1" width="32" data-view-component="true" style="color: var(--vp-c-text-1)"><path fill="currentColor" d="M12 1C5.923 1 1 5.923 1 12c0 4.867 3.149 8.979 7.521 10.436.55.096.756-.233.756-.522 0-.262-.013-1.128-.013-2.049-2.764.509-3.479-.674-3.699-1.292-.124-.317-.66-1.293-1.127-1.554-.385-.207-.936-.715-.014-.729.866-.014 1.485.797 1.691 1.128.99 1.663 2.571 1.196 3.204.907.096-.715.385-1.196.701-1.471-2.448-.275-5.005-1.224-5.005-5.432 0-1.196.426-2.186 1.128-2.956-.111-.275-.496-1.402.11-2.915 0 0 .921-.288 3.024 1.128a10.193 10.193 0 0 1 2.75-.371c.936 0 1.871.123 2.75.371 2.104-1.43 3.025-1.128 3.025-1.128.605 1.513.221 2.64.111 2.915.701.77 1.127 1.747 1.127 2.956 0 4.222-2.571 5.157-5.019 5.432.399.344.743 1.004.743 2.035 0 1.471-.014 2.654-.014 3.025 0 .289.206.632.756.522C19.851 20.979 23 16.854 23 12c0-6.077-4.922-11-11-11Z"></path></svg>',
        },
        link: "https://github.com/ayasrah/mahaam",
      },
      //   {
      //     icon: {
      //       svg: '<svg class="kOqhQd" aria-hidden="true" viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg"><path fill="none" d="M0,0h40v40H0V0z"></path><g><path d="M19.7,19.2L4.3,35.3c0,0,0,0,0,0c0.5,1.7,2.1,3,4,3c0.8,0,1.5-0.2,2.1-0.6l0,0l17.4-9.9L19.7,19.2z" fill="#EA4335"></path><path d="M35.3,16.4L35.3,16.4l-7.5-4.3l-8.4,7.4l8.5,8.3l7.5-4.2c1.3-0.7,2.2-2.1,2.2-3.6C37.5,18.5,36.6,17.1,35.3,16.4z" fill="#FBBC04"></path><path d="M4.3,4.7C4.2,5,4.2,5.4,4.2,5.8v28.5c0,0.4,0,0.7,0.1,1.1l16-15.7L4.3,4.7z" fill="#4285F4"></path><path d="M19.8,20l8-7.9L10.5,2.3C9.9,1.9,9.1,1.7,8.3,1.7c-1.9,0-3.6,1.3-4,3c0,0,0,0,0,0L19.8,20z" fill="#34A853"></path></g></svg>',
      //     },
      //     link: "https://play.google.com/store/apps/details?id=ayasrah.mahaam",
      //   },
      //   {
      //     icon: {
      //       svg: '<svg role="img" width="16" height="16" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" style="vertical-align: middle; margin-right: 8px;"><path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>',
      //     },
      //     link: "https://apps.apple.com/us/app/mahaam/id6502533759",
      //   },
      //   { icon: "x", link: "https://x.com/ayasrah0" },
      //   {
      //     icon: {
      //       svg: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24" focusable="false"><rect width="24" height="24" fill="#0A66C2" rx="4" /><path fill="#FFFFFF" d="M8 19H5v-9h3zM6.5 8.25A1.75 1.75 0 118.3 6.5a1.78 1.78 0 01-1.8 1.75zM19 19h-3v-4.74c0-1.42-.6-1.93-1.38-1.93A1.74 1.74 0 0013 14.19a.66.66 0 000 .14V19h-3v-9h2.9v1.3a3.11 3.11 0 012.7-1.4c1.55 0 3.36.86 3.36 3.66z"/></svg>',
      //     },
      //     link: "https://www.linkedin.com/in/ayasrah/",
      //   },
    ],

    // Footer
    footer: false,

    // Search
    search: {
      provider: "local",
    },

    // Edit link
    // editLink: {
    //   pattern: "https://github.com/ayasrah/mahaam/edit/main/mahaam-docs/docs/:path",
    //   text: "Edit",
    // },
  },

  // Markdown configuration
  markdown: {
    lineNumbers: true,
  },
});
