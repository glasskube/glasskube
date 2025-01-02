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
        slug: '/',
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
        'byoc-definition',
        'saas-definition',
        'on-premises-definition',
        'self-managed-software',
        'isv',
        'air-gapped',
      ],
    },
  ],
};

export default sidebars;
