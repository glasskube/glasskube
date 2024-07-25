import React, {FC} from 'react';
import Link from '@docusaurus/Link';

interface DemoButtonProps {
  additionalClassNames: string;
}

const DemoButton: FC<DemoButtonProps> = ({additionalClassNames}) => (
  <Link
    className={`button button--accent ${additionalClassNames}`}
    to="https://cal.glasskube.eu/team/founder/demo">
    Book free demo
  </Link>
);

export default DemoButton;
