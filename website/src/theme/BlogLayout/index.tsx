import type {Props} from '@theme/BlogLayout';
import Layout from '@theme/Layout';

export default function BlogLayout(props: Props): JSX.Element {
  const {children, ...layoutProps} = props;

  return <Layout {...layoutProps}>{children}</Layout>;
}
