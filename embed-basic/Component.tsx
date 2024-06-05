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

import { useState, useEffect } from 'react';
import { Button } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

import EmbedModal from './modal';
import { useRenderEmbed } from './hooks';

const Component = ({ editor, previewElement }) => {
  const [show, setShowState] = useState(false);
  const { t } = useTranslation('plugin', {
    keyPrefix: 'basic_embed.frontend',
  });
  useRenderEmbed(previewElement);

  useEffect(() => {
    if (!editor) return;
    editor.addKeyMap({
      'Ctrl-m': handleShow,
    });
  }, [editor]);

  const handleShow = () => {
    setShowState(true);
  };

  const handleConfirm = ({ title, url }) => {
    setShowState(false);
    const cursor = editor.getCursor();
    if (cursor.ch !== 0) {
      editor.replaceSelection('\n');
    }
    const embed = `\n[${title}](${url} "@embed")\n`;
    editor.replaceSelection(embed);
    editor.focus();
  };
  return (
    <div className="toolbar-item-wrap">
      <Button
        variant="link"
        className="p-0 b-0 btn-no-border toolbar text-body"
        onClick={handleShow}
        title={`${t('label')} (Ctrl+m)`}>
        <i className="bi bi-window" />
      </Button>
      <EmbedModal
        show={show}
        setShowState={setShowState}
        onConfirm={handleConfirm}
      />
    </div>
  );
};

export default Component;
