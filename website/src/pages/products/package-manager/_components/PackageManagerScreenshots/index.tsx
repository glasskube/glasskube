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
    title: 'Focus on simplicity and reliability',
    description: (
      <>
        After <Link href={'/docs/getting-started/install/'}>installation</Link>,
        use <code>glasskube serve</code> to browse and find all your favorite
        packages in one place. You do not need to look for a Helm repository to
        find a specific package.
      </>
    ),
    screenshotAltText: 'Glasskube overview page',
    lightScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/d3e764e4-7ee0-4ad8-b0eb-d0ebf219160f',
    darkScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/323994d6-6b08-4dca-ac59-d29ae6b37f94',
  },
  {
    title: 'Package configurations',
    description: (
      <>
        <Link href={'/docs/design/package-config/'}>Configure packages</Link>{' '}
        with typesafe input values via the UI and with an interactive questioner
        via the cli. Your also able to inject values from other packages,
        ConfigMaps and Secrets. Say Good-Bye to un-typed and un-documented{' '}
        <code>values.yaml</code> files.
      </>
    ),
    screenshotAltText: 'Glasskube package configuration page',
    lightScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/28ceea95-c66d-4f62-8fe4-d4d1be160ad6',
    darkScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/df6bd7d4-7cac-435b-b3a0-31c3cab6069b',
  },
  {
    title: 'Dependency Management',
    description: (
      <>
        Glasskube{' '}
        <Link href={'/docs/design/dependency-management/'}>
          packages are dependency aware
        </Link>
        , so they can be used and referenced by multiple other packages. They
        will also get installed in the correct namespace. This is how umbrella
        charts should have worked from the beginning.
      </>
    ),
    screenshotAltText: 'Glasskube dependency page',
    lightScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/b7d65ac6-de61-4771-b81c-29b7b1926f77',
    darkScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/9588b3fc-2a87-454e-97ff-b0f7558717bc',
  },
  {
    title: 'Safe Package updates',
    description: (
      <>
        Preview and perform pending Updates to your desired version with a
        single click of a button. All updates have been tested the Glasskube
        test suite before being available in the public respository.
      </>
    ),
    screenshotAltText: 'Glasskube package updates',
    lightScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/70dd3e59-0a3e-46ee-b193-6e91159ab9d4',
    darkScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/a6e6dc72-9919-4d15-addf-bc709ec76d9d',
  },
  {
    title: 'Reactions and comments',
    description: (
      <>
        Discuss and upvote your favorite Kubernetes package on{' '}
        <Link href="https://github.com/glasskube/glasskube/discussions/categories/packages">
          GitHub
        </Link>{' '}
        or right inside the Glasskube UI.
      </>
    ),
    screenshotAltText: 'Glasskube package reactions',
    lightScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/538419ab-3852-4342-9020-0103da00fc21',
    darkScreenshotUrl:
      'https://github.com/glasskube/glasskube/assets/3041752/56f08373-fbbe-46fd-820e-fb637114336b',
  },
  {
    title: 'Multiple Repositories and private packages',
    description: (
      <>
        Use{' '}
        <Link href={'/docs/design/repositories/'}>multiple repositories</Link>{' '}
        and publish your own private packages. This could be your companies
        Internal services packages, so all developers will have the up-to-date
        and easily configured internal services.
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

export default function PackageManagerScreenshots(): JSX.Element {
  return (
    <section className={styles.screenshots}>
      <div className="container margin-top--lg">
        <div className="row">
          <div className="col text--center">
            <Heading as="h2">
              A Package Manager built with Developer Experience in mind
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
