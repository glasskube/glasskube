import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type {Options as IdealImageOptions} from '@docusaurus/plugin-ideal-image';
import type * as Preset from '@docusaurus/preset-classic';
import {EnumChangefreq} from 'sitemap';

const config: Config = {
  title: 'Glasskube',
  tagline: 'The easiest way to distribute enterprise software',
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
      './custom-blog-plugin',
      {id: 'blog', routeBasePath: 'blog', path: './blog'},
    ],
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
      },
    ],
    [
      '@docusaurus/plugin-content-docs',
      {
        id: 'glossary',
        path: 'glossary',
        routeBasePath: 'glossary',
        sidebarPath: './sidebar-glossary.ts',
        editUrl: 'https://github.com/glasskube/glasskube/tree/main/website/',
      },
    ],
    [
      'posthog-docusaurus',
      {
        apiKey: 'phc_EloQUW6cgfbTc0pI9c5CXElhQ4gVGRoBsrUAoakJVoQ',
        appUrl: 'https://p.glasskube.eu',
        ui_host: 'https://eu.posthog.com',
        enableInDevelopment: false,
      },
    ],
  ],
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/glasskube/glasskube/tree/main/website/',
        },
        blog: false,
        theme: {
          customCss: ['./src/css/custom.css'],
        },
        sitemap: {
          changefreq: EnumChangefreq.DAILY,
          priority: 1,
          ignorePatterns: ['/telemetry/', '/blog/authors/', '/blog/archive/', '/blog/tags/**'],
          filename: 'sitemap.xml',
        },
      } satisfies Preset.Options,
    ],
  ],
  markdown: {
    mermaid: true,
  },
  themeConfig: {
    colorMode: {
      defaultMode: 'dark',
      respectPrefersColorScheme: true,
    },
    docs: {
      sidebar: {
        hideable: true,
        autoCollapseCategories: true,
      },
    },
    announcementBar: {
      id: 'announcementBar-1', // Increment on change
      /* x-release-please-start-version */
      content: `🎉 We just released v0.26.0 of our Open Source Kubernetes Package Manager on <a href="https://github.com/glasskube/glasskube/" target="_blank">⭐ GitHub ⭐</a>.`,
      /* x-release-please-end */
      isCloseable: false,
    },
    image:
      'https://opengraph.githubassets.com/3fbd03d4d860275ee154ca566f24ecce9243e229fe367523fbcab52e8b43db3f/glasskube/glasskube',
    navbar: {
      title: 'Glasskube',
      logo: {
        alt: 'Glasskube Logo',
        src: 'img/glasskube-logo.svg',
      },
      items: [
        {label: 'Software Distribution', to: '/software-distribution/'},
        {label: 'Blog', to: '/blog/'},
        {label: 'Pricing', to: '/pricing'},
        {type: 'custom-wrapper', position: 'right'},
        // {label: 'Login', to: 'https://glasskube.cloud/', position: 'right'},
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Products',
          items: [
            {label: 'Glasskube Cloud', to: '/'},
            {label: 'Package Manager', to: '/products/package-manager/'},
            {label: 'Private Package Repository', to: '/products/hub/'},
          ],
        },
        {
          title: 'Resources',
          items: [
            {label: 'Blog', to: '/blog/'},
            {label: 'Glossary', to: '/glossary/'},
            {label: 'Package Manager Docs', to: '/docs/'},
            {label: 'Package Manager Guides', to: '/guides/'},
          ],
        },
        {
          title: 'Community',
          items: [
            {label: 'GitHub', href: 'https://github.com/glasskube/glasskube'},
            {label: 'Discord', href: 'https://discord.gg/SxH6KUCGH7'},
            {label: 'LinkedIn', href: 'https://www.linkedin.com/company/glasskube/'},
            {label: 'Twitter / X', href: 'https://x.com/glasskube'},
          ],
        },
        {
          title: 'More',
          items: [
            {label: 'Blog', to: '/blog/'},
            {label: 'Contact', to: '/contact/'},
            {label: 'Talk to founders', href: 'https://cal.glasskube.com/team/founder/30min'},
            {label: 'Signup for the wait list', href: 'https://glasskube.cloud/'},
          ],
        },
      ],
      copyright: `<img src="/img/glasskube-logo-white.png" class="footer-logo" alt="Glasskube Logo"/><br>Copyright © ${new Date().getFullYear()} Glasskube, Inc.<br>Built with Docusaurus.`,
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
