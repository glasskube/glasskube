import React from 'react';
import styles from './styles.module.css';
import {useBlogPost} from '@docusaurus/plugin-content-blog/client';
import useBaseUrl from '@docusaurus/useBaseUrl';
import BlogPostItemHeaderAuthors from '../../BlogPostItem/Header/Authors';
import BlogPostItemHeaderInfo from '../../BlogPostItem/Header/Info';
import Link from '@docusaurus/Link';
import {Props} from '@theme/BlogPostItem';

export function LatestBlogPostItem({children}: Props) {
  const post = useBlogPost();
  return (
    <div className={styles.container}>
      <img
        src={useBaseUrl(post.frontMatter.image)}
        className={styles.image}
        alt=""
      />
      <div className={styles.metaContainer}>
        <h2 className={styles.title}>{post.metadata.title}</h2>

        <div className={styles.info}>
          <BlogPostItemHeaderAuthors />
          <BlogPostItemHeaderInfo />
        </div>

        <div className={styles.content}>{children}</div>
        <Link to={post.metadata.permalink}>Read more {`->`}</Link>
      </div>
    </div>
  );
}
