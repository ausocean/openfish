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
      { text: 'API Documentation', link: '/api/' },
      { text: 'AusOcean', link: 'https://www.ausocean.org/'}
    ],

    sidebar: [  
      { text: 'Project Overview', link: '/project-overview' },
      { text: 'Getting Started', link: '/getting-started' },
      {
        text: 'API Documentation',
        items: [
          { text: 'Introduction', link: '/api/'},
          { text: 'API Usage', link: '/api/api-usage', items: [
            { text: 'Capture Sources', link: '/api/capture-sources'},
            { text: 'Video Streams', link: '/api/video-streams'},
            { text: 'Annotations', link: '/api/annotations'},
          ]},

        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/ausocean/openfish' }
    ],

    search: {
      provider: 'local'
    }
  }
})




