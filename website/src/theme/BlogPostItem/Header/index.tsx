import React from "react";
import BlogPostItemHeaderTitle from "@theme/BlogPostItem/Header/Title";
import BlogPostItemHeaderInfo from "@theme/BlogPostItem/Header/Info";
import BlogPostItemHeaderAuthors from "@theme/BlogPostItem/Header/Authors";
import {useBlogPost} from '@docusaurus/plugin-content-blog/client';

export default function BlogPostItemHeader(): JSX.Element {
  return (
    <header className="card__header">
      <BlogPostItemHeaderTitle />
      <BlogPostItemHeaderAuthors />
      <BlogPostItemHeaderInfo />
    </header>
  );
}