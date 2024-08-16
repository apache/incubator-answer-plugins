/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import { useEffect } from 'react';
import hljs from 'highlight.js';

import githubLightCss from 'highlight.js/styles/github.css?inline';
import githubDarkCss from 'highlight.js/styles/github-dark.css?inline';

const useHighlightCode = (element: HTMLElement | null) => {
  useEffect(() => {
    if (!element) {
      return;
    }

    const applyThemeCSS = (theme: string) => {
      const existingStyleElement = document.querySelector('style[data-theme-style="highlight"]');
      if (existingStyleElement) {
        existingStyleElement.remove();
      }

      const styleElement = document.createElement('style');
      styleElement.setAttribute('data-theme-style', 'highlight');
      document.head.appendChild(styleElement);

      if (theme === 'dark') {
        styleElement.innerHTML = githubDarkCss;
      } else {
        styleElement.innerHTML = githubLightCss;
      }

      // Highlight code blocks
      element.querySelectorAll('pre code').forEach((block) => {
        hljs.highlightElement(block as HTMLElement);
      });
    };

    // Get and apply the initial theme
    const currentTheme = document.documentElement.getAttribute('data-bs-theme') || 'light';
    applyThemeCSS(currentTheme);

    // Observe DOM changes（e.g. the content of a code block changes）
    const contentObserver = new MutationObserver(() => {
      const newTheme = document.documentElement.getAttribute('data-bs-theme') || 'light';
      applyThemeCSS(newTheme);
    });

    contentObserver.observe(element, {
      childList: true, // Observe element changes
      subtree: true,   // Observe whole subtree
    });

    // Observe theme changes
    const themeObserver = new MutationObserver(() => {
      const newTheme = document.documentElement.getAttribute('data-bs-theme') || 'light';
      applyThemeCSS(newTheme);
    });

    themeObserver.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['data-bs-theme'],
    });

    return () => {
      contentObserver.disconnect();
      themeObserver.disconnect();
    };
  }, [element]);
};

export { useHighlightCode };