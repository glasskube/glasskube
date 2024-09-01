import React from "react";
import { BlogPostProvider } from "@docusaurus/plugin-content-blog/client";
import BlogPostItem from "@theme/BlogPostItem";
import type { Props } from "@theme/BlogPostItems";
import { LatestBlogPostItem } from "./LatestBlogPostItem/LatestBlogPostItem";
import clsx from "clsx";
import styles from "./styles.module.css";

export default function BlogPostItems({
  items,
  metadata,
  component: BlogPostItemComponent = BlogPostItem,
}: Props & { metadata: any }): JSX.Element {
  const [latestBlog, ...rest] = items;
  console.log(items);
  const showLatest = metadata?.page === 1;

  return (
    <>
      {showLatest && (
        <div className={styles.latestContainer}>
          <div className="container">
            <h1 className={styles.title}>Latest</h1>
            <BlogPostProvider
              key={latestBlog.content.metadata.permalink}
              content={latestBlog.content}
            >
              <LatestBlogPostItem>{latestBlog.content}</LatestBlogPostItem>
            </BlogPostProvider>
          </div>
        </div>
      )}
      <div className="container margin-vert--lg">
        <div className={styles.itemGrid}>
          {(showLatest ? rest : items).map(({ content: BlogPostContent }) => (
            <BlogPostProvider
              key={BlogPostContent.metadata.permalink}
              content={BlogPostContent}
            >
              <BlogPostItemComponent>
                <BlogPostContent />
              </BlogPostItemComponent>
            </BlogPostProvider>
          ))}
        </div>
      </div>
    </>
  );
}