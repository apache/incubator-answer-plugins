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

import { useEffect, useState } from 'react';

interface Config {
  platform: string;
  enable: boolean;
}
const useRenderEmbed = (element: HTMLElement) => {
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
        return `<iframe width="100%" height="380" src="https://www.youtube.com/embed/${videoId}" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`;
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

        return [styleElement, blockquoteElement, scriptElement];
      },
    },
    {
      name: 'CodePen',
      regexs: [
        /https:\/\/codepen\.io\/[a-zA-Z0-9_]+\/pen\/([a-zA-Z0-9_]+)/,
        /https:\/\/codepen\.io\/[a-zA-Z0-9_]+\/full\/([a-zA-Z0-9_]+)/,
      ],
      embed: (penId) => {
        return `<iframe width="100%" height="380" scrolling="no" title="CodePen Embed" src="https://codepen.io/${penId}/embed/preview/${penId}?height=265&theme-id=0&default-tab=result" frameborder="no" allowtransparency="true" allowfullscreen="true"></iframe>`;
      },
    },
    {
      name: 'JSFiddle',
      regexs: [
        /https:\/\/jsfiddle\.net\/[a-zA-Z0-9_]+\/([a-zA-Z0-9_]+)/,
        /https:\/\/jsfiddle\.net\/[a-zA-Z0-9_]+\/([a-zA-Z0-9_]+)\/embed/,
      ],
      embed: (fiddleId: string) => {
        return `<iframe width="100%" height="380" src="https://jsfiddle.net/${fiddleId}/embedded/" allowfullscreen="allowfullscreen" allowpaymentrequest frameborder="0"></iframe>`;
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
        console.log(scriptUrl);
        return `<iframe 
        width="100%"
        height="350"    
        src="data:text/html;charset=utf-8,
        <head><base target='_blank' /></head>
        <body style='margin:0;'><script src='${scriptUrl}'></script>
        </body>">`;
      },
    },
    {
      name: 'Figma',
      regexs: [
        /https:\/\/www\.figma\.com\/design\/[a-zA-Z0-9_]+\/([a-zA-Z0-9_]+)/,
        /https:\/\/www\.figma\.com\/file\/[a-zA-Z0-9_]+\/([a-zA-Z0-9_]+)/,
      ],
      embed: (_, url) => {
        return `<iframe style="border: none;" width="100%" height="450
" src="https://www.figma.com/embed?embed_host=share&url=${url}" allowfullscreen></iframe>`;
      },
    },
    {
      name: 'Excalidraw',
      regexs: [
        /https:\/\/excalidraw\.com\/#json=([a-zA-Z0-9_,-]+)/,
        /https:\/\/excalidraw\.com\/([a-zA-Z0-9_,-]+)/,
      ],
      embed: (excalidrawId: string) => {
        return `<iframe width="100%" height="380" src="https://excalidraw.com/${excalidrawId}/embed" frameborder="0"></iframe>`;
      },
    },
    {
      name: 'Loom',
      regexs: [
        /https:\/\/www\.loom\.com\/embed\/([a-zA-Z0-9_]+)/,
        /https:\/\/www\.loom\.com\/share\/([a-zA-Z0-9_]+)/,
      ],
      embed: (loomId: string) => {
        return `<iframe width="100%" height="380" src="https://www.loom.com/embed/${loomId}" frameborder="0"></iframe>`;
      },
    },
    {
      name: 'Dropbox',
      regexs: [
        /https:\/\/www\.dropbox\.com\/s\/([a-zA-Z0-9_]+)\/[a-zA-Z0-9_]+/,
      ],
      embed: (dropboxId: string) => {
        return `<iframe width="100%" height="380" src="https://www.dropbox.com/s/${dropboxId}?raw=1" frameborder="0"></iframe>`;
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
    let html: string | HTMLElement | HTMLElement[] = '';

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

  const render = () => {
    if (!element) {
      return;
    }

    const links = element.querySelectorAll('a');
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
      if (embed) {
        if (typeof embed === 'string') {
          link.innerHTML = embed;
        } else if (Array.isArray(embed)) {
          link.innerHTML = '';
          embed.forEach((item) => {
            link.appendChild(item);
          });
        } else {
          link.innerHTML = '';
          link.appendChild(embed);
        }
      } else {
        link.innerHTML = `
          <div class="border rounded p-3">
            <div class="text-secondary">${url}</div>
            <div class="text-body">${link.textContent}</div>
          </div>
        `;
      }
    });
  };

  const getConfig = () => {
    fetch('/answer/api/v1/embed/config')
      .then((response) => response.json())
      .then((result) => setConfigs(result.data));
  };
  useEffect(() => {
    getConfig();
  }, []);

  useEffect(() => {
    if (!element) {
      return;
    }

    if (!configs) {
      return;
    }

    render();
    const observer = new MutationObserver(() => {
      render();
    });

    observer.observe(element, {
      childList: true,
    });
  }, [element, configs]);
};

export { useRenderEmbed };
