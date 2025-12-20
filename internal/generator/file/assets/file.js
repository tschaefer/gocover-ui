// Constants
const HIGHLIGHT_CONFIG = {
  hashPrefix: '#L',
  scrollBehavior: 'smooth',
  scrollBlock: 'center'
};

// State Management
const state = {
  currentHighlightedLine: null
};

// DOM Utilities
const DOMUtils = (() => {
  function getLineElements(lineNum) {
    return {
      line: document.getElementById('line-' + lineNum),
      linenum: document.getElementById('linenum-' + lineNum)
    };
  }

  function toggleElementHighlight(element, shouldHighlight) {
    if (!element) return;

    if (shouldHighlight) {
      element.classList.add('highlighted');
    } else {
      element.classList.remove('highlighted');
    }
  }

  return { getLineElements, toggleElementHighlight };
})();

// URL Management
const URLManager = (() => {
  function setLineHash(lineNum) {
    window.location.hash = HIGHLIGHT_CONFIG.hashPrefix.substring(1) + lineNum;
  }

  function clearHash() {
    history.pushState('', document.title, window.location.pathname + window.location.search);
  }

  function getLineNumberFromHash() {
    const hash = window.location.hash;
    if (!hash || !hash.startsWith(HIGHLIGHT_CONFIG.hashPrefix)) {
      return null;
    }

    const lineNum = parseInt(hash.substring(HIGHLIGHT_CONFIG.hashPrefix.length));
    return (lineNum && !isNaN(lineNum)) ? lineNum : null;
  }

  return { setLineHash, clearHash, getLineNumberFromHash };
})();

// Line Highlighting
const LineHighlighter = (() => {
  function clearHighlight(lineNum) {
    if (!lineNum) return;

    const elements = DOMUtils.getLineElements(lineNum);
    DOMUtils.toggleElementHighlight(elements.line, false);
    DOMUtils.toggleElementHighlight(elements.linenum, false);
  }

  function applyHighlight(lineNum) {
    if (!lineNum) return false;

    const elements = DOMUtils.getLineElements(lineNum);
    if (!elements.line || !elements.linenum) {
      return false;
    }

    DOMUtils.toggleElementHighlight(elements.line, true);
    DOMUtils.toggleElementHighlight(elements.linenum, true);
    return true;
  }

  function toggleHighlight(lineNum) {
    const elements = DOMUtils.getLineElements(lineNum);

    if (!elements.line || !elements.linenum) {
      return;
    }

    if (state.currentHighlightedLine && state.currentHighlightedLine !== lineNum) {
      clearHighlight(state.currentHighlightedLine);
    }

    const isHighlighted = elements.line.classList.contains('highlighted');

    DOMUtils.toggleElementHighlight(elements.line, !isHighlighted);
    DOMUtils.toggleElementHighlight(elements.linenum, !isHighlighted);

    if (!isHighlighted) {
      state.currentHighlightedLine = lineNum;
      URLManager.setLineHash(lineNum);
    } else {
      state.currentHighlightedLine = null;
      URLManager.clearHash();
    }
  }

  function scrollToLine(lineNum) {
    const elements = DOMUtils.getLineElements(lineNum);
    if (elements.line) {
      elements.line.scrollIntoView({
        behavior: HIGHLIGHT_CONFIG.scrollBehavior,
        block: HIGHLIGHT_CONFIG.scrollBlock
      });
    }
  }

  return { toggleHighlight, applyHighlight, scrollToLine };
})();

// Application Initialization
function init() {
  const lineNum = URLManager.getLineNumberFromHash();
  if (lineNum) {
    const success = LineHighlighter.applyHighlight(lineNum);
    if (success) {
      state.currentHighlightedLine = lineNum;
      LineHighlighter.scrollToLine(lineNum);
    }
  }
}

// Make toggleHighlight available globally for onclick handlers in HTML
window.toggleHighlight = LineHighlighter.toggleHighlight;

// Initialize on page load
window.addEventListener('DOMContentLoaded', init);

