/**
 * Layout module - Shared navbar with role-based navigation
 * Theo chuc_nang.pdf: Admin, Manager, Staff, Resident
 */

import authManager from './auth-manager.js';
import notificationAPI from '../api/notifications.js';
import { ROLES } from '../utils/constants.js';

const ROLE_LABELS = {
  admin: 'Quản trị viên',
  manager: 'Ban quản lý',
  staff: 'Nhân viên',
  resident: 'Cư dân',
};

/** Menu items cho Admin/Manager/Staff */
const ADMIN_MENU = [
  { id: 'dashboard', label: 'Dashboard', href: '/pages/dashboard.html', roles: ['admin', 'manager', 'staff'] },
  { id: 'gis', label: 'GIS 3D', href: '/index.html', roles: ['admin', 'manager', 'staff'], external: true },
  { id: 'buildings', label: 'Tòa nhà', href: '/pages/buildings.html', roles: ['admin', 'manager', 'staff'] },
  { id: 'residents', label: 'Cư dân', href: '/pages/residents.html', roles: ['admin', 'manager', 'staff'] },
  { id: 'contracts', label: 'Hợp đồng', href: '/pages/contracts.html', roles: ['admin', 'manager', 'staff'] },
  { id: 'invoices', label: 'Hóa đơn', href: '/pages/invoices.html', roles: ['admin', 'manager', 'staff'] },
  { id: 'maintenance', label: 'Bảo trì', href: '/pages/maintenance.html', roles: ['admin', 'manager', 'staff'] },
  { id: 'reports', label: 'Báo cáo', href: '/pages/reports.html', roles: ['admin', 'manager'] },
  { id: 'users', label: 'Tài khoản', href: '/pages/users.html', roles: ['admin'] },
  { id: 'notifications', label: 'Thông báo', href: '/pages/notifications.html', roles: ['admin', 'manager', 'staff', 'resident'] },
];

/** Menu items cho Resident Portal */
const RESIDENT_MENU = [
  { id: 'resident-dashboard', label: 'Tổng quan', href: '/pages/resident-dashboard.html', roles: ['resident'] },
  { id: 'gis', label: 'Căn hộ 3D', href: '/index.html', roles: ['resident'], external: true },
  { id: 'contracts', label: 'Hợp đồng', href: '/pages/contracts.html', roles: ['resident'] },
  { id: 'invoices', label: 'Hóa đơn', href: '/pages/invoices.html', roles: ['resident'] },
  { id: 'maintenance', label: 'Bảo trì', href: '/pages/maintenance.html', roles: ['resident'] },
  { id: 'notifications', label: 'Thông báo', href: '/pages/notifications.html', roles: ['resident'] },
];

function getCurrentPageId() {
  const path = window.location.pathname;
  const file = path.split('/').pop().replace('.html', '');
  return file || 'dashboard';
}

function buildNavItems(role) {
  const menu = role === ROLES.RESIDENT ? RESIDENT_MENU : ADMIN_MENU;
  const currentPage = getCurrentPageId();

  return menu
    .filter((item) => item.roles.includes(role))
    .map((item) => {
      const active = item.id === currentPage ? ' active' : '';
      const target = item.external ? ' target="_blank"' : '';
      return `<li><a href="${item.href}" class="nav-link${active}"${target}>${item.label}</a></li>`;
    })
    .join('');
}

function buildNavbarHTML(role, user, unreadCount) {
  const homeHref = role === ROLES.RESIDENT
    ? '/pages/resident-dashboard.html'
    : '/pages/dashboard.html';

  const badge = unreadCount > 0
    ? `<span class="notif-badge">${unreadCount > 99 ? '99+' : unreadCount}</span>`
    : '';

  return `
    <nav class="navbar">
      <div class="navbar-left">
        <button class="menu-toggle" id="menuToggle" title="Mở/đóng menu">☰</button>
        <a class="logo" href="${homeHref}">🏢 Vinhomes</a>
      </div>
      <div class="navbar-right">
        <a href="/pages/notifications.html" class="notif-btn" title="Thông báo">
          🔔${badge}
        </a>
        <div class="user-info">
          <div class="name">${user?.name || user?.username || 'User'}</div>
          <div class="role">${ROLE_LABELS[role] || role}</div>
        </div>
        <a class="avatar" href="/pages/profile.html" title="Tài khoản">
          ${(user?.name || user?.username || 'U')[0].toUpperCase()}
        </a>
        <button class="btn-logout" id="layoutLogoutBtn">Đăng xuất</button>
      </div>
    </nav>
    <aside class="sidebar" id="sidebar">
      <div class="sidebar-header">
        <span class="sidebar-title">Menu</span>
        <button class="sidebar-close" id="sidebarClose">✕</button>
      </div>
      <ul class="sidebar-menu">${buildNavItems(role)}</ul>
    </aside>
    <div class="sidebar-overlay" id="sidebarOverlay"></div>
  `;
}

/**
 * Khởi tạo layout: render navbar + gắn sự kiện logout
 * @param {object} opts
 * @param {string[]} opts.requiredRoles - vai trò được phép truy cập trang
 */
async function initLayout(opts = {}) {
  const { requiredRoles = null } = opts;

  if (!authManager.guardRoute(requiredRoles)) return false;

  const role = authManager.getUserRole();
  const user = authManager.getCurrentUser();

  // Tải số thông báo chưa đọc (không block nếu lỗi)
  let unreadCount = 0;
  try {
    const res = await notificationAPI.getUnreadNotifications();
    unreadCount = (res.data || []).length;
  } catch (_) { /* backend chưa chạy */ }

  const placeholder = document.getElementById('navbar-placeholder');
  const html = buildNavbarHTML(role, user, unreadCount);

  if (placeholder) {
    placeholder.innerHTML = html;
  } else {
    document.body.insertAdjacentHTML('afterbegin', html);
  }

  const logoutBtn = document.getElementById('layoutLogoutBtn');
  if (logoutBtn) {
    logoutBtn.addEventListener('click', () => {
      if (confirm('Bạn có chắc muốn đăng xuất?')) authManager.logout();
    });
  }

  // Sidebar toggle functionality
  const menuToggle = document.getElementById('menuToggle');
  const sidebar = document.getElementById('sidebar');
  const sidebarClose = document.getElementById('sidebarClose');
  const sidebarOverlay = document.getElementById('sidebarOverlay');

  if (menuToggle && sidebar) {
    menuToggle.addEventListener('click', () => {
      sidebar.classList.add('show');
      sidebarOverlay?.classList.add('show');
    });
  }

  if (sidebarClose && sidebar) {
    sidebarClose.addEventListener('click', () => {
      sidebar.classList.remove('show');
      sidebarOverlay?.classList.remove('show');
    });
  }

  if (sidebarOverlay && sidebar) {
    sidebarOverlay.addEventListener('click', () => {
      sidebar.classList.remove('show');
      sidebarOverlay.classList.remove('show');
    });
  }

  // Close sidebar when clicking on menu links
  const sidebarLinks = sidebar?.querySelectorAll('a');
  sidebarLinks?.forEach(link => {
    link.addEventListener('click', () => {
      sidebar.classList.remove('show');
      sidebarOverlay?.classList.remove('show');
    });
  });

  return true;
}

export { initLayout, ROLE_LABELS, getCurrentPageId };
