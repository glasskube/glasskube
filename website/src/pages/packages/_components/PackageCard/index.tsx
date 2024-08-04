import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import Image from '@theme/IdealImage';
import {
  type Package,
  Tag,
  TagList,
  Tags,
  type TagType,
} from '@site/src/data/packages';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import {sortBy} from '@site/src/utils/jsUtils';

const TagComp = React.forwardRef<HTMLLIElement, Tag>(
  ({label, color}: {label: string; color: string}, ref) => (
    <li ref={ref} className={styles.tag}>
      <span className={styles.textLabel}>{label.toLowerCase()}</span>
      <span className={styles.colorLabel} style={{backgroundColor: color}} />
    </li>
  ),
);
TagComp.displayName = 'TagComp';

function PackageTag({tags}: {tags: TagType[]}) {
  const tagObjects = tags.map(tag => ({tag, ...Tags[tag]}));

  // Keep same order for all tags
  const tagObjectsSorted = sortBy(tagObjects, tagObject =>
    TagList.indexOf(tagObject.tag),
  );

  return (
    <>
      {tagObjectsSorted.map((tagObject, index) => {
        return <TagComp key={index} {...tagObject} />;
      })}
    </>
  );
}

function PackageCard({user}: {user: Package}) {
  return (
    <li key={user.name} className="card shadow--tl">
      <div className={clsx('card__image', styles.packageCardImage)}>
        <Image img={user.iconUrl} alt={user.name} />
      </div>
      <div className="card__body">
        <div className={clsx(styles.packageCardHeader)}>
          <Heading as="h4" className={styles.packageCardTitle}>
            <Link href={user.websiteUrl}>{user.name}</Link>
          </Heading>
          <Link
            href={user.sourceUrl}
            className={clsx(
              'button button--secondary button--sm',
              styles.packageCardSrcBtn,
            )}>
            source code
          </Link>
        </div>
        <p className={styles.packageCardBody}>{user.shortDescription}</p>
      </div>
      <ul className={clsx('card__footer', styles.cardFooter)}>
        <PackageTag tags={user.tags} />
      </ul>
    </li>
  );
}

export default React.memo(PackageCard);
