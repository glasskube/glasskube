import React from 'react';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import styles from './styles.module.css';

interface GlossaryItemProps {
  term: string;
}

export default function GlossaryItem({ term }: GlossaryItemProps): JSX.Element {
  const {siteConfig} = useDocusaurusContext();
  const glossaryItems = siteConfig.customFields.glossaryItems as Array<{term: string, fileName: string}>;
  
  const item = glossaryItems.find(item => item.term === term);
  const link = item ? `/docs/glossary/${item.fileName}` : '#';

  return (
    <div className={styles.glossaryItem}>
      <div className={styles.glossaryItemInner}>
        <h3>{term}</h3>
        <p className={styles.description}>
          {/* Add a brief description here if available */}
        </p>
        <Link to={link} className={styles.moreButton}>
          More
        </Link>
      </div>
    </div>
  );
}