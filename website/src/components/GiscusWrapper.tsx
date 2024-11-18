import React, {useEffect} from 'react';
import Giscus from '@giscus/react';
import {useColorMode} from '@docusaurus/theme-common';

export default function GiscusWrapper({
  category,
  categoryId,
}: {
  category: string;
  categoryId: string;
}) {
  const {colorMode} = useColorMode();

  const handleGiscusMessage = ev => {
    if (ev.origin !== 'https://giscus.app') return;
    if (!(typeof ev.data === 'object' && ev.data.giscus)) return;

    const giscusData = ev.data.giscus;

    // Make sure that the message is one that contains the discussion metadata.
    // This is necessary because other message types are also available.
    if (giscusData['discussion'] && giscusData['viewer']) {
      const username = giscusData['viewer']['login'];
      if (username.includes('giscus') && username.includes('bot')) {
        return;
      }
      const githubUrl = giscusData['viewer']['url'];
      if (window['posthog']) {
        window['posthog'].setPersonProperties({
          github_url: githubUrl,
        });
      }
    }
  };

  useEffect(() => {
    window.addEventListener('message', handleGiscusMessage);
    return () => {
      window.removeEventListener('message', handleGiscusMessage);
    };
  }, []);

  return (
    <Giscus
      repo="glasskube/glasskube"
      repoId="R_kgDOLDumDw"
      category={category}
      categoryId={categoryId}
      mapping="title"
      strict="0"
      reactionsEnabled="1"
      emitMetadata="1"
      inputPosition="top"
      theme={colorMode}
      lang="en"
      loading="lazy"
    />
  );
}

export function BlogDiscussion() {
  return <GiscusWrapper category="Blog" categoryId="DIC_kwDOLDumD84CfCte" />;
}
export function Discussion({
  category,
  categoryId,
}: {
  category: string;
  categoryId: string;
}) {
  return <GiscusWrapper category={category} categoryId={categoryId} />;
}
