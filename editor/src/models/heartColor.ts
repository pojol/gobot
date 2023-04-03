
import { useState } from 'react';

export default () => {
    const [heartColor, setHeatColor] = useState("#BDCDD6");

    return { heartColor, setHeatColor };
};
