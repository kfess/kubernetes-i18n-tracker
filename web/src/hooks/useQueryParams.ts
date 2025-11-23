import { useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';

export function useQueryParams<T>(
  key: string,
  defaultValue: T,
  serialize: (value: T) => string = String,
  deserialize: (value: string) => T = (value) => value as T
) {
  const [searchParams, setSearchParams] = useSearchParams();

  const value = searchParams.get(key) ? deserialize(searchParams.get(key) as string) : defaultValue;

  const setValue = useCallback(
    (newValue: T) => {
      const next = new URLSearchParams(searchParams);
      next.set(key, serialize(newValue));
      setSearchParams(next);
    },
    [key, serialize, searchParams, setSearchParams]
  );

  return [value, setValue] as const;
}
