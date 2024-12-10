export const PricingFAQs: PricingFAQ[] = [
  {
    id: 'glasskube-cloud',
    question: 'What is Glasskube Cloud?',
    answer:
      'Glasskube Cloud is the easiest way for ISVs to distribute their software to their customers. It has multiple modules: E.g. the Software Distribution Platform and Customer Portal',
  },
  {
    id: 'software-distribution-platform',
    question: 'What is the Glasskube Cloud Software Distribution Platform?',
    answer:
      'A platform for ISVs to manage their software application, release channels, versions, customers, and their deployment targets.',
  },
  {
    id: 'customer-portal',
    question: 'What is the Glasskube Cloud Customer Portal?',
    answer:
      'The customer portal is a white labeled or co-branded portal for your customers to receive installation instructions, updates, and support.',
  },
  {
    id: 'application',
    question: 'What is an application?',
    answer:
      'An ISV might distribute multiple applications. E.g. a monitoring platform and data collection agent.',
  },
  {
    id: 'release-channel',
    question: 'What is a release channel?',
    answer:
      'An application can have multiple release channels. E.g. a stable and a beta channel. Each channel can have multiple versions.' +
      ' A channel is assigned to a deployment target.',
  },
  {
    id: 'deployment-target',
    question: 'What is a deployment target?',
    answer:
      'Your customers can have multiple deployment targets. Deployment targets are usually associated with a release channel.',
  },
];

export type PricingFAQ = {
  id: string;
  question: string;
  answer: string;
};
