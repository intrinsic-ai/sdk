// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/**
 * Styles for the world viewer's container, and for the elements that the world
 * viewer creates outside of Angular templates, such as the overlays a
 * ViewController adds over its canvas.
 *
 * The styles are registered once per document (or shadow root) as a
 * constructable stylesheet. All selectors are wrapped in `:where()`, so they
 * have zero specificity and any rule that targets these classes overrides
 * them. Apps can therefore restyle the overlays with plain CSS, e.g.:
 *
 * ```css
 * .intr-world-viewer-tooltip { background-color: navy; }
 * ```
 *
 * The colors are also defined as CSS variables on the container, so apps can
 * override them, e.g. to match their theme:
 *
 * ```css
 * .intr-world-viewer-container {
 *   --intr-world-viewer-tooltip-background-color: navy;
 * }
 * ```
 */

/** Class of the element that a world viewer renders into. */
export const CONTAINER_CLASS = 'intr-world-viewer-container';
/** Class of the canvas that a view renders into. */
export const CANVAS_CLASS = 'intr-world-viewer-canvas';
/** Class of the view widget (gimbal) canvas. */
export const VIEW_WIDGET_CLASS = 'intr-world-viewer-view-widget';
/** Class of the tooltip element. */
export const TOOLTIP_CLASS = 'intr-world-viewer-tooltip';
/** Class of the text inside the tooltip element. */
export const TOOLTIP_TEXT_CLASS = 'intr-world-viewer-tooltip-text';
/** Class of the frame rate display. */
export const FPS_DISPLAY_CLASS = 'intr-world-viewer-fps-display';

const STYLES = `
:where(.${CONTAINER_CLASS}) {
  --intr-world-viewer-view-widget-background-color: rgba(255, 255, 255, 0.8);
  --intr-world-viewer-tooltip-background-color: #000;
  --intr-world-viewer-tooltip-text-color: #fff;
  --intr-world-viewer-fps-display-color: green;

  /* Set this as the reference element for the overlays' absolute positioning. */
  position: relative;
  /* Prevent scrollbars from messing up layout when resizing the window. */
  overflow: hidden;
}

:where(.${CANVAS_CLASS}) {
  position: absolute;
  inset: 0;
  outline: 0;
}

:where(.${VIEW_WIDGET_CLASS}) {
  position: absolute;
  bottom: 0;
  left: 50%;
  margin-bottom: 15px;
  border-radius: 50%;
  background-color: var(--intr-world-viewer-view-widget-background-color);
}

:where(.${TOOLTIP_CLASS}) {
  box-shadow:
    0px 0.333333px 1.5px rgba(0, 0, 0, 0.039),
    0px 1.25px 5px rgba(0, 0, 0, 0.19);
  position: absolute;
  pointer-events: none;
  background-color: var(--intr-world-viewer-tooltip-background-color);
  border-radius: 4px;
  padding: 4px 8px;
  max-width: 290px;
  overflow-wrap: break-word;
}

:where(.${TOOLTIP_TEXT_CLASS}) {
  padding: 0 4px 0 2px;
  font-size: 1rem;
  font-weight: bold;
  color: var(--intr-world-viewer-tooltip-text-color);
}

:where(.${FPS_DISPLAY_CLASS}) {
  color: var(--intr-world-viewer-fps-display-color);
  position: absolute;
  top: 0;
  left: 0;
}
`;

// The stylesheet of each document, created on first use. A constructed
// stylesheet can only be adopted within the document whose window created it,
// so we key these by their owning document.
const styleSheets = new WeakMap<Document, CSSStyleSheet>();

/**
 * Registers the world viewer styles with the document or shadow root that
 * contains `element`. Note that this may be in an iframe. Does nothing if they
 * are already registered there, or if `element` is not in a document.
 */
export function registerWorldViewerStyles(element: Element) {
  const root = element.getRootNode();
  // Detached elements have no document or shadow root to register with.
  if (!isDocumentOrShadowRoot(root)) return;
  const styleSheet = getStyleSheet(element.ownerDocument);
  if (styleSheet && !root.adoptedStyleSheets.includes(styleSheet)) {
    root.adoptedStyleSheets.push(styleSheet);
  }
}

function getStyleSheet(doc: Document): CSSStyleSheet | undefined {
  let styleSheet = styleSheets.get(doc);
  if (!styleSheet) {
    // NOTE: Documents without a window (e.g. from a DOMParser) can't be styled.
    const view = doc.defaultView;
    if (!view) return undefined;
    styleSheet = new view.CSSStyleSheet();
    styleSheet.replaceSync(STYLES);
    styleSheets.set(doc, styleSheet);
  }
  return styleSheet;
}

// NOTE: We use this instead of `instanceof Document || instanceof ShadowRoot`
// because those constructors are specific to this window, so the check would
// fail for nodes in iframes.
function isDocumentOrShadowRoot(
  node: Node,
): node is Node & DocumentOrShadowRoot {
  return 'adoptedStyleSheets' in node;
}
