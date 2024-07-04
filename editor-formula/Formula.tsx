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

import { FC, useState } from 'react';
import { Button, Dropdown } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

import 'katex/dist/katex.min.css';

import { useRenderFormula } from './hooks';

export interface FormulaProps {
  editor;
  previewElement: HTMLElement;
}

const Formula: FC<FormulaProps> = ({ editor, previewElement }) => {
  useRenderFormula(previewElement);
  const { t } = useTranslation('plugin', {
    keyPrefix: 'formula_editor.frontend',
  });
  const [isLocked, setLockState] = useState(false);

  const handleMouseEnter = () => {
    if (isLocked) {
      return;
    }
    setLockState(true);
  };

  const handleMouseLeave = () => {
    setLockState(false);
  };
  const formulaList = [
    {
      type: 'line',
      label: t('options.inline'),
    },
    {
      type: 'block',
      label: t('options.block'),
    },
  ];

  const handleClick = (type: string, label: string) => {
    if (!editor) {
      return;
    }
    const { wrapText } = editor;
    if (type === 'line') {
      wrapText('\\\\( ', ' \\\\)', label);
    } else {
      const cursor = editor.getCursor();

      wrapText('\n$$\n', '\n$$\n', label);

      editor.setSelection(
        { line: cursor.line + 2, ch: 0 },
        { line: cursor.line + 2, ch: label.length },
      );
    }
    editor?.focus();
  };
  return (
    <div className="toolbar-item-wrap">
      <Dropdown>
        <Dropdown.Toggle
          type="button"
          as={Button}
          variant="link"
          title={t('title')}
          className="p-0 b-0 btn-no-border toolbar text-body">
          <i className="bi bi-123" />
        </Dropdown.Toggle>
        <Dropdown.Menu
          onMouseEnter={handleMouseEnter}
          onMouseLeave={handleMouseLeave}>
          {formulaList.map((formula) => {
            return (
              <Dropdown.Item
                key={formula.label}
                onClick={(e) => {
                  e.preventDefault();
                  handleClick(formula.type, formula.label);
                }}>
                {formula.label}
              </Dropdown.Item>
            );
          })}
        </Dropdown.Menu>
      </Dropdown>
    </div>
  );
};

export default Formula;
