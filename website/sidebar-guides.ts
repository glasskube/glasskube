import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {

  guides: [
    {
      type: 'category',
      label: 'Categories',
      link: {
        type: 'generated-index',
        title: 'Categories',
        description: "⚠️ Learn how to install your favorite Kubernetes add-on's using the Glasskube package manager ⚠️",
        slug: '/categories/',
      },
      items: [
        'cert-manager',
        'ingress-nginx',
      ],
    },
  ]
};

export default sidebars;