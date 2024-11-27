import React from 'react';
import DemoButton from '../buttons/DemoButton';
import styles from './styles.module.css';
import clsx from 'clsx';
import SignupForWaitlistButton from '@site/src/components/buttons/SignupForWaitlistButton';

export default function BlogCTA(): JSX.Element {
  return (
    <div
      className={clsx('margin-top--xl margin-bottom--lg', styles.ctaWrapper)}>
      <div className="container">
        <div className="row">
          <div className="col col--10 col--offset-1">
            <div className={styles.ctaContent}>
              <h2>Need help running Kubeflow in production?</h2>
              <p className="margin-bottom--lg">
                Book a call with us to discuss your Kubernetes deployment and
                scaling needs for your Cloud-Native AI.
              </p>
              <div className={styles.buttonContainer}>
                <SignupForWaitlistButton additionalClassNames="margin-right--md" />
                <DemoButton additionalClassNames="" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
