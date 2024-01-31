import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type {Options as IdealImageOptions} from '@docusaurus/plugin-ideal-image';
import type * as Preset from '@docusaurus/preset-classic';
import {EnumChangefreq} from 'sitemap';

const config: Config = {
  title: 'Glasskube.dev',
  tagline: '🧊 The missing Package Manager for Kubernetes 📦',
  favicon: 'img/favicon.png',
  trailingSlash: true,

  // Set the production url of your site here
  url: 'https://glasskube.dev',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'glasskube', // Usually your GitHub org/user name.
  projectName: 'glasskube', // Usually your repo name.

  onBrokenLinks: 'warn',
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
    '@docusaurus/theme-mermaid',
    [
      '@docusaurus/plugin-ideal-image',
      /** @type {import("@docusaurus/plugin-ideal-image").PluginOptions} */
      {
        quality: 70,
        max: 1030, // max resized image's size.
        min: 640, // min resized image's size. if original is lower, use that size.
        steps: 2, // the max number of images generated between min and max (inclusive)
        disableInDev: false,
      } satisfies IdealImageOptions,
    ],
    [
      'content-docs',
      {
        id: 'guides',
        path: 'guides',
        routeBasePath: 'guides',
        editCurrentVersion: true,
        sidebarPath: './sidebar-guides.ts',
        showLastUpdateAuthor: true,
        showLastUpdateTime: true,
      } satisfies DocsOptions,
    ],
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
          customCss: ['./src/css/custom.css']
        },
        sitemap: {
          changefreq: EnumChangefreq.DAILY,
          priority: 1,
          ignorePatterns: ['/blog/archive','/blog/tags', '/blog/tags/**'],
          filename: 'sitemap.xml',
        }
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
    colorMode: {
      respectPrefersColorScheme: true,
    },
    docs: {
      sidebar: {
        hideable: true,
        autoCollapseCategories: true,
      },
    },
    announcementBar: {
      id: 'announcementBar-0', // Increment on change
      // content: '⭐️ If you like <code>glasskube</code>, give it a star on <a target="_blank" rel="noopener noreferrer" href="https://github.com/glasskube/glasskube">GitHub</a> and follow us on <a target="_blank" rel="noopener noreferrer" href="https://x.com/glasskube">X</a> ⭐️',
      content: `🎉️ <a target="_blank" href="https://github.com/glasskube/glasskube"><code>glasskube/glasskube</code></a> is launching its technical concept on GitHub 🥳️ <a target="_blank" rel="noopener noreferrer" href="https://github.com/glasskube/glasskube">Leave a star to support us</a> ⭐️`,
      isCloseable: false
    },
    image: 'https://opengraph.githubassets.com/3fbd03d4d860275ee154ca566f24ecce9243e229fe367523fbcab52e8b43db3f/glasskube/glasskube',
    navbar: {
      title: 'Glasskube',
      logo: {
        alt: 'Glasskube Logo',
        src: 'img/glasskube-logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docs',
          position: 'left',
          label: 'Docs',
        },
        {
          to: '/guides/cert-manager',
          position: 'left',
          label: 'Guides',
          activeBaseRegex: `/guides/`,
        },
        { to: '/blog', label: 'Blog', position: 'left' },
        { to: '/roadmap', label: 'Roadmap', position: 'left' },
        { to: '/packages', label: 'Packages', position: 'left' },

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
              label: 'Contact / Book appoitment',
              href: 'https://cal.glasskube.eu/team/founder/30min'
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
      copyright: `<img src="/img/glasskube-logo-white.png" class="footer-logo"/><br>Copyright © ${new Date().getFullYear()} Glasskube<br>Built with Docusaurus.`,
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
