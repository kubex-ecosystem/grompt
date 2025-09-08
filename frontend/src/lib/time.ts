export const formatRelativeTime = (timestamp: number, locale: string): string => {
    try {
        const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' });
        const seconds = Math.floor((Date.now() - timestamp) / 1000);

        if (seconds < 60) return rtf.format(Math.floor(-seconds), 'second');
        if (seconds < 3600) return rtf.format(-Math.floor(seconds / 60), 'minute');
        if (seconds < 86400) return rtf.format(-Math.floor(seconds / 3600), 'hour');
        if (seconds < 2592000) return rtf.format(-Math.floor(seconds / 86400), 'day');
        if (seconds < 31536000) return rtf.format(-Math.floor(seconds / 2592000), 'month');
        return rtf.format(-Math.floor(seconds / 31536000), 'year');
    } catch (e) {
        console.error("Error formatting relative time", e);
        return new Date(timestamp).toLocaleDateString();
    }
};
