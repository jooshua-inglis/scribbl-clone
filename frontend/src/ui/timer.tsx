import { useEffect, useState } from 'react';

export default function CountDown({ time }: { time: Date }) {
    const [now, setNow] = useState(new Date);
    useEffect(() => {
        const i = setInterval(() => {
            setNow(new Date());
        }, 1000);
        return () => {
            clearInterval(i);
        };
    }, []);
    return <>{Math.round((time.getTime() - now.getTime()) / 1000)}</>;
}
