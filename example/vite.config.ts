import config from '@openfish/ui/vite.config'
import { mergeConfig } from 'vite'

export default mergeConfig(config, {
  root: '.',
})
