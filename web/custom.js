(() => {
  const getColorSchemeQuery = () =>
    window.matchMedia('(prefers-color-scheme: dark)');
  const getPreferredTheme = () =>
    getColorSchemeQuery().matches ? 'dark' : 'light';
  const setPreferredTheme = () =>
    document.body.setAttribute('data-bs-theme', getPreferredTheme());
  setPreferredTheme();
  getColorSchemeQuery().addEventListener('change', () => setPreferredTheme());
})();

var sseOnline = true;
document.addEventListener('htmx:sseError', function () {
  sseOnline = false;
  document
    .getElementById('sse-error-container')
    .classList.remove('visually-hidden');
  document.getElementById('sse-error-container-message').innerHTML =
    'You are disconnected from the server. Make sure to run <code>glasskube serve</code> and refresh this page!';
});
document.addEventListener('htmx:sseOpen', function () {
  if (!sseOnline) {
    sseOnline = true;
    const msg = document.getElementById('sse-error-container-message');
    msg.innerText =
      'You have been disconnected for a while. Please refresh this page to make sure you are up to date!';
  }
});

window.giscusReported = false;
function handleGiscusMessage(ev) {
  if (window.giscusReported) return;
  if (ev.origin !== 'https://giscus.app') return;
  if (!(typeof ev.data === 'object' && ev.data.giscus)) return;
  console.log('giscusReported', window.giscusReported);

  const giscusData = ev.data.giscus;
  if (giscusData['discussion'] && giscusData['viewer']) {
    const username = giscusData['viewer']['login'];
    if (username.includes('giscus') && username.includes('bot')) {
      return;
    }
    const githubUrl = giscusData['viewer']['url'];
    const formData = new FormData();
    formData.append('githubUrl', githubUrl);
    fetch('', {
      method: 'POST',
      body: formData,
    });
    window.giscusReported = true;
  }
}
window.addEventListener('message', handleGiscusMessage);
