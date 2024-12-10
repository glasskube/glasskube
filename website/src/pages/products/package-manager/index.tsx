import Layout from '@theme/Layout';
import React from 'react';
import Heading from '@theme/Heading';
import styles from './index.module.css';
import PackageManagerScreenshots from '@site/src/pages/products/package-manager/_components/PackageManagerScreenshots';
import Link from '@docusaurus/Link';
import clsx from 'clsx';
import Typewriter from 'typewriter-effect';
import PackageManagerFeatures from '@site/src/pages/products/package-manager/_components/PackageManagerFeatures';
import Testimonials from '@site/src/components/Testimonials';
import DefaultCTA from '@site/src/components/cta/DefaultCTA/defaultCTA';
import NewsletterSignup from '@site/src/components/NewsletterSignup';

function PackageManagerHeader() {
  return (
    <section className="margin-top--lg margin-bottom--lg text--center">
      <Heading as="h1" className={styles.heroHeading}>
        Glasskube Package Manager
      </Heading>
      <p>
        <strong>
          Use our GUI, CLI or GitOps integration and get started for free on
          GitHub.
        </strong>
      </p>
      <Link
        className="button button--secondary"
        to="https://github.com/glasskube/glasskube/">
        Get started on GitHub
      </Link>
      &nbsp;
      <Link
        className="button button--primary"
        to="https://discord.gg/p7uYfnxZFd">
        Join our Discord
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
        <section>
          <div className="container margin-top--lg">
            <div className="row row--no-gutters">
              <div className={clsx('col', styles.heroCol, styles.typewriter)}>
                <pre>
                  <Typewriter
                    onInit={typewriter => {
                      typewriter
                        .changeDeleteSpeed(25)
                        .changeDelay(75)
                        .typeString(
                          'brew <span class="typewriter-command">install</span> <span class="typewriter-argument">glasskube/tap/glasskube</span>',
                        )
                        .pauseFor(1500)
                        .deleteAll(25)
                        .typeString(
                          'glasskube <span class="typewriter-command">install </span>',
                        )
                        .typeString(
                          '<span class="typewriter-argument">cert-manager</span>',
                        )
                        .pauseFor(1500)
                        .deleteChars('cert-manager'.length)
                        .typeString(
                          '<span class="typewriter-argument">ingress-nginx</span>',
                        )
                        .pauseFor(1500)
                        .deleteChars('ingress-nginx'.length)
                        .typeString(
                          '<span class="typewriter-argument">kubernetes-dashboard</span>',
                        )
                        .deleteChars('kubernetes-dashboard'.length)
                        .typeString(
                          '<span class="typewriter-argument">[your-package]</span>',
                        )
                        .start();
                    }}
                  />
                </pre>
              </div>
            </div>
          </div>
        </section>
        <PackageManagerFeatures />
        <PackageManagerScreenshots />
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
