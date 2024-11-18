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

(() => {
  const modal = document.getElementById('modal-container');
  modal.addEventListener('show.bs.modal', (evt) => {
    // https://getbootstrap.com/docs/5.3/components/modal/#events
    // "hidden.bs.modal" is too early to clear innerHTML – the form submission from inside the modal would be cancelled
    modal.innerHTML = '';
  });
})();

function setSSEDisconnected() {
  const elem = document.getElementById('disconnected-toast');
  if (elem && !elem.classList.contains('show')) {
    document.getElementById('disconnected-toast').classList.add('show');
  }
}
document.addEventListener('htmx:sseError', function (evt) {
  console.log('htmx:sseError', evt);
  setSSEDisconnected();
});
document.addEventListener('htmx:sseClose', function (evt) {
  console.log('htmx:sseClose', evt);
  setSSEDisconnected();
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
    fetch('/giscus', {
      method: 'POST',
      body: formData,
    });
    window.giscusReported = true;
  }
}
window.addEventListener('message', handleGiscusMessage);
