import Layout from '@theme/Layout';
import React from 'react';
import HubFeatures from '@site/src/pages/products/hub/_components/HubFeatures';
import Heading from '@theme/Heading';
import DemoButton from '@site/src/components/buttons/DemoButton';
import styles from './index.module.css';
import DefaultCTA from '@site/src/components/cta/DefaultCTA/defaultCTA';
import Testimonials from '@site/src/components/Testimonials';
import NewsletterSignup from '@site/src/components/NewsletterSignup';

function HubHeader() {
  return (
    <section className="margin-top--lg margin-bottom--lg text--center">
      <Heading as="h1" className={styles.heroHeading}>
        Private Package Repositories for Kubernetes
      </Heading>
      <DemoButton additionalClassNames={'light button--lg margin-top--lg'} />
    </section>
  );
}

export default function Home(): JSX.Element {
  return (
    <Layout
      title="Glasskube Hub - Private Glasskube Package Repository for Kubernetes"
      description="Distribute Glasskube packages for Kubernetes via private hoste Glasskube package repositories">
      <main className="margin-vert--lg">
        <HubHeader />
        <HubFeatures />
        <div className="container">
          <div className="row">
            <div className="col col--10 col--offset-1">
              <DefaultCTA />
            </div>
          </div>
        </div>
        <Testimonials />
        <div className="container">
          <div className="row">
            <div className="col col--10 col--offset-1">
              <DefaultCTA />
            </div>
          </div>
        </div>
        <NewsletterSignup />
      </main>
    </Layout>
  );
}
