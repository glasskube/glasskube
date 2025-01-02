import {translate} from '@docusaurus/Translate';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import React from 'react';
import clsx from 'clsx';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faChalkboardUser, faPhone} from '@fortawesome/free-solid-svg-icons';
import DemoButton from '@site/src/components/buttons/DemoButton';
import ContactForm from '@site/src/pages/contact/ContactForm';
import Link from '@docusaurus/Link';

const TITLE = translate({message: 'Get in touch'});
const DESCRIPTION = translate({
  message:
    "Want to get in touch? We'd love to hear from you. Here's how you can reach us.",
});

function ContactOptions() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          <div className={clsx('col', 'margin-top--lg')}>
            <div className={styles.iconBorder}>
              <div className={clsx('text--center')}>
                <FontAwesomeIcon
                  icon={faPhone}
                  size="4x"
                  className={styles.iconHeight}
                />
              </div>
              <div className="text--center padding-horiz--md margin-top--lg">
                <Heading as="h3" className="">
                  Talk to Founders
                </Heading>
                <p>Interested in Glasskube? Talk to the founders directly.</p>
              </div>
            </div>
            <div>
              <div className="row">
                <div className="col text--center">
                  <b>Glasskube, Inc.</b>
                  <br />
                  Chicago / United States
                  <br />
                  <Link href="tel:%2B18723092840">+1 (872) 309-2840</Link>
                </div>
                <div className="col col--6 text--center">
                  <b>Glasskube Labs GmbH</b>
                  <br />
                  Vienna / Austria
                  <br />
                  <Link href="tel:%2B4367761038250">+43 677 610 382 50</Link>
                </div>
              </div>
            </div>
          </div>
          <div className={clsx('col col--6', 'margin-top--lg')}>
            <div className={styles.iconBorder}>
              <div className={clsx('text--center')}>
                <FontAwesomeIcon
                  icon={faChalkboardUser}
                  size="4x"
                  className={styles.iconHeight}
                />
              </div>
              <div className="text--center padding-horiz--md margin-top--lg">
                <Heading as="h3" className="">
                  Book a demo
                </Heading>
                <p>Directly schedule a demo of Glasskube Cloud.</p>
              </div>
            </div>
            <div className="text--center">
              <DemoButton additionalClassNames={'button--lg'}></DemoButton>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

function ContactHeader() {
  return (
    <section className="margin-top--lg text--center">
      <Heading as="h1">{TITLE}</Heading>
      <p>
        <strong>{DESCRIPTION}</strong>
      </p>
    </section>
  );
}

export default function PackagePage(): JSX.Element {
  return (
    <Layout title={TITLE} description={DESCRIPTION}>
      <main className="margin-vert--lg">
        <ContactHeader />
        <ContactOptions />
        <ContactForm />
      </main>
    </Layout>
  );
}
