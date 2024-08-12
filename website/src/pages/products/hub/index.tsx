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
        Glasskube Hub
      </Heading>
      <p>
        <strong>
          The easiest way to create, manage, and deliver your Cloud Native
          applications.
        </strong>
      </p>
      <DemoButton additionalClassNames={'light'} />
    </section>
  );
}

export default function Home(): JSX.Element {
  return (
    <Layout
      title="Glasskube Hub"
      description="The easiest way to create, manage, and deliver your Cloud Native applications.">
      <main className="margin-vert--lg">
        <HubHeader />
        <HubFeatures />
      </main>
    </Layout>
  );
}
