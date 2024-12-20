import React, {FC} from 'react';
import Link from '@docusaurus/Link';

interface DemoButtonProps {
  additionalClassNames: string;
}

const DemoButton: FC<DemoButtonProps> = ({additionalClassNames}) => (
  <Link
    className={`glasskube-telemetry-demo button button--accent ${additionalClassNames}`}
    to="https://cal.glasskube.com/team/gk/demo">
    Book free demo
  </Link>
);

export default DemoButton;
