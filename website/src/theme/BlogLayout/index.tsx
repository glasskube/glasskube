import type { Props } from "@theme/BlogLayout";
import Layout from "@theme/Layout";
import { LatestBlogPostItem } from "../BlogPostItems/LatestBlogPostItem/LatestBlogPostItem";

export default function BlogLayout(props: Props): JSX.Element {
  const { toc, children, ...layoutProps } = props;

  return <Layout {...layoutProps}>{children}</Layout>;
}