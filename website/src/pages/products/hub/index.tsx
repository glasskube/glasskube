import Layout from '@theme/Layout';
import React from 'react';
import HubFeatures from '@site/src/pages/products/hub/_components/HubFeatures';
import Heading from '@theme/Heading';
import DemoButton from '@site/src/components/buttons/DemoButton';
import styles from './index.module.css';

function HubHeader() {
  return (
    <section className="margin-top--lg margin-bottom--lg text--center">
      <Heading as="h1" className={styles.heroHeading}>
        Application Delivery Platform for Kubernetes
      </Heading>
      <DemoButton additionalClassNames={'light button--lg margin-top--lg'} />
    </section>
  );
}

export default function Home(): JSX.Element {
  return (
    <Layout
      title="Application Delivery Platform"
      description="Sell software to enterprises and deploy to private infrastructure, different cloud environments, and
          on-premises.">
      <main className="margin-vert--lg">
        <HubHeader />
        <HubFeatures />
      </main>
    </Layout>
  );
}
