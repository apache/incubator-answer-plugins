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
import { Modal, Form, Button, InputGroup, Spinner } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';
import ReCAPTCHA from "react-google-recaptcha";

import ReactDOM from 'react-dom/client';

import { languageKeys } from './common';
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

  const { t, i18n } = useTranslation('plugin', {
    keyPrefix: 'google_v2_captcha.frontend',
  });

  const refKey = useRef<CaptchaKey>(captchaKey);
  const refCallback = useRef<SubmitCallback>();
  const pending = useRef(false);
  const autoInitCaptchaData = /email/i.test(refKey.current);

  const [isLoading, setIsLoading] = useState(true);
  const [stateShow, setStateShow] = useState(false);
  const [googleKey, setGoogleKey] = useState('');
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

  const getCaptchaKey = () => {
    fetch(`/answer/api/v1/captcha/config`)
      .then((resp) => {
        return resp.json();
      })
      .then((data) => {
        setGoogleKey(data?.data?.config.key);
      })
  }

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

  const handleChange = (token) => {
    setImgCode({
      value: token || '',
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

  useEffect(() => {
    if (stateShow) {
      getCaptchaKey();
      setTimeout(() => {
        setIsLoading(false);
      }, 800)
    } else {
      setTimeout(() => {
        setIsLoading(true);
      }, 100)
    }
  }, [stateShow]);

  useLayoutEffect(() => {
    refImgCode.current = imgCode;
    refCaptcha.current = captcha;
  }, [captcha, imgCode]);

  useEffect(() => {
    refRoot.current?.render(
      <Modal
        style={{ width: '336px', margin: '0 auto', right: 0 }}
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
              <InputGroup>
                {isLoading && !captchaKey && <div style={{ height: '78px' }} />}
                <div className={isLoading ? 'w-100 text-center d-block' : 'w-100 text-center d-none'} style={{ position: 'absolute', top: 0, left: 0, height: '78px', lineHeight: '78px' }}>
                  <Spinner animation="border" variant="secondary" />
                </div>
                {googleKey && (
                  <ReCAPTCHA
                    sitekey={googleKey}
                    theme="light"
                    size="normal"
                    className={isLoading ? 'invisible' : 'visible'}
                    hl={languageKeys[i18n.language] || 'en'}
                    onChange={(token) =>handleChange(token)}
                    onErrored={() => resetImgCode()}
                    onExpired={() => resetImgCode()}
                  />
                )}
                <Form.Control
                  type="text"
                  autoComplete="off"
                  className="d-none"
                  isInvalid={imgCode?.isInvalid}
                />
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
