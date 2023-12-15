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

// @ts-ignore
import katexRender from 'katex/contrib/auto-render/auto-render';

const useRenderFormula = (element: HTMLElement) => {
  const render = (element) => {
    katexRender(element, {
      delimiters: [
        { left: '$$', right: '$$', display: true },
        { left: '$$<br>', right: '<br>$$', display: true },
        {
          left: '\\begin{equation}',
          right: '\\end{equation}',
          display: true,
        },
        { left: '\\begin{align}', right: '\\end{align}', display: true },
        { left: '\\begin{alignat}', right: '\\end{alignat}', display: true },
        { left: '\\begin{gather}', right: '\\end{gather}', display: true },
        { left: '\\(', right: '\\)', display: false },
        { left: '\\[', right: '\\]', display: true },
      ],
    });
  };
  useEffect(() => {
    if (!element) {
      return;
    }

    render(element);
    const observer = new MutationObserver(() => {
      render(element);
    });

    observer.observe(element, {
      childList: true,
      attributes: true,
      subtree: true,
    });
  }, [element]);
};

export { useRenderFormula };
