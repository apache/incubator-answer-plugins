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

import { useEffect, useRef, useState, useLayoutEffect } from 'react';
import { Modal, Form, Button, InputGroup } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

import ReactDOM from 'react-dom/client';

import type {
  FormValue,
  ImgCodeRes,
  CaptchaKey,
  FieldError,
  ImgCodeReq,
} from './interface';
export interface Props {
  captchaKey: CaptchaKey;
  commonProps?: any;
}

const checkImgCode = ({
  captchaKey,
  commonProps,
}: Props) => {
  return new Promise<ImgCodeRes>((resolve) => {
    fetch(`/answer/api/v1/user/action/record?action=${captchaKey}`, {
      headers: {
        ...commonProps?.headers
      }
    })
      .then((resp) => {
        return resp.json();
      })
      .then((data) => {
        console.log('checkImgCode', data);
        resolve(data.data);
      });
  });
};

type SubmitCallback = {
  (): void;
};

const Index = ({
  captchaKey,
  commonProps,
}: Props) => {
  const refRoot = useRef<ReactDOM.Root | null>(null);
  if (refRoot.current === null) {
    refRoot.current = ReactDOM.createRoot(document.createElement('div'));
  }

  const { t } = useTranslation('plugin', {
    keyPrefix: 'basic_captcha.frontend',
  });
  const refKey = useRef<CaptchaKey>(captchaKey);
  const refCallback = useRef<SubmitCallback>();
  const pending = useRef(false);
  const autoInitCaptchaData = /email/i.test(refKey.current);

  const [stateShow, setStateShow] = useState(false);
  const [captcha, setCaptcha] = useState<ImgCodeRes>({
    captcha_id: '',
    captcha_img: '',
    verify: false,
  });
  const [imgCode, setImgCode] = useState<FormValue>({
    value: '',
    isInvalid: false,
    errorMsg: '',
  });
  const refCaptcha = useRef(captcha);
  const refImgCode = useRef(imgCode);

  const fetchCaptchaData = () => {
    pending.current = true;
    checkImgCode({
      captchaKey: refKey.current,
      commonProps,
    })
      .then((resp) => {
        setCaptcha(resp);
      })
      .finally(() => {
        pending.current = false;
      });
  };

  const resetCapture = () => {
    setCaptcha({
      captcha_id: '',
      captcha_img: '',
      verify: false,
    });
  };

  const resetImgCode = () => {
    setImgCode({
      value: '',
      isInvalid: false,
      errorMsg: '',
    });
  };
  const resetCallback = () => {
    refCallback.current = undefined;
  };

  const show = () => {
    if (!stateShow) {
      setStateShow(true);
    }
  };
  /**
   * There are some cases where the React scheduler cancels the execution of some functions,
   * which prevents them from closing properly:
   *  for example, if the parent component uninstalls the child component directly,
   *  and the `captchaModal.close()` call is inside the child component.
   * In this case, call `await captchaModal.close()` and wait for the close action to complete.
   */
  const close = () => {
    setStateShow(false);
    resetCapture();
    resetImgCode();
    resetCallback();

    const p = new Promise<void>((resolve) => {
      setTimeout(resolve, 50);
    });
    return p;
  };

  const handleCaptchaError = (fel: FieldError[] = []) => {
    console.log('handleCaptchaError', fel);
    const captchaErr = fel.find((o) => {
      return o.error_field === 'captcha_code';
    });

    const ri = refImgCode.current;
    if (captchaErr) {
      /**
       * `imgCode.value` No value but a validation error is received,
       * indicating that it is the first time the interface has returned a CAPTCHA error,
       * triggering the CAPTCHA logic. There is no need to display the error message at this point.
       */
      if (ri.value) {
        setImgCode({
          ...ri,
          isInvalid: true,
          errorMsg: captchaErr.error_msg,
        });
      }
      fetchCaptchaData();
      show();
    } else {
      close();
    }
    // Assist business logic in filtering CAPTCHA error messages when necessary
    return captchaErr;
  };

  const handleChange = (evt) => {
    evt.preventDefault();
    setImgCode({
      value: evt.target.value || '',
      isInvalid: false,
      errorMsg: '',
    });
  };

  const getCaptcha = () => {
    const rc = refCaptcha.current;
    const ri = refImgCode.current;
    const r = {
      verify: !!rc?.verify,
      captcha_id: rc?.captcha_id,
      captcha_code: ri.value,
    };

    return r;
  };

  const resolveCaptchaReq = (req: ImgCodeReq) => {
    const r = getCaptcha();
    if (r.verify) {
      req.captcha_code = r.captcha_code;
      req.captcha_id = r.captcha_id;
    }
  };

  const handleSubmit = (evt) => {
    evt.preventDefault();
    if (!imgCode.value) {
      return;
    }

    if (refCallback.current) {
      refCallback.current();
    }
  };

  useEffect(() => {
    if (autoInitCaptchaData) {
      fetchCaptchaData();
    }
  }, []);

  useLayoutEffect(() => {
    refImgCode.current = imgCode;
    refCaptcha.current = captcha;
  }, [captcha, imgCode]);

  useEffect(() => {
    refRoot.current?.render(
      <Modal
        size="sm"
        title="Captcha"
        show={stateShow}
        onHide={() => close()}
        centered>
        <Modal.Header closeButton>
          <Modal.Title as="h5">{t('title')}</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form noValidate onSubmit={handleSubmit}>
            <Form.Group controlId="code" className="mb-3">
              <div className="mb-3 p-2 d-flex align-items-center justify-content-center bg-light rounded-2">
                <img
                  src={captcha?.captcha_img}
                  alt="captcha img"
                  width="auto"
                  height="60px"
                />
              </div>
              <InputGroup>
                <Form.Control
                  type="text"
                  autoComplete="off"
                  placeholder={t('placeholder')}
                  isInvalid={imgCode?.isInvalid}
                  onChange={handleChange}
                  value={imgCode.value}
                />
                <Button
                  onClick={fetchCaptchaData}
                  variant="outline-secondary"
                  title={t('refresh', { keyPrefix: 'btns' })}
                  style={{
                    borderTopRightRadius: '0.375rem',
                    borderBottomRightRadius: '0.375rem',
                  }}>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="16"
                    height="16"
                    fill="currentColor"
                    className="bi bi-arrow-repeat"
                    viewBox="0 0 16 16">
                    <path d="M11.534 7h3.932a.25.25 0 0 1 .192.41l-1.966 2.36a.25.25 0 0 1-.384 0l-1.966-2.36a.25.25 0 0 1 .192-.41m-11 2h3.932a.25.25 0 0 0 .192-.41L2.692 6.23a.25.25 0 0 0-.384 0L.342 8.59A.25.25 0 0 0 .534 9" />
                    <path
                      fillRule="evenodd"
                      d="M8 3c-1.552 0-2.94.707-3.857 1.818a.5.5 0 1 1-.771-.636A6.002 6.002 0 0 1 13.917 7H12.9A5 5 0 0 0 8 3M3.1 9a5.002 5.002 0 0 0 8.757 2.182.5.5 0 1 1 .771.636A6.002 6.002 0 0 1 2.083 9z"
                    />
                  </svg>
                </Button>

                <Form.Control.Feedback type="invalid">
                  {imgCode?.errorMsg}
                </Form.Control.Feedback>
              </InputGroup>
            </Form.Group>

            <div className="d-grid">
              <Button type="submit" disabled={!imgCode.value}>
                {t('verify')}
              </Button>
            </div>
          </Form>
        </Modal.Body>
      </Modal>,
    );
  });

  const r = {
    close,
    show,
    check: (submitFunc: SubmitCallback) => {
      if (pending.current) {
        return false;
      }
      refCallback.current = submitFunc;
      if (captcha?.verify) {
        show();
        return false;
      }
      return submitFunc();
    },
    getCaptcha,
    resolveCaptchaReq,
    fetchCaptchaData,
    handleCaptchaError,
  };

  return r;
};

export default Index;
