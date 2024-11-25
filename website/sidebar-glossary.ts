import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  glossary: [
    {
      type: 'category',
      label: 'Glossary',
      link: {
        type: 'generated-index',
        title: 'Kubernetes Glossary',
        description:
          'Learn about common Kubernetes and cloud-native terminology',
        slug: '/categories/glossary',
      },
      items: [
        'devops',
        'helm',
        'helm-chart',
        'kubernetes-configuration-management',
        'kubernetes-package-management',
        'kubernetes',
        'kubernetes-operator',
        'kustomize',
        'BYOC-bring-your-own-cloud',
        'SaaS-software-as-a-service',
        'self-managed-software',
      ],
    },
  ],
};

export default sidebars;
