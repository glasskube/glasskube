import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';
import React from 'react';
import Typewriter from 'typewriter-effect';
import BrowserOnly from '@docusaurus/BrowserOnly';
import AsciinemaPlayer from '../components/asciinema-player';


function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();

  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <div className="row row--no-gutters">
          <div className={clsx('col', styles.heroCol)}>
            <Heading as="h1" className="hero__title">
              <pre>
                <Typewriter
                  onInit={(typewriter) => {
                    typewriter
                      .changeDeleteSpeed(25)
                      .changeDelay(75)
                      .typeString('brew install <span class="typewriter-command">glasskube/tap/glasskube</span>')
                      .pauseFor(1500)
                      .deleteAll(25)
                      .typeString('glasskube install ')
                      .typeString('<span class="typewriter-command">cert-manager</span>')
                      .pauseFor(1500)
                      .deleteChars('cert-manager'.length)
                      .typeString('<span class="typewriter-command">ingress-nginx</span>')
                      .pauseFor(1500)
                      .deleteChars('ingress-nginx'.length)
                      .typeString('<span class="typewriter-command">kubernetes-dashboard</span>')
                      .deleteChars('kubernetes-dashboard'.length)
                      .typeString('<span class="typewriter-command">[your-package]</span>')
                      .start();
                  }}
                />
              </pre>
            </Heading>
            <p className="hero__subtitle">{siteConfig.tagline}</p>
            <div className={styles.buttons}>
              <Link
                className="button button--secondary button--lg"
                to="/docs/">
                🚀 Get started
              </Link>
              <Link
                className="button button--outline button--lg"
                to="https://discord.gg/SxH6KUCGH7">
                Join the community
              </Link>
            </div>
            <div className={styles.producthunt}>
              <a href="https://www.producthunt.com/products/glasskube?utm_source=badge-featured&utm_medium=badge&utm_souce=badge-glasskube" target="_blank">
                <img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=452879&theme=light"
                  alt="Glasskube - &#0032;🧊&#0032;The&#0032;next&#0032;generation&#0032;Package&#0032;Manager&#0032;for&#0032;Kubernetes&#0032;📦 | Product Hunt"
                  style={{width: '250px', height: '54px'}}/>
              </a>
            </div>
          </div>
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
              Learn how to install cert-manager in less than 2 minutes using Glasskube
            </Heading>
            <AsciinemaPlayer
              src='/cast/634355.cast'
              rows='22'
              idleTimeLimit={7}
              poster='npt:0:19'
              controls={false}/>
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
        <HomepageVideo/>
        <HomepageNewsletter/>
      </main>
    </Layout>
  );
}
