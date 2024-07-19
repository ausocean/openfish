import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  base: '/openfish/',
  title: "OpenFish Documentation",
  description: "OpenFish is an open-source system written in Golang for classifying marine species. Tasks involve importing video or image data, classifying and annotating data (both manually and automatically), searching, and more. ",
  ignoreDeadLinks: [/^https?:\/\/localhost:5173/],
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Project Overview', link: '/project-overview' },
      { text: 'API Documentation', link: '/developer-docs/api' },
      { text: 'AusOcean', link: 'https://www.ausocean.org/'}
    ],

    sidebar: [  
      { text: 'Project Overview', link: '/project-overview' },
      { 
        text: 'Developer Documentation',
        items: [
          { text: 'Getting Started', link: '/developer-docs/getting-started' },
          { text: 'API Documentation', link: '/developer-docs/api' },
          { text: 'Contributing', link: '/developer-docs/contributing' },
        ]
      },
      
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/ausocean/openfish' }
    ],

    search: {
      provider: 'local'
    }
  }
})




