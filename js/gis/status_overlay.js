/**
 * GIS Apartment Status Overlay — UC 2.7.8–2.7.12
 * Tô màu căn hộ theo trạng thái từ backend API
 */
(function () {
  const API_BASE = 'http://localhost:8080/api/v1';

  const STATUS_COLORS = {
    empty: [76, 175, 80, 0.92],       // Xanh lá — Trống
    rented: [33, 150, 243, 0.92],     // Xanh dương — Đang thuê
    expiring: [255, 193, 7, 0.95],    // Vàng — Sắp hết HĐ
    maintenance: [244, 67, 54, 0.95], // Đỏ — Bảo trì
    reserved: [255, 152, 0, 0.9],
  };

  const STATUS_LABELS = {
    empty: 'Trống',
    rented: 'Đang thuê',
    expiring: 'Sắp hết HĐ',
    maintenance: 'Bảo trì',
    reserved: 'Đặt trước',
  };

  const TOWER_META = {
    W1: { cx: 35, cy: -28, w: 50, d: 20, baseZ: 25.2, floorH: 3.3, face: 'north' },
    W2: { cx: 10, cy: 28, w: 50, d: 20, baseZ: 25.2, floorH: 3.3, face: 'south' },
    W3: { cx: -69, cy: 0, w: 20, d: 50, baseZ: 3.5, floorH: 3.3, face: 'east' },
  };

  let allApartments = [];
  let activeFilters = { empty: true, rented: true, expiring: true, maintenance: true };
  let activeTowers = { W1: true, W2: true, W3: true };
  let floorRange = { min: 1, max: 40 };
  let overlayEnabled = true;

  function waitForEngine(cb, tries) {
    if (window.WPEngine && window.WPEngine.addBox && window.redrawAll) {
      cb();
      return;
    }
    if (tries <= 0) return;
    setTimeout(() => waitForEngine(cb, tries - 1), 300);
  }

  async function fetchApartmentStatus() {
    try {
      const res = await fetch(`${API_BASE}/gis/apartment-status`);
      if (!res.ok) throw new Error('HTTP ' + res.status);
      const json = await res.json();
      allApartments = json.data || [];
      updateStatusBadge(allApartments.length);
    } catch (err) {
      console.warn('[GIS Status] Không tải được dữ liệu backend:', err.message);
      allApartments = [];
      updateStatusBadge(0, true);
    }
  }

  function updateStatusBadge(count, error) {
    const el = document.getElementById('gisStatusCount');
    if (!el) return;
    if (error) {
      el.textContent = 'Backend offline';
      el.style.color = '#f44336';
    } else {
      el.textContent = count + ' căn hộ';
      el.style.color = '#4CAF50';
    }
  }

  function getUnitPosition(tower, floorNum, unitIdx, totalUnits) {
    const t = TOWER_META[tower];
    if (!t) return null;
    const z = t.baseZ + (floorNum - 1) * t.floorH + 0.3;
    const h = t.floorH * 0.75;
    const slot = (unitIdx - 1) / Math.max(totalUnits, 1);
    const offset = (slot - 0.5) * 0.85;

    if (t.face === 'north') {
      return { x: t.cx + offset * t.w, y: t.cy + t.d / 2 + 0.5, z, w: t.w / totalUnits * 0.85, d: 0.6, h };
    }
    if (t.face === 'south') {
      return { x: t.cx + offset * t.w, y: t.cy - t.d / 2 - 0.5, z, w: t.w / totalUnits * 0.85, d: 0.6, h };
    }
    // east (W3)
    return { x: t.cx + t.w / 2 + 0.5, y: t.cy + offset * t.d, z, w: 0.6, d: t.d / totalUnits * 0.85, h };
  }

  function drawStatusOverlay() {
    if (!window.WPEngine || !overlayEnabled) {
      if (window.WPEngine?.clearStatusOverlay) window.WPEngine.clearStatusOverlay();
      return;
    }

    window.WPEngine.clearStatusOverlay();

    const floorGroups = {};
    allApartments.forEach((apt) => {
      const key = apt.tower_id + '-' + apt.floor_number;
      if (!floorGroups[key]) floorGroups[key] = [];
      floorGroups[key].push(apt);
    });

    allApartments.forEach((apt) => {
      if (!activeTowers[apt.tower_id]) return;
      if (apt.floor_number < floorRange.min || apt.floor_number > floorRange.max) return;

      const status = apt.display_status || apt.current_status;
      if (!activeFilters[status]) return;

      const key = apt.tower_id + '-' + apt.floor_number;
      const totalUnits = floorGroups[key]?.length || 4;
      const pos = getUnitPosition(apt.tower_id, apt.floor_number, apt.unit_index || 1, totalUnits);
      if (!pos) return;

      const color = STATUS_COLORS[status] || STATUS_COLORS.empty;

      window.WPEngine.addStatusBox({
        x: pos.x, y: pos.y, z: pos.z,
        width: pos.w, depth: pos.d, height: pos.h,
        color: color,
        metallic: 0.1, roughness: 0.4,
        attributes: {
          name: apt.apartment_code,
          id: apt.apartment_code,
          status: status,
          statusLabel: STATUS_LABELS[status] || status,
          area: apt.area,
          resident: apt.resident_name || '',
          contractEnd: apt.contract_end || '',
          type: 'apartment_status',
        },
        popup: {
          title: '{name}',
          content: [{
            type: 'fields',
            fieldInfos: [
              { fieldName: 'statusLabel', label: 'Trạng thái' },
              { fieldName: 'area', label: 'Diện tích (m²)' },
              { fieldName: 'resident', label: 'Cư dân' },
              { fieldName: 'contractEnd', label: 'Hết HĐ' },
            ],
          }],
        },
        type: 'apartment_status',
      });
    });
  }

  function setupFilterUI() {
    const filterStatuses = document.querySelectorAll('[data-status-filter]');
    filterStatuses.forEach((chk) => {
      chk.addEventListener('change', () => {
        activeFilters[chk.dataset.statusFilter] = chk.checked;
        redrawWithStatus();
      });
    });

    ['W1', 'W2', 'W3'].forEach((tid) => {
      const chk = document.getElementById('filterTower' + tid);
      if (chk) {
        chk.addEventListener('change', () => {
          activeTowers[tid] = chk.checked;
          redrawWithStatus();
        });
      }
    });

    const floorMin = document.getElementById('filterFloorMin');
    const floorMax = document.getElementById('filterFloorMax');
    const floorMinVal = document.getElementById('filterFloorMinVal');
    const floorMaxVal = document.getElementById('filterFloorMaxVal');

    function updateFloorRange() {
      if (floorMin) {
        floorRange.min = parseInt(floorMin.value, 10);
        if (floorMinVal) floorMinVal.textContent = floorRange.min;
      }
      if (floorMax) {
        floorRange.max = parseInt(floorMax.value, 10);
        if (floorMaxVal) floorMaxVal.textContent = floorRange.max;
      }
      redrawWithStatus();
    }

    if (floorMin) floorMin.addEventListener('input', updateFloorRange);
    if (floorMax) floorMax.addEventListener('input', updateFloorRange);

    const chkOverlay = document.getElementById('chkStatusOverlay');
    if (chkOverlay) {
      chkOverlay.addEventListener('change', () => {
        overlayEnabled = chkOverlay.checked;
        redrawWithStatus();
      });
    }

    const btnRefresh = document.getElementById('btnRefreshStatus');
    if (btnRefresh) {
      btnRefresh.addEventListener('click', async () => {
        btnRefresh.textContent = '⟳ Đang tải...';
        await fetchApartmentStatus();
        redrawWithStatus();
        btnRefresh.textContent = '⟳ Làm mới dữ liệu';
      });
    }
  }

  function redrawWithStatus() {
    if (typeof window.redrawAll === 'function') {
      window.redrawAll();
    }
    drawStatusOverlay();
  }

  // Hook redrawAll để vẽ lại overlay sau mỗi lần redraw mô hình
  function hookRedraw() {
    const orig = window.redrawAll;
    window.redrawAll = function () {
      orig();
      if (overlayEnabled) drawStatusOverlay();
    };
  }

  function enhanceClickReadout() {
    document.addEventListener('click', () => {
      // handled in app.js hitTest — extend via mutation observer on selection box
    });
  }

  async function init() {
    setupFilterUI();
    await fetchApartmentStatus();
    waitForEngine(() => {
      hookRedraw();
      redrawWithStatus();
      // Auto refresh mỗi 60s
      setInterval(async () => {
        await fetchApartmentStatus();
        if (overlayEnabled) drawStatusOverlay();
      }, 60000);
    }, 30);
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }

  window.WPStatusOverlay = {
    refresh: fetchApartmentStatus,
    redraw: redrawWithStatus,
    getApartments: () => allApartments,
  };
})();
