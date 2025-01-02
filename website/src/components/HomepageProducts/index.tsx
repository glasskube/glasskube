import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import React from 'react';
import ThemedImage from '@theme/ThemedImage';
import Link from '@docusaurus/Link';

type ScreenshotItem = {
  title: string;
  description: JSX.Element;
  screenshotAltText: string;
  lightScreenshotUrl: string;
  darkScreenshotUrl: string;
};

const ScreenshotList: ScreenshotItem[] = [
  {
    title: 'Software Distributor Platform',
    description: (
      <>
        Onboard customers faster, distribute updates easily and troubleshoot
        issues with confidence.
        <br />
        <Link
          className="button button--secondary margin-top--md"
          to="https://glasskube.cloud/">
          Get started free
        </Link>
      </>
    ),
    screenshotAltText: 'Software distribution platform',
    lightScreenshotUrl:
      '/img/screenshots/glasskube-cloud-software-distributor-platform-light.png',
    darkScreenshotUrl:
      '/img/screenshots/glasskube-cloud-software-distributor-platform-dark.png',
  },
  {
    title: 'Customer Portal',
    description: (
      <>
        Give your customers a simple, but powerful portal to simplify their
        installations and stay on top of their deployments.
        <br />
        <Link
          className="button button--secondary margin-top--md"
          to="https://glasskube.cloud/">
          Get started free
        </Link>
      </>
    ),
    screenshotAltText: 'Glasskube overview page',
    lightScreenshotUrl:
      '/img/screenshots/glasskube-cloud-customer-portal-light.png',
    darkScreenshotUrl:
      '/img/screenshots/glasskube-cloud-customer-portal-dark.png',
  },
];

function Screenshot(item: ScreenshotItem) {
  return (
    <>
      <div className={clsx('col col--3')}>
        <div className={styles.sticky}>
          <Heading as="h3" className={styles.borderLeft}>
            {item.title}
          </Heading>
          <p>{item.description}</p>
        </div>
      </div>
      <div className={clsx('col col--9')}>
        <ThemedImage
          alt={item.screenshotAltText}
          sources={{
            light: item.lightScreenshotUrl,
            dark: item.darkScreenshotUrl,
          }}
        />
      </div>
    </>
  );
}

export default function HomepageProducts(): JSX.Element {
  return (
    <section className={styles.screenshots}>
      <div className="container margin-top--lg">
        <div className="row">
          <div className="col text--center">
            <Heading as="h2">
              From onboarding to update distribution, and support—all in one
              platform.
            </Heading>
          </div>
        </div>
        {ScreenshotList.map((item, idx) => (
          <div className="row" key={idx}>
            <Screenshot {...item} />
          </div>
        ))}
      </div>
    </section>
  );
}
