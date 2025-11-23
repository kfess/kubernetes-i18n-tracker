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
      setSearchParams((prev) => {
        const next = new URLSearchParams(prev);
        next.set(key, serialize(newValue));
        return next;
      });
    },
    [key, serialize, setSearchParams]
  );

  return [value, setValue] as const;
}

export function useUpdateQueryParams() {
  const [, setSearchParams] = useSearchParams();

  return useCallback(
    (params: Record<string, string>) => {
      setSearchParams((prev) => {
        const next = new URLSearchParams(prev);
        Object.entries(params).forEach(([key, value]) => {
          next.set(key, value);
        });
        return next;
      });
    },
    [setSearchParams]
  );
}
