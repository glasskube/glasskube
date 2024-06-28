import clsx from 'clsx';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';
import React from 'react';
import Typewriter from 'typewriter-effect';
import BrowserOnly from '@docusaurus/BrowserOnly';
import HomepageScreenshots from '@site/src/components/HomepageScreenshots';
import TalkToFoundersButton from '@site/src/components/buttons/TalkToFoundersButton';
import SignupForWaitlistButton from '@site/src/components/buttons/SignupForWaitlistButton';
import useBaseUrl from '@docusaurus/core/lib/client/exports/useBaseUrl';
import Image from '@theme/IdealImage';


function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();

  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <div className="row row--no-gutters">
          <div className="col">
            <Heading as="h1" className={styles.heroTitle}>
              {siteConfig.tagline}
            </Heading>
          </div>

        </div>
        <div className="row row--no-gutters">
          <div className={clsx('col', styles.heroCol)}>
            <div className={styles.buttons}>
              <TalkToFoundersButton additionalClassNames={'button--lg light'}/>
              <SignupForWaitlistButton additionalClassNames={'button--lg'}/>
            </div>
            <div className={styles.yc}>
              <h4>Backed by</h4>
              <a
                href="https://www.ycombinator.com/companies/glasskube"
                target="_blank">
                <Image
                  alt="Glasskube backed by Y Combinator"
                  className={styles.ycImg}
                  img={useBaseUrl('/img/yc/yc.svg')}/>
              </a>
            </div>
          </div>
        </div>
        <div className="row row--no-gutters">
          <div className="col">
            <div className={styles.lottiePlayerWrapper}>
              <BrowserOnly fallback={<div className={styles.lottieFallback}>Loading...</div>}>
                {() => {
                  const Player = require('@lottiefiles/react-lottie-player').Player;
                  return <Player
                    autoplay
                    loop
                    src="/animations/home.json"
                    style={{height: '480px'}}
                  />
                }}
              </BrowserOnly>
            </div>
          </div>
        </div>
        <div className="row row--no-gutters">
          <div className={clsx('col', styles.heroCol, styles.typewriter)}>
              <pre>
                <Typewriter
                  onInit={(typewriter) => {
                    typewriter
                      .changeDeleteSpeed(25)
                      .changeDelay(75)
                      .typeString('brew <span class="typewriter-command">install</span> <span class="typewriter-argument">glasskube/tap/glasskube</span>')
                      .pauseFor(1500)
                      .deleteAll(25)
                      .typeString('glasskube <span class="typewriter-command">install </span>')
                      .typeString('<span class="typewriter-argument">cert-manager</span>')
                      .pauseFor(1500)
                      .deleteChars('cert-manager'.length)
                      .typeString('<span class="typewriter-argument">ingress-nginx</span>')
                      .pauseFor(1500)
                      .deleteChars('ingress-nginx'.length)
                      .typeString('<span class="typewriter-argument">kubernetes-dashboard</span>')
                      .deleteChars('kubernetes-dashboard'.length)
                      .typeString('<span class="typewriter-argument">[your-package]</span>')
                      .start();
                  }}
                />
              </pre>
          </div>
        </div>

      </div>
    </header>
  );
}

function HomepageVideo() {
  const {siteConfig} = useDocusaurusContext();

  return (
    <div className={clsx('container-fluid', 'text--center', styles.backgroundSecondary)}>
      <div className="container text--center">
        <div className="row">
          <div className="col col--8 col--offset-2 margin-vert--lg">
            <Heading as={'h2'} className={styles.colorWhite}>
              Learn how to use Glasskube in less than 2 minutes
            </Heading>
            <iframe width="100%" height="460" src="https://www.youtube-nocookie.com/embed/aIeTHGWsG2c?si=KUcqvY4coU89GmdK" title="YouTube video player" frameBorder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerPolicy="strict-origin-when-cross-origin" allowFullScreen></iframe>
          </div>
        </div>
      </div>
    </div>

  );
}

function loadScript() {
  if (typeof window === "undefined") {
    return null;
  }

  const elementId = "hs-script";
  if (document.getElementById(elementId) === null) {
    const script = document.createElement("script");
    script.type = "text/javascript";
    script.src = "https://js-eu1.hs-scripts.com/144345473.js";
    script.id = elementId;
    document.head.appendChild(script);
  }
}

function HomepageNewsletter() {
  const {siteConfig} = useDocusaurusContext();

  return (
    <div className="container text--center">
      <div className="row">
        <div className="col col--6 col--offset-3 margin-vert--lg">
          <div>
            <Heading as={'h2'}>
              Glasskube Newsletter
            </Heading>
            <p>Sign-Up to get the latest product updates and release notes!</p>
            <NewsletterForm/>
          </div>
        </div>
      </div>
    </div>

  );
}

class NewsletterForm extends React.Component<any, { value: string }> {
  constructor(props) {
    super(props);
    this.state = {value: ''};

    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  handleChange(event) {
    if (!this.state.value) {
      loadScript();
    }
    this.setState({value: event.target.value});
  }

  async handleSubmit(event) {
    event.preventDefault()
    await fetch('https://cms.glasskube.eu/api/ezforms/submit', {
      method: 'POST',
      body: JSON.stringify({
        token: '',
        formName: 'newsletter',
        formData: {
          email: this.state.value,
        }
      }),
      headers: {
        'Content-type': 'application/json; charset=UTF-8',
      },
    })
      .then((response) => response.json())
      .then((data) => {
        alert('Email successfully added to our newsletter: ' + this.state.value);
      })
      .catch((err) => {
        console.error(err.message);
      });
  }

  render() {
    return (
      <form onSubmit={this.handleSubmit} className={styles.newsletterForm}>
        <input type="email" id="email" name="email" required
               placeholder="your-email@corp.com"
               value={this.state.value} onChange={this.handleChange}
               className={styles.emailInput}/>
        <button
          className="button button--secondary button--lg"
          type="submit">
          Subscribe
        </button>
      </form>
    );
  }
}

export default function Home(): JSX.Element {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={siteConfig.tagline}
      description='Featuring a GUI and a CLI. Glasskube packages are dependency aware, GitOps ready and get automatic updates via a central public package repository.'>
      <HomepageHeader/>
      <main>
        <HomepageFeatures/>
        <HomepageScreenshots/>
        <HomepageVideo/>
        <HomepageNewsletter/>
      </main>
    </Layout>
  );
}
