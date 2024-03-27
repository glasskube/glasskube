import type { SidebarsConfig } from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  guides: [
    {
      type: 'category',
      label: 'Package Installation Guides',
      link: {
        type: 'generated-index',
        title: 'Package Installation Guides',
        description:
          '⚠️ Learn how to install your favorite Kubernetes add-ons using the Glasskube package manager ⚠️',
        slug: '/categories/',
      },
      items: ['cert-manager', 'ingress-nginx'],
    },
  ],
};

export default sidebars;
