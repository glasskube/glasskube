import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Glasskube.dev',
  tagline: '🧊 Kubernetes Package Management the easy way 🔥',
  favicon: 'img/favicon.png',

  // Set the production url of your site here
  url: 'https://glasskube.dev',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'glasskube', // Usually your GitHub org/user name.
  projectName: 'glasskube', // Usually your repo name.

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },
  plugins: [
    'docusaurus-plugin-matomo',
    '@docusaurus/theme-mermaid'
  ],
  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl: 'https://github.com/glasskube/glasskube/tree/main/website/',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl: 'https://github.com/glasskube/glasskube/tree/main/website/',
        },
        theme: {
          customCss: [ './src/css/custom.css']
        },
      } satisfies Preset.Options,
    ],
  ],
  themes: [
    [
      require.resolve("@easyops-cn/docusaurus-search-local"),
      /** @type {import("@easyops-cn/docusaurus-search-local").PluginOptions} */
      ({
        hashed: true,
        indexBlog: false,
        docsRouteBasePath: '/'
      })
    ]
  ],
  markdown: {
    mermaid: true,
  },
  themeConfig: {
    colorMode : {
      respectPrefersColorScheme: true,
    },
    announcementBar: {
      id: 'announcementBar-0', // Increment on change
      // content: '⭐️ If you like <code>glasskube</code>, give it a star on <a target="_blank" rel="noopener noreferrer" href="https://github.com/glasskube/glasskube">GitHub</a> and follow us on <a target="_blank" rel="noopener noreferrer" href="https://x.com/glasskube">X</a> ⭐️',
      content: `🎉️ <a target="_blank" href="https://github.com/glasskube/glasskube"><code>glasskube/glasskube</code></a> is pre launching on GitHub 🥳️`,
      isCloseable: false
    },
    image: 'img/glasskube-social-card.jpg',
    navbar: {
      title: 'Glasskube',
      logo: {
        alt: 'Glasskube Logo',
        src: 'img/glasskube-logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'tutorialSidebar',
          position: 'left',
          label: 'Docs',
        },
        {to: '/roadmap', label: 'Roadmap', position: 'left'},
        {to: '/blog', label: 'Blog', position: 'left'},
        {
          href: 'https://github.com/glasskube/glasskube',
          label: 'GitHub',
          position: 'right',
        },
        {
          href: 'https://x.com/glasskube',
          label: 'Twitter / X',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Architecture',
              to: '/docs/',
            },
            {
              label: 'Getting started',
              to: '/docs/getting-started/install',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'Discord',
              href: 'https://discord.gg/SxH6KUCGH7',
            },
            {
              label: 'Twitter / X',
              href: 'https://x.com/glasskube',
            },
            {
              label: 'LinkedIn',
              href: 'https://www.linkedin.com/company/glasskube/',
            },
            {
              label: 'GitHub',
              href: 'https://github.com/glasskube/glasskube',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'Blog',
              to: '/blog',
            },
            {
              label: 'Imprint',
              to: '/imprint',
            },
            {
              label: 'Data privacy policy',
              to: '/data-privacy',
            },
          ],
        },
      ],
      copyright: `<img src="/img/glasskube-logo-white.png" width="25%" style="margin: 2rem"/><br>Copyright © ${new Date().getFullYear()} Glasskube<br>Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
    matomo: {
      matomoUrl: 'https://a.glasskube.eu/',
      siteId: '5',
      phpLoader: 'matomo.php',
      jsLoader: 'matomo.js',
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
