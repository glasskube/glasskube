import React from 'react';
import Layout from '@theme/Layout';
import GlossaryPage from '@theme/GlossaryPage';

export default function Glossary(): JSX.Element {
  return (
    <Layout
      title="Glossary"
      description="Glasskube Glossary - Key terms and definitions">
      <GlossaryPage />
    </Layout>
  );
}