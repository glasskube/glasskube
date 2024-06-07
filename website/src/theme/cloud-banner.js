import ExecutionEnvironment from '@docusaurus/ExecutionEnvironment';

export function onRouteDidUpdate({location, previousLocation}) {
  // Don't execute if we are still on the same page; the lifecycle may be fired
  // because the hash changes (e.g. when navigating between headings)
  if (location.pathname !== previousLocation?.pathname) {
    if (ExecutionEnvironment.canUseDOM) {
      try {
        const id = window.posthog.get_distinct_id();
        const elem = document.querySelector('#banner-cloud-link');
        if (elem) {
          elem.href = elem.href + id;
        }
      } catch (e) {}
    }
  }
}
