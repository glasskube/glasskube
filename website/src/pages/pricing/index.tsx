import {translate} from '@docusaurus/Translate';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';
import Pricing from './_components/PricingComponent';
import PricingFaq from '@site/src/pages/pricing/_components/PricingFaqComponent';

const TITLE = translate({message: 'Glasskube Pricing'});
const DESCRIPTION = translate({
  message: 'Choose the Right Plan for Your Cloud Security and DevOps Needs',
});

function PricingHeader() {
  return (
    <section className="margin-top--lg margin-bottom--lg text--center">
      <Heading as="h1">{TITLE}</Heading>
      <p>{DESCRIPTION}</p>
    </section>
  );
}

export default function PackagePage(): JSX.Element {
  return (
    <Layout title={TITLE} description={DESCRIPTION}>
      <main className="margin-vert--lg">
        <PricingHeader />
        <Pricing />
        <div className="margin-vert--xl" />
        <PricingFaq />
        <div className="margin-vert--md" />
      </main>
    </Layout>
  );
}
