// Constants
const DONUT_CONFIG = {
  outerRadius: 130,
  innerRadius: 80,
  hoverScale: 1.15,
  gapAngleDegrees: 0.5,
  minSliceAngle: 0.001
};

const COLOR_CONFIG = {
  red: { r: 239, g: 68, b: 68 },
  yellow: { r: 245, g: 158, b: 11 },
  green: { r: 34, g: 197, b: 94 }
};

// Device Detection
const isTouchDevice = () => {
  return 'ontouchstart' in window ||
         navigator.maxTouchPoints > 0 ||
         navigator.msMaxTouchPoints > 0;
};

// State Management
const state = {
  currentPath: [],
  currentNode: fileTree
};

// Tooltip Module
const Tooltip = (() => {
  let tooltipElement = null;

  function create() {
    tooltipElement = document.createElement('div');
    tooltipElement.style.position = 'fixed';
    tooltipElement.style.padding = '8px 10px';
    tooltipElement.style.background = 'rgba(2,6,23,0.9)';
    tooltipElement.style.borderRadius = '10px';
    tooltipElement.style.color = 'rgba(255, 255, 255, 1)';
    tooltipElement.style.pointerEvents = 'none';
    tooltipElement.style.boxShadow = '0 4px 6px rgba(0,0,0,0.3)';
    tooltipElement.style.display = 'none';
    document.body.appendChild(tooltipElement);
  }

  function show(x, y, text) {
    if (!tooltipElement) {
      create();
    }
    tooltipElement.textContent = text;
    tooltipElement.style.left = `${x + 12}px`;
    tooltipElement.style.top = `${y + 12}px`;
    tooltipElement.style.display = 'block';
  }

  function hide() {
    if (tooltipElement) {
      tooltipElement.style.display = 'none';
    }
  }

  return { show, hide };
})();

// Color Utilities
const ColorUtils = (() => {
  function interpolate(start, end, ratio) {
    return Math.round(start + (end - start) * ratio);
  }

  function getCoverageColr(pct) {
    const clampedPct = Math.max(0, Math.min(100, pct));

    if (clampedPct <= 50) {
      const ratio = clampedPct / 50;
      const r = COLOR_CONFIG.red.r;
      const g = interpolate(COLOR_CONFIG.red.g, COLOR_CONFIG.yellow.g, ratio);
      const b = interpolate(COLOR_CONFIG.red.b, COLOR_CONFIG.yellow.b, ratio);
      return `rgb(${r}, ${g}, ${b})`;
    } else {
      const ratio = (clampedPct - 50) / 50;
      const r = interpolate(COLOR_CONFIG.yellow.r, COLOR_CONFIG.green.r, ratio);
      const g = interpolate(COLOR_CONFIG.yellow.g, COLOR_CONFIG.green.g, ratio);
      const b = interpolate(COLOR_CONFIG.yellow.b, COLOR_CONFIG.green.b, ratio);
      return `rgb(${r}, ${g}, ${b})`;
    }
  }

  return { getCoverageColr };
})();

// Navigation Module
const Navigation = (() => {
  function navigateInto(node) {
    state.currentPath.push(node.name);
    state.currentNode = node;
    FileTreeRenderer.render();
  }

  function navigateToPath(index) {
    if (index === -1) {
      state.currentPath = [];
      state.currentNode = fileTree;
    } else {
      state.currentPath = state.currentPath.slice(0, index + 1);
      state.currentNode = fileTree;

      for (let i = 0; i <= index; i++) {
        const childName = state.currentPath[i];
        state.currentNode = state.currentNode.children.find(c => c.name === childName);
        if (!state.currentNode) {
          state.currentPath = [];
          state.currentNode = fileTree;
          break;
        }
      }
    }
    FileTreeRenderer.render();
  }

  function navigateToFile(localPath) {
    const htmlPath = localPath.replace(/\.[^.]+$/, '.html');
    window.location.href = `tree/${htmlPath}`;
  }

  return { navigateInto, navigateToPath, navigateToFile };
})();

// DOM Helper Functions
const DOMHelpers = (() => {
  function addTooltip(element, label) {
    if (isTouchDevice()) {
      return element;
    }

    element.addEventListener('mouseenter', (e) => {
      Tooltip.show(e.clientX, e.clientY, label);
    });
    element.addEventListener('mouseleave', () => {
      Tooltip.hide();
    });
    element.addEventListener('mousemove', (e) => {
      Tooltip.show(e.clientX, e.clientY, label);
    });
    element.addEventListener('click', (_e) => {
      Tooltip.hide();
    });

    return element;
  }

  function createStatCell(value, label) {
    const cell = document.createElement('td');
    cell.className = 'file-table-stat';
    cell.textContent = (value || 0).toString();

    return addTooltip(cell, label);
  }

  function createBreadcrumbItem(name, index) {
    const item = document.createElement('span');
    item.className = 'breadcrumb-item';
    item.textContent = name;
    item.dataset.index = index;
    item.style.cursor = 'pointer';
    item.addEventListener('click', () => {
      Navigation.navigateToPath(parseInt(item.dataset.index, 10));
    });
    return item;
  }

  function createBreadcrumbSeparator() {
    const sep = document.createElement('span');
    sep.className = 'breadcrumb-sep';
    sep.textContent = ' / ';
    return sep;
  }

  return { addTooltip, createStatCell, createBreadcrumbItem, createBreadcrumbSeparator };
})();

// File Tree Renderer
const FileTreeRenderer = (() => {
  function render() {
    const browser = document.getElementById('file-browser');
    if (!browser || !fileTree) return;

    browser.innerHTML = '';
    browser.appendChild(renderBreadcrumb());
    browser.appendChild(renderFileList());
    DonutChart.render();
  }

  function renderBreadcrumb() {
    const breadcrumb = document.createElement('div');
    breadcrumb.className = 'breadcrumb';

    breadcrumb.appendChild(DOMHelpers.createBreadcrumbItem('root', -1));

    state.currentPath.forEach((name, index) => {
      breadcrumb.appendChild(DOMHelpers.createBreadcrumbSeparator());
      breadcrumb.appendChild(DOMHelpers.createBreadcrumbItem(name, index));
    });

    return breadcrumb;
  }

  function renderFileList() {
    const items = state.currentNode.children || [];

    if (items.length === 0) {
      const emptyDiv = document.createElement('div');
      emptyDiv.style.padding = '20px';
      emptyDiv.style.color = 'var(--text-muted)';
      emptyDiv.textContent = 'No files';
      return emptyDiv;
    }

    const table = document.createElement('table');
    table.className = 'file-table';

    const thead = document.createElement('thead');
    const headerRow = document.createElement('tr');

    const headers = ['', 'Lines', 'Covered', 'Partial', 'Missed', 'Coverage'];
    headers.forEach((text, _index) => {
      const th = document.createElement('th');
      th.textContent = text;
      th.className = 'file-table-stat-header';
      headerRow.appendChild(th);
    });
    thead.appendChild(headerRow);
    table.appendChild(thead);

    const tbody = document.createElement('tbody');
    items.forEach(child => {
      tbody.appendChild(renderTreeItem(child));
    });
    table.appendChild(tbody);

    return table;
  }

  function renderTreeItem(child) {
    const row = document.createElement('tr');
    row.className = 'file-table-row';

    const nameCell = document.createElement('td');
    nameCell.className = 'file-table-name';
    const nameSpan = document.createElement('span');
    nameSpan.className = `tree-name${child.isDir ? ' dir' : ''}`;
    nameSpan.textContent = child.name;
    nameCell.appendChild(nameSpan);
    row.appendChild(nameCell);

    row.appendChild(DOMHelpers.createStatCell(child.trackedLines, 'Lines'));
    row.appendChild(DOMHelpers.createStatCell(child.coveredLines, 'Covered'));
    row.appendChild(DOMHelpers.createStatCell(child.partialLines, 'Partial'));
    row.appendChild(DOMHelpers.createStatCell(child.missedLines, 'Missed'));

    const pct = child.coveragePct || 0;
    const color = ColorUtils.getCoverageColr(pct);
    const coverageCell = document.createElement('td');
    coverageCell.className = 'file-table-coverage';
    coverageCell.textContent = `${pct.toFixed(1)}%`;
    coverageCell.style.color = color;
    DOMHelpers.addTooltip(coverageCell, 'Coverage');
    row.appendChild(coverageCell);

    row.style.cursor = 'pointer';
    row.addEventListener('click', () => {
      if (child.isDir) {
        Navigation.navigateInto(child);
      } else if (child.file?.localPath) {
        Navigation.navigateToFile(child.file.localPath);
      }
    });

    return row;
  }

  return { render };
})();

// Donut Chart Renderer
const DonutChart = (() => {
  function render() {
    const svg = document.getElementById('donut');
    if (!svg || !state.currentNode) return;

    svg.innerHTML = '';

    const items = state.currentNode.children || [];

    if (items.length === 0) {
      renderEmptyState(svg);
      return;
    }

    renderArcs(svg, items);
    renderCenterText(svg);
  }

  function renderEmptyState(svg) {
    const centerText = createSVGText(0, 6, '0 items', '18', '700');
    svg.appendChild(centerText);
  }

  function renderCenterText(svg) {
    const dirName = state.currentPath.length > 0
      ? state.currentPath[state.currentPath.length - 1]
      : 'root';
    const centerText = createSVGText(0, 6, dirName, '18', '700');
    svg.appendChild(centerText);
  }

  function createSVGText(x, y, text, fontSize, fontWeight) {
    const textElement = document.createElementNS("http://www.w3.org/2000/svg", "text");
    textElement.setAttribute('x', x);
    textElement.setAttribute('y', y);
    textElement.setAttribute('fill', 'rgba(255,255,255,0.95)');
    textElement.setAttribute('font-size', fontSize);
    textElement.setAttribute('text-anchor', 'middle');
    textElement.setAttribute('font-weight', fontWeight);
    textElement.textContent = text;
    return textElement;
  }

  function renderArcs(svg, items) {
    const total = items.reduce((acc, item) => acc + (item.trackedLines || 0), 0);
    const allZero = total === 0;
    const gapAngle = (DONUT_CONFIG.gapAngleDegrees * Math.PI) / 180;

    let angle = -Math.PI / 2;

    items.forEach((item) => {
      const arc = createArc(item, total, allZero, angle, gapAngle, items.length);
      svg.appendChild(arc.path);
      angle = arc.nextAngle;
    });
  }

  function createArc(item, total, allZero, startAngle, gapAngle, itemCount) {
    const ratio = allZero ? (1 / itemCount) : ((item.trackedLines || 0) / total);
    let sliceAngle = Math.max(DONUT_CONFIG.minSliceAngle, ratio * 2 * Math.PI);

    if (itemCount === 1) {
      sliceAngle = 2 * Math.PI - DONUT_CONFIG.minSliceAngle;
    } else {
      sliceAngle = Math.max(DONUT_CONFIG.minSliceAngle, sliceAngle - gapAngle);
    }

    const endAngle = startAngle + sliceAngle;
    const nextAngle = endAngle + (itemCount === 1 ? 0 : gapAngle);

    const path = createArcPath(item, startAngle, endAngle, sliceAngle);

    return { path, nextAngle };
  }

  function createArcPath(item, startAngle, endAngle, sliceAngle) {
    const path = document.createElementNS("http://www.w3.org/2000/svg", "path");
    const pathData = getArcPathData(startAngle, endAngle, sliceAngle, DONUT_CONFIG.outerRadius);

    path.setAttribute('d', pathData);
    path.setAttribute('fill', ColorUtils.getCoverageColr(item.coveragePct || 0));
    path.setAttribute('stroke', 'rgba(255,255,255,0.02)');
    path.style.cursor = 'pointer';
    path.style.transition = 'd 0.2s ease-out';

    attachArcEventHandlers(path, item, startAngle, endAngle, sliceAngle);

    return path;
  }

  function attachArcEventHandlers(path, item, startAngle, endAngle, sliceAngle) {
    const hoverRadius = DONUT_CONFIG.outerRadius * DONUT_CONFIG.hoverScale;

    path.addEventListener('click', () => {
      if (item.isDir) {
        Navigation.navigateInto(item);
      } else if (item.file?.localPath) {
        Navigation.navigateToFile(item.file.localPath);
      }
    });

    if (isTouchDevice()) {
      return;
    }

    path.addEventListener('mouseover', (e) => {
      const hoverPathData = getArcPathData(startAngle, endAngle, sliceAngle, hoverRadius);
      path.setAttribute('d', hoverPathData);
      path.setAttribute('opacity', '0.9');

      const label = `${item.name} - ${(item.coveragePct || 0).toFixed(1)}%`;
      Tooltip.show(e.clientX, e.clientY, label);
    });

    path.addEventListener('mouseout', () => {
      const normalPathData = getArcPathData(startAngle, endAngle, sliceAngle, DONUT_CONFIG.outerRadius);
      path.setAttribute('d', normalPathData);
      path.removeAttribute('opacity');
      Tooltip.hide();
    });
  }

  function getArcPathData(startAngle, endAngle, sliceAngle, outerRadius) {
    const large = sliceAngle > Math.PI ? 1 : 0;
    const innerRadius = DONUT_CONFIG.innerRadius;

    const x1 = Math.cos(startAngle) * outerRadius;
    const y1 = Math.sin(startAngle) * outerRadius;
    const x2 = Math.cos(endAngle) * outerRadius;
    const y2 = Math.sin(endAngle) * outerRadius;
    const x3 = Math.cos(endAngle) * innerRadius;
    const y3 = Math.sin(endAngle) * innerRadius;
    const x4 = Math.cos(startAngle) * innerRadius;
    const y4 = Math.sin(startAngle) * innerRadius;

    return [
      `M ${x1} ${y1}`,
      `A ${outerRadius} ${outerRadius} 0 ${large} 1 ${x2} ${y2}`,
      `L ${x3} ${y3}`,
      `A ${innerRadius} ${innerRadius} 0 ${large} 0 ${x4} ${y4}`,
      'Z'
    ].join(' ');
  }

  return { render };
})();

// Application Initialization
function init() {
  FileTreeRenderer.render();
}

init();
