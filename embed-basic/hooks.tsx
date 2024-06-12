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

import {
  useEffect,
  useState,
  RefObject,
  ReactElement,
  isValidElement,
} from 'react';
import { createRoot } from 'react-dom/client';
import {
  GithubGistEmbed,
  CodePenEmbed,
  YouTubeEmbed,
  JSFiddleEmbed,
  FigmaEmbed,
  ExcalidrawEmbed,
  LoomEmbed,
  DropboxEmbed,
  TwitterEmbed,
} from './components';

interface Config {
  platform: string;
  enable: boolean;
}

const useRenderEmbed = (
  element: HTMLElement | RefObject<HTMLElement> | null,
) => {
  const [configs, setConfigs] = useState<Config[] | null>(null);

  const embeds = [
    {
      name: 'YouTube',
      regexs: [
        /https:\/\/youtu\.be\/([a-zA-Z0-9_-]{11})/,
        /https:\/\/www\.youtube\.com\/watch\?v=([a-zA-Z0-9_-]{11})/,
        /https:\/\/www\.youtube\.com\/embed\/([a-zA-Z0-9_-]{11})/,
      ],
      embed: (videoId: string) => {
        return <YouTubeEmbed videoId={videoId} />;
      },
    },
    {
      name: 'Twitter',
      regexs: [
        /https:\/\/twitter\.com\/[a-zA-Z0-9_]+\/status\/([a-zA-Z0-9_]+)/,
        /https:\/\/x\.com\/[a-zA-Z0-9_]+\/status\/([a-zA-Z0-9_]+)/,
      ],
      embed: (_, url, title = '') => {
        const blockquoteElement = document.createElement('blockquote');
        blockquoteElement.classList.add('twitter-tweet');

        const anchorElement = document.createElement('a');
        anchorElement.href = url.replace('x.com', 'twitter.com');

        anchorElement.textContent = title;
        blockquoteElement.appendChild(anchorElement);
        const scriptElement = document.createElement('script');
        scriptElement.src = 'https://platform.twitter.com/widgets.js';
        scriptElement.async = true;

        const styleElement = document.createElement('style');
        styleElement.innerHTML = `
          .twitter-tweet {
            display: block;
            margin: 0 auto;
          }
        `;

        return (
          <TwitterEmbed
            url={url.replace('x.com', 'twitter.com')}
            title={title}
          />
        );
      },
    },
    {
      name: 'CodePen',
      regexs: [
        /https:\/\/codepen\.io\/[a-zA-Z0-9_]+\/pen\/([a-zA-Z0-9_]+)/,
        /https:\/\/codepen\.io\/[a-zA-Z0-9_]+\/full\/([a-zA-Z0-9_]+)/,
      ],
      embed: (penId) => {
        return <CodePenEmbed penId={penId} />;
      },
    },
    {
      name: 'JSFiddle',
      regexs: [
        /https:\/\/jsfiddle\.net\/[a-zA-Z0-9_]+\/([a-zA-Z0-9_]+)/,
        /https:\/\/jsfiddle\.net\/[a-zA-Z0-9_]+\/([a-zA-Z0-9_]+)\/embed/,
      ],
      embed: (fiddleId: string) => {
        return <JSFiddleEmbed fiddleId={fiddleId} />;
      },
    },
    {
      name: 'GithubGist',
      regexs: [
        /https:\/\/gist\.github\.com\/[a-zA-Z0-9_]+\/([a-zA-Z0-9_]+)/,
        /https:\/\/gist\.github\.com\/[a-zA-Z0-9_]+\/([a-zA-Z0-9_]+)\.js/,
      ],
      embed: (_, url) => {
        const scriptUrl = url.indexOf('.js') > -1 ? url : `${url}.js`;
        return <GithubGistEmbed scriptUrl={scriptUrl} />;
      },
    },
    {
      name: 'Figma',
      regexs: [
        /https:\/\/www\.figma\.com\/design\/[a-zA-Z0-9_]+\/([a-zA-Z0-9_]+)/,
        /https:\/\/www\.figma\.com\/file\/[a-zA-Z0-9_]+\/([a-zA-Z0-9_]+)/,
      ],
      embed: (_, url) => {
        return <FigmaEmbed url={url} />;
      },
    },
    {
      name: 'Excalidraw',
      regexs: [
        /https:\/\/excalidraw\.com\/#json=([a-zA-Z0-9_,-]+)/,
        /https:\/\/excalidraw\.com\/([a-zA-Z0-9_,-]+)/,
      ],
      embed: (excalidrawId: string) => {
        return <ExcalidrawEmbed excalidrawId={excalidrawId} />;
      },
    },
    {
      name: 'Loom',
      regexs: [
        /https:\/\/www\.loom\.com\/embed\/([a-zA-Z0-9_]+)/,
        /https:\/\/www\.loom\.com\/share\/([a-zA-Z0-9_]+)/,
      ],
      embed: (loomId: string) => {
        return <LoomEmbed loomId={loomId} />;
      },
    },
    {
      name: 'Dropbox',
      regexs: [
        /https:\/\/www\.dropbox\.com\/s\/([a-zA-Z0-9_]+)\/[a-zA-Z0-9_]+/,
      ],
      embed: (dropboxId: string) => {
        return <DropboxEmbed dropboxId={dropboxId} />;
      },
    },
  ];

  const filteredEmbeds = embeds.filter((embed) => {
    const finded = configs?.find(
      (config) => config.platform === embed.name && config.enable,
    );
    return finded;
  });

  const renderEmbed = (
    url: string,
    title: string,
  ): string | HTMLElement | HTMLElement[] => {
    let html: string | HTMLElement | HTMLElement[] | ReactElement = '';

    filteredEmbeds.forEach((embed) => {
      if (html) return;
      embed.regexs.forEach((regex) => {
        if (html) return;
        const match = url.match(regex);
        if (match) {
          html = embed.embed(match[1], url, title);
        }
      });
    });

    return html;
  };

  const render = (targetElement) => {
    if (!element) {
      return;
    }

    const links = targetElement.querySelectorAll('a');
    let hasDefaultStyle = false;
    links.forEach((link) => {
      const url = link.getAttribute('href') || '';
      const title = link.getAttribute('title') || '';
      if (!url) {
        return;
      }
      if (title !== '@embed') {
        return;
      }
      const embed = renderEmbed(url, link?.textContent || '');
      if (isValidElement(embed)) {
        const parentElement = link.parentElement as HTMLElement;
        parentElement.classList.add('position-relative');
        parentElement.style.height = '128px';
        createRoot(parentElement).render(embed);
      } else {
        hasDefaultStyle = true;
        link.innerHTML = `
          <div class="card embed-light">
            <div class="card-body">
              <div class="text-secondary small mb-1">${url}</div>
              <div class="text-body fw-bold">${link.textContent}</div>
            </div>
          </div>
        `;
      }
    });
    // default card style add embed-ligh class for hover bg-light
    let styleElement = document.querySelector('style#embed-style');
    if (!styleElement) {
      styleElement = document.createElement('style');
      styleElement.id = 'embed-style';
      if (hasDefaultStyle) {
        styleElement.textContent =  `
         .embed-light:hover {
           --bs-bg-opacity: 1;
           background-color: rgba(var(--bs-light-rgb), var(--bs-bg-opacity)) !important;
         }
        `
        // style 插入 header
        const head = document.querySelector('head') as HTMLElement;
        head.appendChild(styleElement);
      }
    }
  };

  const getConfig = () => {
    fetch('/answer/api/v1/embed/config')
      .then((response) => response.json())
      .then((result) => setConfigs(result.data));
  };
  useEffect(() => {
    getConfig();

    return () => {
      const styleEle = document.querySelector('style#embed-style');
      const head = document.querySelector('head') as HTMLElement;
      if (styleEle) {
        head.removeChild(styleEle);
      }
    }
  }, []);

  useEffect(() => {
    if (!element) {
      return;
    }

    if (!configs) {
      return;
    }

    let targetElement;
    if (element instanceof HTMLElement) {
      targetElement = element;
    } else {
      targetElement = element.current;
    }
    render(targetElement);
    const observer = new MutationObserver(() => {
      render(targetElement);
    });

    observer.observe(targetElement, {
      childList: true,
    });

    return () => {
      observer.disconnect();
    };
  }, [element, configs]);
};

export { useRenderEmbed };
