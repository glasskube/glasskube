import React, {FC} from 'react';
import Link from '@docusaurus/Link';

interface DemoButtonProps {
  additionalClassNames: string;
}

const DemoButton: FC<DemoButtonProps> = ({additionalClassNames}) => (
  <Link
    className={`glasskube-telemetry-demo button button--accent ${additionalClassNames}`}
    to="https://cal.glasskube.com/team/gk/demo">
    Get a demo
  </Link>
);

export default DemoButton;
