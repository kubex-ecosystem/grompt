import { useEffect } from 'react';
import { Idea } from '../types';

interface UseUrlSharingProps {
    setIdeas: (ideas: Idea[]) => void;
    setPurpose: (purpose: string) => void;
    setGeneratedPrompt: (prompt: string) => void;
    setTokenUsage: (usage: null) => void;
    setError: (error: string | null) => void;
}

export const useUrlSharing = ({ setIdeas, setPurpose, setGeneratedPrompt, setTokenUsage, setError }: UseUrlSharingProps) => {
    useEffect(() => {
        const loadSharedPrompt = () => {
            try {
                const hash = window.location.hash;
                if (hash.startsWith('#prompt=')) {
                    const encodedData = hash.substring('#prompt='.length);
                    const decodedJson = atob(encodedData);
                    const data = JSON.parse(decodedJson) as { ideas: Idea[], purpose: string, prompt: string };

                    if (data.ideas && data.purpose && data.prompt) {
                        setIdeas(data.ideas);
                        setPurpose(data.purpose);
                        setGeneratedPrompt(data.prompt);
                        setTokenUsage(null); // Tokens are not shared in link
                        setError(null);
                        window.history.replaceState(null, document.title, window.location.pathname + window.location.search);
                    }
                }
            } catch (e) {
                console.error("Failed to parse shared prompt from URL", e);
                setError("The shared link appears to be invalid or corrupted.");
                window.history.replaceState(null, document.title, window.location.pathname + window.location.search);
            }
        };
        loadSharedPrompt();
    }, [setIdeas, setPurpose, setGeneratedPrompt, setTokenUsage, setError]);
};
