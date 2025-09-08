import { useState, useEffect, useCallback } from 'react';
import { HistoryItem } from '../types';
import { getHistory, saveHistory } from '../services/storageService';

export const usePromptHistory = () => {
    const [promptHistory, setPromptHistory] = useState<HistoryItem[]>([]);

    // Load history from the storage service on initial render
    useEffect(() => {
        const loadHistory = async () => {
            const history = await getHistory();
            setPromptHistory(history);
        };
        loadHistory();
    }, []);

    // Save history to the storage service whenever it changes
    useEffect(() => {
        // Avoid saving on initial empty state before anything is loaded
        if (promptHistory.length === 0) {
            const checkStorage = async () => {
                const existing = await getHistory();
                if (existing.length === 0) return;
                saveHistory(promptHistory);
            }
            checkStorage();
        } else {
             saveHistory(promptHistory);
        }
    }, [promptHistory]);

    const deleteFromHistory = useCallback((id: string) => {
        setPromptHistory(prev => prev.filter(item => item.id !== id));
    }, []);

    const clearHistory = useCallback(() => {
        setPromptHistory([]);
        // The useEffect will handle saving the empty array
    }, []);

    return { promptHistory, setPromptHistory, deleteFromHistory, clearHistory };
};