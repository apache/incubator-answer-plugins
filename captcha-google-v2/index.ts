import Captcha from './Captcha';
import i18nConfig from './i18n';
import useCaptcha from './useCaptcha';

export default {
  info: {
    type: 'captcha',
    slug_name: 'captcha_google_v2',
  },
  component: Captcha,
  i18nConfig,
  hooks: {
    useCaptcha: useCaptcha,
  },
};
