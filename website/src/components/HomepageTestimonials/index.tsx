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
    name: 'John Doe',
    role: 'CTO',
    company: 'TechCorp',
    image: 'https://example.com/john-doe.jpg',
    text: 'Glasskube has revolutionized our Kubernetes package management. It\'s intuitive, efficient, and a game-changer for our DevOps team.',
  },
  {
    name: 'Jane Smith',
    role: 'Lead DevOps Engineer',
    company: 'InnoSystems',
    image: 'https://example.com/jane-smith.jpg',
    text: 'The simplicity and power of Glasskube have significantly improved our deployment processes. It\'s now an essential tool in our infrastructure.',
  },
];

function Testimonial({name, role, company, image, text}: TestimonialItem) {
  return (
    <div className={clsx('col col--6')}>
      <div className={styles.testimonial}>
        <div className={styles.testimonialContent}>
          <p className={styles.testimonialText}>"{text}"</p>
          <div className={styles.testimonialAuthor}>
            <Image img={image} className={styles.testimonialImage} />
            <div>
              <p className={styles.testimonialName}>{name}</p>
              <p className={styles.testimonialRole}>{role}, {company}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function HomepageTestimonials(): JSX.Element {
  return (
    <section className={styles.testimonials}>
      <div className="container">
        <Heading as="h2" className={clsx('text--center', styles.heading)}>
          What Our Users Say
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
