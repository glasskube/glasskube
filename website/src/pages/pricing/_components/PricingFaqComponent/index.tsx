import React from 'react';
import {PricingFAQ, PricingFAQs} from '@site/src/data/pricing';
import clsx from 'clsx';

import styles from './styles.module.css';


function FAQ({faq}: { faq: PricingFAQ }) {
  return (
    <>
      <div id={faq.id} className={clsx('margin-bottom--xl', styles.marginTop50)}/>
      <div className="card">
        <div className="card__header">
          <h3 className="anchor">{faq.question}
            <a className="hash-link" href={'#' + faq.id}>​</a>
          </h3>
        </div>
        <div className="card__body">
          <p dangerouslySetInnerHTML={{__html: faq.answer}}>
          </p>
        </div>
      </div>
    </>
  )
}


function PricingFaq() {
  return (
    <section>
      <div className="container">
        <div className="row">
          <div className="col">
            <h2 className="text--center">FAQ</h2>
            {PricingFAQs.map((faq) => (
              <FAQ key={faq.id} faq={faq}/>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}

export default React.memo(PricingFaq);
