import Link from '@docusaurus/Link';
import DemoButton from '@site/src/components/buttons/DemoButton';
import SignupForWaitlistButton from '@site/src/components/buttons/SignupForWaitlistButton';
import React from 'react';
import styles from './styles.module.css';
import ContactSalesButton from '@site/src/components/buttons/ContactSalesButton';

function Pricing() {
  return (
    <section>
      <div className="container">
        <div className="row">
          <div className="col col--4">
            <div className="card margin-top--lg shadow--md">
              <div className="card__header text--center">
                <h3>Open Source Kubernetes Package Manager</h3>
                <h4>Distribute your community edition for free</h4>
              </div>
              <hr />
              <div className="card__body">
                <p>
                  Get your open-source project ready for Kubernetes and easily
                  deploy across clusters.
                </p>
                <ul className={styles.pricingItemList}>
                  <li>Never maintain a Helm charts every again</li>
                  <li>
                    Compatibility across different Kubernetes environments
                  </li>
                  <li>Installation and configuration via CLI and UI</li>
                  <li>
                    Utilize managed dependencies like PostgreSQL, Redis and many
                    more
                  </li>
                  <li>Distribute updates and patches via Glasskube</li>
                  <li>GitHub & Discord community support</li>
                </ul>
              </div>
              <div className="card__footer">
                <Link
                  className="button button--secondary button--block button--lg"
                  to="https://discord.gg/p7uYfnxZFd">
                  Submit on Discord
                </Link>
              </div>
            </div>
          </div>
          <div className="col col--4">
            <div className="card shadow--md">
              <div className="card__header text--center">
                <h3>Application Delivery Platform for Kubernetes</h3>
                First 5 customer deployments included
                <br />
                <h4>Starting at $18k Platform fee per year</h4>
              </div>
              <hr />
              <div className="card__body">
                <p>
                  Distribute your commercial software into customer managed
                  Kubernetes environments in private infrastructure, different
                  cloud environments, and on-premises (even air-gapped).
                </p>
                <p>Everything from the free plan and:</p>
                <ul className={styles.pricingItemList}>
                  <li>Integrated license management</li>
                  <li>Insights and Reporting</li>
                  <li>Advanced Compatibility Testing</li>
                  <li>Resolve Technical Issues with Playbooks</li>
                  <li>BYOC</li>
                </ul>
              </div>
              <div className="card__footer">
                <SignupForWaitlistButton
                  additionalClassNames={'button--lg button--block'}
                />
                <div className="margin-top--md" />
                <DemoButton additionalClassNames={'button--lg button--block'} />
              </div>
            </div>
          </div>
          <div className="col">
            <div className="card col--4 margin-top--lg shadow--md">
              <div className="card__header text--center">
                <h3>Enterprise hardened Kubernetes Packages</h3>
                <h4>Custom quote</h4>
              </div>
              <hr />
              <div className="card__body">
                <p>
                  Tailor-made solutions for the highest level of cloud security
                  and reliability and personal support.
                </p>
                <p>Glasskube offers:</p>
                <ul className={styles.pricingItemList}>
                  <li>Free DevOps assessment call</li>
                  <li>SLAs and personal onboarding</li>
                  <li>Dedicated success manager</li>
                </ul>
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
