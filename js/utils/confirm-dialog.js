/**
 * Popup xác nhận dùng chung (thay cho window.confirm)
 */

let dialogReady = false;

function ensureDialog() {
  if (dialogReady) return;
  document.body.insertAdjacentHTML(
    'beforeend',
    `
    <div class="confirm-overlay" id="confirmOverlay" aria-hidden="true">
      <div class="confirm-dialog" role="dialog" aria-modal="true" aria-labelledby="confirmTitle">
        <div class="confirm-icon" id="confirmIcon">⚠️</div>
        <h3 class="confirm-title" id="confirmTitle">Xác nhận</h3>
        <p class="confirm-message" id="confirmMessage"></p>
        <div class="confirm-actions">
          <button type="button" class="confirm-btn-cancel" id="confirmCancelBtn">Hủy</button>
          <button type="button" class="confirm-btn-ok" id="confirmOkBtn">Xác nhận</button>
        </div>
      </div>
    </div>
    `
  );
  dialogReady = true;
}

/**
 * @param {object} opts
 * @param {string} opts.title
 * @param {string} opts.message
 * @param {string} [opts.confirmText]
 * @param {string} [opts.cancelText]
 * @param {'default'|'danger'} [opts.variant]
 * @returns {Promise<boolean>}
 */
function showConfirm(opts = {}) {
  const {
    title = 'Xác nhận',
    message = 'Bạn có chắc chắn muốn tiếp tục?',
    confirmText = 'Xác nhận',
    cancelText = 'Hủy',
    variant = 'default',
  } = opts;

  ensureDialog();

  const overlay = document.getElementById('confirmOverlay');
  const titleEl = document.getElementById('confirmTitle');
  const messageEl = document.getElementById('confirmMessage');
  const iconEl = document.getElementById('confirmIcon');
  const okBtn = document.getElementById('confirmOkBtn');
  const cancelBtn = document.getElementById('confirmCancelBtn');

  titleEl.textContent = title;
  messageEl.textContent = message;
  okBtn.textContent = confirmText;
  cancelBtn.textContent = cancelText;

  iconEl.textContent = variant === 'danger' ? '🚪' : '❓';
  okBtn.className = variant === 'danger' ? 'confirm-btn-ok danger' : 'confirm-btn-ok';

  return new Promise((resolve) => {
    let settled = false;

    const finish = (result) => {
      if (settled) return;
      settled = true;
      overlay.classList.remove('show');
      overlay.setAttribute('aria-hidden', 'true');
      document.removeEventListener('keydown', onKeydown);
      okBtn.removeEventListener('click', onConfirm);
      cancelBtn.removeEventListener('click', onCancel);
      overlay.removeEventListener('click', onOverlayClick);
      resolve(result);
    };

    const onConfirm = () => finish(true);
    const onCancel = () => finish(false);
    const onOverlayClick = (e) => {
      if (e.target === overlay) finish(false);
    };
    const onKeydown = (e) => {
      if (e.key === 'Escape') finish(false);
      if (e.key === 'Enter') finish(true);
    };

    okBtn.addEventListener('click', onConfirm);
    cancelBtn.addEventListener('click', onCancel);
    overlay.addEventListener('click', onOverlayClick);
    document.addEventListener('keydown', onKeydown);

    overlay.classList.add('show');
    overlay.setAttribute('aria-hidden', 'false');
    cancelBtn.focus();
  });
}

export { showConfirm };
