import React from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import Image from '@theme/IdealImage';

type TestimonialItem = {
  name: string;
  role: string;
  company: string;
  image: string;
  text: string;
};

const TestimonialList: TestimonialItem[] = [
  {
    name: 'Denzell Ford',
    role: 'CTO',
    company: 'Trieve (YC W24)',
    image:
      'https://github.com/user-attachments/assets/cb098c0a-d681-4755-ba88-b24e5a5daad4',
    text:
      'Glasskube is able to package Trieve for generic Kubernetes clusters and also provide specific overlays for GCP and AWS for one-line installations.' +
      ' It effortlessly manages our complex stack—Qdrant, PostgreSQL, Keycloak, Embedding servers and more—making our on-premises Kubernetes' +
      ' deployments of our Search and RAG API simple and efficient.',
  },
  {
    name: 'Mathias Nöbauer',
    role: 'CEO',
    company: 'A1 Digital',
    image:
      'https://github.com/user-attachments/assets/8e25ddbe-38f8-4dac-a519-3c901449760c',
    text:
      'Glasskube is the perfect solution for powering private infrastructure deployments. It simplifies the process, ' +
      'making it easy to deploy and manage complex applications',
  },
  {
    name: 'François Massot',
    role: 'Co-founder',
    company: 'Quickwit',
    image:
      'https://github.com/user-attachments/assets/3f94b763-b679-4cb7-9e5b-9c52381e3f10',
    text:
      'Glasskube`s intuitive interface has made deploying and managing Quickwit on Kubernetes a breeze. ' +
      'The UX is outstanding. Features like type-safe package configuration and dependency aware are a game changer ' +
      'for managing packages in Kubernetes.',
  },
  {
    name: 'Siggi Simonarson',
    role: 'Founder',
    company: 'BuildBuddy (YC W20)',
    image:
      'https://github.com/user-attachments/assets/fdae75c6-5c31-4c06-98c5-668932f8d1e2',
    text: 'Maintaining Helm charts takes a lot of time and effort. It`s time for the next generation of package manager for Kubernetes.',
  },
  {
    name: 'Chris Lo',
    role: 'Co-founder',
    company: 'Tracecat (YC W24)',
    image:
      'https://github.com/user-attachments/assets/2d61a326-e9bc-4982-a800-15c356d99459',
    text:
      'Glasskube made it possible to move from being self-hosted only on AWS ECS to being deployed everywhere ' +
      'Kubernetes runs—even in the most protected environments.',
  },
  {
    name: 'Antoine Coetsier',
    role: 'Co-founder and COO',
    company: 'EXOSCALE',
    image:
      'https://github.com/user-attachments/assets/432be015-ae42-46b1-bc14-0ae6a79a49ce',
    text:
      'Glasskube`s managed open-source tools are a crucial part of our marketplace, enabling secure, fully automated ' +
      'application deployment and management on Exoscale`s cloud.',
  },
];

function Testimonial({name, role, company, image, text}: TestimonialItem) {
  return (
    <div className={clsx('col col--6 margin-top--lg')}>
      <div className={styles.testimonial}>
        <div className={styles.testimonialContent}>
          <p className={styles.testimonialText}>&ldquo;{text}&rdquo;</p>
          <div className={styles.testimonialAuthor}>
            <Image img={image} className={styles.testimonialImage} />
            <div>
              <p className={styles.testimonialName}>{name}</p>
              <p className={styles.testimonialRole}>
                {role}, {company}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function Testimonials(): JSX.Element {
  return (
    <section className={styles.testimonials}>
      <div className="container">
        <Heading as="h2" className={clsx('text--center', styles.heading)}>
          What our users say
        </Heading>
        <div className="row">
          {TestimonialList.map((props, idx) => (
            <Testimonial key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
