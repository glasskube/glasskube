import Layout from '@theme/Layout';
import React from 'react';
import Heading from '@theme/Heading';
import styles from './index.module.css';
import PackageManagerScreenshots from '@site/src/pages/products/package-manager/_components/PackageManagerScreenshots';
import Link from '@docusaurus/Link';

function PackageManagerHeader() {
  return (
    <section className="margin-top--lg margin-bottom--lg text--center">
      <Heading as="h1" className={styles.heroHeading}>
        Glasskube Package Manager
      </Heading>
      <p>
        <strong>Use our GUI, CLI or GitOps integration and get started for free on GitHub.</strong>
      </p>
      <Link
        className="button button--secondary"
        to="https://github.com/glasskube/glasskube/">
        Get started on GitHub
      </Link>
    </section>
  );
}

export default function Home(): JSX.Element {
  return (
    <Layout
      title="Glasskube Package Manager"
      description="The next generation Package Manager for Kubernetes">
      <main className="margin-vert--lg">
        <PackageManagerHeader />
        <PackageManagerScreenshots />
      </main>
    </Layout>
  );
}
