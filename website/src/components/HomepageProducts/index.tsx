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
    title: 'Glasskube Package Manager',
    description: (
      <>
        Use our GUI, CLI or GitOps integration and get started for free on
        GitHub. <br />
        <Link
          className="button button--secondary margin-top--md"
          to="/products/package-manager">
          Learn more
        </Link>
      </>
    ),
    screenshotAltText: 'Glasskube overview page',
    lightScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/d3e764e4-7ee0-4ad8-b0eb-d0ebf219160f',
    darkScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/323994d6-6b08-4dca-ac59-d29ae6b37f94',
  },
  {
    title: 'Glasskube Hub',
    description: (
      <>
        The easiest way to create, manage, and deliver your Cloud Native
        applications. <br />
        <Link
          className="button button--secondary margin-top--md"
          to="/products/hub">
          Learn more
        </Link>
      </>
    ),
    screenshotAltText: 'Glasskube settings page',
    lightScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/e8565cf7-6089-4b82-b169-91f5b3ef4c33',
    darkScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/cf1f1983-78c8-4bb4-9d47-86dbf3a16c4e',
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
              Let Glasskube take care about your Kubernetes infrastructure
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
