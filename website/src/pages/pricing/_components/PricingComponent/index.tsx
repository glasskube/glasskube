import Link from '@docusaurus/Link';
import DemoButton from '@site/src/components/buttons/DemoButton';
import SignupForWaitlistButton from '@site/src/components/buttons/SignupForWaitlistButton';
import React from 'react';
import styles from './styles.module.css';
import ContactSalesButton from '@site/src/components/buttons/ContactSalesButton';
import clsx from 'clsx';

function Pricing() {
  return (
    <section>
      <div className="container">
        <div className="row">
          <div className="col col--4">
            <div className={clsx('card', 'shadow--md', styles.pricingCardSide)}>
              <div
                className={clsx(
                  'card__header',
                  'text--center',
                  styles.pricingCardHeader,
                  styles.pricingCardHeaderSide,
                )}>
                <h3>Free Forever</h3>
                (no Credit Card required)
              </div>
              <hr className={styles.hr} />
              <div className="card__body">
                <h4 className={styles.pricingSectionHeader}>
                  Easy Customer Onboarding
                </h4>
                <ul className={styles.pricingItemList}>
                  <li>Application & Customer Management</li>
                  <li>
                    Co-branded Customer Portal with interactive installation
                    instructions
                  </li>
                </ul>
                <h4 className={styles.pricingSectionHeader}>
                  Software Distribution
                </h4>
                <ul>
                  <li>Support for docker-compose based application</li>
                  <li>Support for Helm based applications</li>
                  <li>Customer configurations & secrets</li>
                  <li>Release Channels</li>
                </ul>
                <h4 className={styles.pricingSectionHeader}>
                  Monitoring & Customer Support
                </h4>
                <ul>
                  <li>Deployment Target Healthcheck</li>
                  <li>Installed versions overview</li>
                </ul>
              </div>
              <div className="card__footer">
                <Link
                  className="button button--secondary button--block button--lg"
                  to="https://glasskube.cloud/">
                  Free forever
                </Link>
              </div>
            </div>
          </div>
          <div className="col col--4">
            <div className={clsx('card', 'shadow--md', styles.pricingCardMain)}>
              <div
                className={clsx(
                  'card__header',
                  'text--center',
                  styles.pricingCardHeader,
                  styles.pricingCardHeaderMain,
                )}>
                <h3>Pro</h3>
                Coming soon
              </div>
              <hr className={styles.hr} />
              <div className="card__body">
                <h4 className={styles.pricingSectionHeader}>
                  Easy Customer Onboarding
                </h4>
                All Free features and:
                <ul className={styles.pricingItemList}>
                  <li>API Integration</li>
                  <li>License dependent features</li>
                </ul>
                <h4 className={styles.pricingSectionHeader}>
                  Software Distribution
                </h4>
                All Free features and:
                <ul>
                  <li>Support for Glasskube Kubernetes packages</li>
                  <li>
                    Pre-flight checks to ensure deployment target integrity
                  </li>
                </ul>
                <h4 className={styles.pricingSectionHeader}>
                  Monitoring & Customer Support
                </h4>
                All Free features and:
                <ul>
                  <li>Customers can collect and send support insights</li>
                </ul>
              </div>
              <div className="card__footer">
                <DemoButton additionalClassNames={'button--lg button--block'} />
                <div className="margin-top--md" />
                <SignupForWaitlistButton
                  additionalClassNames={'button--lg button--block'}
                />
              </div>
            </div>
          </div>
          <div className="col">
            <div className={clsx('card', 'shadow--md', styles.pricingCardSide)}>
              <div
                className={clsx(
                  'card__header',
                  'text--center',
                  styles.pricingCardHeader,
                  styles.pricingCardHeaderSide,
                )}>
                <h3>Enterprise</h3>
              </div>
              <hr className={styles.hr} />
              <div className="card__body">
                <h4 className={styles.pricingSectionHeader}>
                  Easy Customer Onboarding
                </h4>
                All Free & Pro features and:
                <ul className={styles.pricingItemList}>
                  <li>Single-Sign-On</li>
                  <li>White Label Customer Portal</li>
                </ul>
                <h4 className={styles.pricingSectionHeader}>
                  Software Distribution
                </h4>
                All Free & Pro features and:
                <ul>
                  <li>Software vulnerability scanning & reporting</li>
                  <li>Air-Gapped software distribution</li>
                  <li>Scheduled updates</li>
                </ul>
                <h4 className={styles.pricingSectionHeader}>
                  Monitoring & Customer Support
                </h4>
                All Free & Pro features and:
                <ul>
                  <li>Customer support SLAs</li>
                </ul>
                <h4 className={styles.pricingSectionHeader}>On-premises</h4>
                You can also self host Glasskube Cloud on your own
                infrastructure.
              </div>
              <div className="card__footer">
                <ContactSalesButton
                  additionalClassNames={'button--lg button--block'}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

export default React.memo(Pricing);
