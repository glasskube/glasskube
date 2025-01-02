import React from 'react';
import DemoButton from '@site/src/components/buttons/DemoButton';
import styles from './styles.module.css';
import clsx from 'clsx';
import SignupForWaitlistButton from '@site/src/components/buttons/SignupForWaitlistButton';

export default function CustomCTA({
  header,
  text,
}: {
  header: string;
  text: string;
}): JSX.Element {
  return (
    <div
      className={clsx('margin-top--xl margin-bottom--lg', styles.ctaWrapper)}>
      <div className="container">
        <div className="row">
          <div className="col col--10 col--offset-1">
            <div className={styles.ctaContent}>
              <h2>{header}</h2>
              <p className="margin-bottom--lg">{text}</p>
              <div className={styles.buttonContainer}>
                <DemoButton additionalClassNames="button--lg" />
                <SignupForWaitlistButton additionalClassNames="button--lg" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
