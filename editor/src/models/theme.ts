
import { useState } from 'react';
import ThemeType from '@/constant/constant';

export default () => {
  const [themeValue, setThemeValue] = useState(
    () => {
      if (localStorage.theme === undefined || localStorage.theme === ThemeType.Light) {
        return ThemeType.Light
      } else {
        return ThemeType.Dark
      }
    }
  );

  return { themeValue, setThemeValue };
};
