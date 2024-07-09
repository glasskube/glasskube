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

(() => {
  const dismissed = sessionStorage.getItem('cloud-info-dismissed') === 'true';
  if (!dismissed) {
    document.body
      .querySelector('#cloud-info')
      .classList.remove('visually-hidden');
  }
  document.body
    .querySelector('#cloud-info-close')
    .addEventListener('click', () => {
      sessionStorage.setItem('cloud-info-dismissed', true);
    });
})();

window.advancedOptions = function (currentContext) {
  return localStorage.getItem('advancedOptions_' + currentContext) === 'true';
};

// TODO fix disconnected when graceful close!!
// hoping for this to be merged & released soon: https://github.com/bigskysoftware/htmx-extensions/pull/31
function setSSEDisconnected() {
  document
    .getElementById('sse-error-container')
    .classList.remove('visually-hidden');
  document.getElementById('sse-error-container-message').innerHTML =
    'You are disconnected from the server. Make sure to run <code>glasskube serve</code> and refresh this page!';
}
document.addEventListener('htmx:sseError', setSSEDisconnected);
document.addEventListener('htmx:sseMessage', function (evt) {
  console.log(evt.detail.type);
  if (evt && evt.detail && evt.detail.type === 'close') {
    console.log('closing');
    setSSEDisconnected();
  }
});

window.giscusReported = false;
function handleGiscusMessage(ev) {
  if (window.giscusReported) return;
  if (ev.origin !== 'https://giscus.app') return;
  if (!(typeof ev.data === 'object' && ev.data.giscus)) return;

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
