<template>
  <div class="voice-selector">
    <el-select
      v-model="selectedId"
      placeholder="选择音色"
      @change="handleChange"
    >
      <el-option
        v-for="voice in voiceList"
        :key="voice.id"
        :label="voice.name"
        :value="voice.id"
      >
        <span>{{ voice.name }}</span>
        <span style="color: #999; font-size: 12px; margin-left: 8px">
          {{ voice.gender }} - {{ voice.tone }}
        </span>
      </el-option>
    </el-select>
    <el-button
      v-if="showPlayButton"
      :icon="VideoPlay"
      circle
      size="small"
      @click="handlePreview"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { VideoPlay } from '@element-plus/icons-vue'
import { useVoiceStore } from '@/store/voice'

const props = defineProps<{
  modelValue?: number
  showPlayButton?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: number): void
  (e: 'preview', voiceId: number): void
}>()

const voiceStore = useVoiceStore()
const voiceList = computed(() => voiceStore.voiceList)

const selectedId = computed({
  get: () => props.modelValue ?? 0,
  set: (val: number) => {
    if (val) {
      emit('update:modelValue', val)
    }
  },
})

function handleChange(value: number) {
  emit('update:modelValue', value)
}

function handlePreview() {
  if (props.modelValue) {
    emit('preview', props.modelValue)
  }
}
</script>

<style scoped>
.voice-selector {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
