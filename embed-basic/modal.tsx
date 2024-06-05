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

import { useState } from 'react';
import { Modal, Button, Form } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

const EmbedModal = ({ show, setShowState, onConfirm }) => {
  const [title, setTitle] = useState({
    value: '',
    isInvalid: false,
    errorMsg: '',
  });
  const [url, setUrl] = useState({
    value: '',
    isInvalid: false,
    errorMsg: '',
  });
  const { t } = useTranslation('plugin', {
    keyPrefix: 'basic_embed.frontend',
  });

  const handleHide = () => {
    setShowState(false);
  };

  const handleChangeTitle = (e) => {
    setTitle({
      value: e.target.value,
      isInvalid: false,
      errorMsg: '',
    });
  };

  const handleChangeUrl = (e) => {
    setUrl({
      value: e.target.value,
      isInvalid: false,
      errorMsg: '',
    });
  };

  const handleSubmit = () => {
    if (!title.value) {
      setTitle({
        ...title,
        isInvalid: true,
        errorMsg: t('required_title'),
      });
      return;
    }
    if (!url.value) {
      setUrl({
        ...url,
        isInvalid: true,
        errorMsg: t('required_url'),
      });
      return;
    }
    const urlRegex =
      /^[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)/;
    if (!urlRegex.test(url.value)) {
      setUrl({
        ...url,
        isInvalid: true,
        errorMsg: t('invalid_url'),
      });
      return;
    }
    onConfirm({
      title: title.value,
      url: url.value,
    });
    setShowState(false);

    setTitle({
      value: '',
      isInvalid: false,
      errorMsg: '',
    });

    setUrl({
      value: '',
      isInvalid: false,
      errorMsg: '',
    });
  };
  return (
    <Modal show={show} onHide={handleHide}>
      <Modal.Header closeButton>
        <Modal.Title>{t('header')}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form>
          <Form.Group className="mb-3" controlId="editor.plugin.embed.title">
            <Form.Label>{t('title')}</Form.Label>
            <Form.Control
              type="text"
              value={title.value}
              isInvalid={title.isInvalid}
              onChange={handleChangeTitle}
            />
            <Form.Control.Feedback type="invalid">
              {title.errorMsg}
            </Form.Control.Feedback>
          </Form.Group>
          <Form.Group className="mb-3" controlId="editor.plugin.embed.url">
            <Form.Label>{t('url')}</Form.Label>
            <Form.Control
              type="url"
              value={url.value}
              isInvalid={url.isInvalid}
              onChange={handleChangeUrl}
            />
            <Form.Control.Feedback type="invalid">
              {url.errorMsg}
            </Form.Control.Feedback>
          </Form.Group>
        </Form>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="link" onClick={handleHide}>
          {t('cancel')}
        </Button>
        <Button variant="primary" onClick={handleSubmit}>
          {t('add')}
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

export default EmbedModal;
