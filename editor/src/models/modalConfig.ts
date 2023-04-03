
import { useState } from 'react';

export default () => {
    const [open, setOpen] = useState(() => {
        if (localStorage.remoteAddr === "" || localStorage.remoteAddr === undefined) {
            return true
        } else {
            return false
        }
    });

    return { open, setOpen };
};
