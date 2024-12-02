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

import { useEffect, useState, FC } from 'react';
import hljs from 'highlight.js';
import { themeStyles } from './themeStyles';
import { pluginHookProps, Request } from './types';

const get = async (url: string) => {
  const response = await fetch(url);
  const { data } = await response.json();
  return data;
};
const useHighlightCode: FC<pluginHookProps> = (
  props:
    | HTMLElement
    | null
    | {
      current: HTMLElement | null;
    },
  request: Request = {
    get,
  },
) => {
  const [selectTheme, setSelectTheme] = useState<string>('stackoverflow');

  // Fetch theme from API
  useEffect(() => {
    request
      .get('/answer/api/v1/render/config')
      .then((result) => {
        if (result.select_theme) {
          setSelectTheme(result.select_theme);
        }
      })
      .catch((error) => {
        console.error('Error fetching theme:', error);
      });
  }, []);

  useEffect(() => {
    let element;
    if (props instanceof HTMLElement) {
      element = props;
    } else if (props && props.current instanceof HTMLElement) {
      element = props.current;
    } else {
      return;
    }

    const applyThemeCSS = async (theme: string) => {
      const existingStyleElement = document.querySelector(
        'style[data-theme-style="highlight"]',
      );
      if (existingStyleElement) existingStyleElement.remove();

      const styleElement = document.createElement('style');
      styleElement.setAttribute('data-theme-style', 'highlight');
      document.head.appendChild(styleElement);

      const themeMode = theme === 'dark' ? 'dark' : 'light';
      const selectedTheme = themeStyles[selectTheme] || themeStyles.default;

      // Dynamically import the corresponding style
      const css = await selectedTheme[themeMode]();
      styleElement.innerHTML = css.default;

      // Apply syntax highlighting
      element.querySelectorAll('pre code').forEach((block) => {
        hljs.highlightElement(block as HTMLElement);
        (block as HTMLElement).style.backgroundColor = 'transparent';
        (block as HTMLElement).style.padding = '0';
      });
    };

    // Get and apply the initial theme
    const currentTheme =
      document.documentElement.getAttribute('data-bs-theme') || 'light';
    applyThemeCSS(currentTheme);

    // Observe DOM changes (e.g., code block content changes)
    const contentObserver = new MutationObserver(() => {
      const newTheme =
        document.documentElement.getAttribute('data-bs-theme') || 'light';
      console.log(
        'Detected code content change, reapplying syntax highlighting, current theme:',
        newTheme,
      );
      applyThemeCSS(newTheme);
    });

    contentObserver.observe(element, {
      childList: true, // Observe changes to child elements
      subtree: true, // Observe the entire subtree
    });

    // Observe theme changes
    const themeObserver = new MutationObserver(() => {
      const newTheme =
        document.documentElement.getAttribute('data-bs-theme') || 'light';
      console.log('Detected theme change:', newTheme);
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
  }, [props, selectTheme]);

  return null;
};

export { useHighlightCode };
