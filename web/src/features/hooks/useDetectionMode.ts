import { useLocalStorage } from '@mantine/hooks';
import { type DetectionMode } from '@/features/types';

export const useDetectionMode = () => {
  const [selectedDetectionMode, setDetectionMode] = useLocalStorage<DetectionMode>({
    key: 'selected-detection-mode',
    defaultValue: 'combined',
    getInitialValueInEffect: false,
  });

  return [selectedDetectionMode, setDetectionMode] as const;
};
