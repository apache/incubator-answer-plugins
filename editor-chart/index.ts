import Chart from './Chart';
import i18nConfig from './i18n';
import { useRenderChart } from './hooks';

export default {
  info: {
    type: 'editor',
    slug_name: 'chart_editor',
  },
  component: Chart,
  i18nConfig,
  hooks: {
    useRender: [useRenderChart],
  },
};
